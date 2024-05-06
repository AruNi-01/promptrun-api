package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"promptrun-api/common/errs"
	"promptrun-api/model"
	"promptrun-api/third_party/lantu_pay"
	"promptrun-api/utils"
	"strconv"
	"time"
)

const (
	LantuWxPayNotifySuccess = "0" // LantuWxPayNotifySuccess 支付成功
	LantuWxPayNotifyFail    = "1" // LantuWxPayNotifyFail 支付失败
)

// BasePayReq 支付请求，基础参数 struct
type BasePayReq struct {
	PromptId     int     `json:"promptId"`
	PromptTitle  string  `json:"promptTitle"`
	SellerId     int     `json:"sellerId"`
	SellerUserId int     `json:"sellerUserId"`
	BuyerId      int     `json:"buyerId"`
	Price        float64 `json:"price"`
}

// LantuWxPayReq 蓝兔微信支付请求参数，同时也作为 Attach 附加数据类型，用于回调通知时获取订单信息
type LantuWxPayReq struct {
	BasePayReq
}

type LantuWxPayQueryOrderReq struct {
	OrderId int64 `json:"orderId"` // 订单 id，从业务传入，跟 order 表的 id 保持一致
}

type LantuWxPayResp struct {
	QRCodeURL string `json:"qrCodeUrl"`
	OrderId   int64  `json:"orderId"`
}

type LantuWxPayQueryOrderResp struct {
	IsPay   bool          `json:"isPay"`   // 是否已支付
	OrderId int64         `json:"orderId"` // 订单 id，从业务传入，跟 order 表的 id 保持一致
	PayTime time.Time     `json:"payTime"` // 支付时间
	Attach  LantuWxPayReq `json:"attach"`  // 附加数据，在支付接口中填写的数据，可作为自定义参数使用。
}

type LantuWxPayNotifyParams struct {
	Code        string `json:"code"`         // 支付结果，0：支付成功，1：支付失败
	Timestamp   string `json:"timestamp"`    // 时间戳
	MchId       string `json:"mch_id"`       // 商户号
	OrderNo     string `json:"order_no"`     // 蓝兔系统订单号
	OutTradeNo  string `json:"out_trade_no"` // 商户订单号，订单 id，跟 order 表的 id 保持一致
	PayNo       string `json:"pay_no"`       // 微信支付订单号
	TotalFee    string `json:"total_fee"`    // 支付金额
	Sign        string `json:"sign"`         // 签名
	PayChannel  string `json:"pay_channel"`  // 支付渠道，wxpay：微信支付，alipay：支付宝支付
	TradeType   string `json:"trade_type"`   // 交易类型，JSAPI：公众号支付，NATIVE：扫码支付，APP：APP 支付，H5：H5 支付，MINIPROGRAM：小程序支付
	SuccessTime string `json:"success_time"` // 支付成功时间
	Attach      string `json:"attach"`       // 附加数据，在支付接口中填写的数据，可作为自定义参数使用。
	OpenId      string `json:"open_id"`      // 支付者信息，在此商户下的唯一标识
}

// BalancePayReq 余额支付请求参数
type BalancePayReq struct {
	BasePayReq
	Wallet model.Wallet `json:"wallet"`
}

type BalancePayResp struct {
	OrderId int64 `json:"orderId"`
}

// LantuWxPay 蓝兔微信支付接口，返回支付二维码
func (r *LantuWxPayReq) LantuWxPay(c *gin.Context) (LantuWxPayResp, *errs.Errs) {
	orderId := utils.GenSnowflakeId()

	lantuWxPayReq := lantu_pay.LantuWxPayReq{
		OrderId:   strconv.FormatInt(orderId, 10),
		TotalFee:  strconv.FormatFloat(r.Price, 'f', 2, 64),
		GoodsDesc: "PromptRun 平台 — " + r.PromptTitle,
		// 附加数据，用于回调通知时获取订单信息
		Attach: func(c *gin.Context, r *LantuWxPayReq) string {
			marshal, err := json.Marshal(r)
			if err != nil {
				utils.Log().Error(c.FullPath(), "json.Marshal error: %s", err.Error())
				return ""
			} else {
				return string(marshal)
			}
		}(c, r),
		TimeExpire: lantu_pay.PayTimeExpire2Minutes,
	}
	resp, err := lantuWxPayReq.LantuWxPay()
	if err != nil {
		utils.Log().Error(c.FullPath(), "lantuWxPayReq.LantuWxPay error: %s", err.Error())
		return LantuWxPayResp{}, errs.NewErrs(errs.ErrPayFacadeError, errors.New("调用蓝兔支付接口失败"))
	}
	return LantuWxPayResp{
		QRCodeURL: resp.Data.QRCodeURL,
		OrderId:   orderId,
	}, nil
}

// LantuWxPayNotify 蓝兔微信支付回调通知
func (r *LantuWxPayNotifyParams) LantuWxPayNotify(c *gin.Context) (bool, *errs.Errs) {
	if r.Code != LantuWxPayNotifySuccess {
		utils.Log().Info(c.FullPath(), "蓝兔微信支付回调通知结果：支付失败")
		return false, nil
	}

	// 从附加数据中获取订单信息
	var orderInfo LantuWxPayReq
	if err := json.Unmarshal([]byte(r.Attach), &orderInfo); err != nil {
		utils.Log().Error(c.FullPath(), "json.Unmarshal error: %s", err.Error())
		return false, nil
	}

	// 若订单已创建（查询蓝兔支付订单接口逻辑 LantuWxPayQueryOrder），则直接返回
	orderId, err := strconv.ParseInt(r.OutTradeNo, 10, 64)
	if err != nil {
		utils.Log().Error(c.FullPath(), "strconv.Atoi error: %s", err.Error())
		return false, nil
	}
	order, e := FindOrderById(c, orderId)
	// 订单已存在，直接返回
	if order.Id != 0 {
		return false, e
	}

	// 支付成功，创建订单
	createOrderReq := CreateOrderReq{
		Id:       orderId,
		PromptId: orderInfo.PromptId,
		SellerId: orderInfo.SellerId,
		BuyerId:  orderInfo.BuyerId,
		Price:    orderInfo.Price,
		CreateTime: func(c *gin.Context, r *LantuWxPayNotifyParams) time.Time {
			timestampInt, err := strconv.ParseInt(r.Timestamp, 10, 64)
			if err != nil {
				fmt.Println("Error parsing timestamp:", err)
				return time.Now()
			}
			return time.Unix(timestampInt, 0)
		}(c, r),
	}
	_, e = createOrderReq.CreateOrder(c)
	if e != nil {
		return false, e
	}
	return true, nil
}

// LantuWxPayQueryOrder 蓝兔支付订单查询接口
func (r *LantuWxPayQueryOrderReq) LantuWxPayQueryOrder(c *gin.Context) (LantuWxPayQueryOrderResp, *errs.Errs) {
	lantuWxPayQueryOrderReq := lantu_pay.LantuWxPayQueryOrderReq{
		OrderId: strconv.FormatInt(r.OrderId, 10),
	}

	// 调用蓝兔支付订单查询接口，获取订单信息
	resp, err := lantuWxPayQueryOrderReq.LantuWxPayQueryOrder()
	if err != nil {
		utils.Log().Error("", "lantuWxPayQueryOrderReq.LantuWxPayQueryOrder error: %s", err.Error())
		return LantuWxPayQueryOrderResp{}, errs.NewErrs(errs.ErrPayOrderQueryError, errors.New("调用蓝兔支付订单查询接口失败"))
	}
	lantuWxPayQueryOrderResp := LantuWxPayQueryOrderResp{
		IsPay:   resp.Data.PayStatus == lantu_pay.PayStatusPaid,
		OrderId: r.OrderId,
		PayTime: func(payTime string) time.Time {
			if payTime == "" {
				return time.Now()
			}
			timestampInt, err := strconv.ParseInt(payTime, 10, 64)
			if err != nil {
				utils.Log().Error("", "strconv.Atoi error: %s", err.Error())
				return time.Now()
			}
			return time.Unix(timestampInt, 0)
		}(resp.Data.SuccessTime),
		Attach: func(c *gin.Context, attach string) LantuWxPayReq {
			orderInfo := LantuWxPayReq{}
			if err := json.Unmarshal([]byte(attach), &orderInfo); err != nil {
				utils.Log().Error(c.FullPath(), "json.Unmarshal error: %s", err.Error())
			}
			return orderInfo
		}(c, resp.Data.Attach),
	}

	// 支付成功，计算钱包、异步创建订单和生成账单。TODO：对于创建失败的订单，可以通过订单查询接口再次创建来补偿
	if lantuWxPayQueryOrderResp.IsPay {
		// 计算卖家在平台的余额和收入总额
		if _, e := CalculateBalanceAndIncome(c, lantuWxPayQueryOrderResp.Attach.SellerUserId,
			lantuWxPayQueryOrderResp.Attach.Price, lantuWxPayQueryOrderResp.Attach.Price); e != nil {
			return LantuWxPayQueryOrderResp{}, e
		}
		// 计算买家在平台的支出总额
		if _, e := CalculateBalanceAndOutcome(c, lantuWxPayQueryOrderResp.Attach.BuyerId,
			0, lantuWxPayQueryOrderResp.Attach.Price); e != nil {
			return LantuWxPayQueryOrderResp{}, e
		}

		go genOrderAndBill(
			c,
			CreateOrderReq{
				// 从附加数据中获取订单信息
				Id:         lantuWxPayQueryOrderResp.OrderId,
				PromptId:   lantuWxPayQueryOrderResp.Attach.PromptId,
				SellerId:   lantuWxPayQueryOrderResp.Attach.SellerId,
				BuyerId:    lantuWxPayQueryOrderResp.Attach.BuyerId,
				Price:      lantuWxPayQueryOrderResp.Attach.Price,
				CreateTime: lantuWxPayQueryOrderResp.PayTime},
			lantuWxPayQueryOrderResp.Attach.SellerUserId,
			lantuWxPayQueryOrderResp.Attach.PromptTitle,
			model.BillChannelBalance,
			model.BillChannelWxPay,
		)

		go PromptSoldMsgNotice(c, lantuWxPayQueryOrderResp.Attach.PromptTitle, lantuWxPayQueryOrderResp.Attach.SellerUserId, lantuWxPayQueryOrderResp.Attach.BuyerId)
	}

	return lantuWxPayQueryOrderResp, nil
}

// BalancePay 余额支付
func (r *BalancePayReq) BalancePay(c *gin.Context) (BalancePayResp, *errs.Errs) {
	// 余额支付
	if r.Wallet.Balance < r.Price {
		return BalancePayResp{}, errs.NewErrs(errs.ErrBalanceNotEnough, errors.New("账户余额不足"))
	}
	// 计算余额和收入，余额增加，收入增加
	if _, e := CalculateBalanceAndIncome(c, r.SellerUserId, r.Price, r.Price); e != nil {
		return BalancePayResp{}, e
	}
	// 计算余额和支出，余额减少，支出增加
	if _, e := CalculateBalanceAndOutcome(c, r.BuyerId, r.Price, r.Price); e != nil {
		return BalancePayResp{}, e
	}

	orderId := utils.GenSnowflakeId()

	// 创建订单，生成账单
	go genOrderAndBill(
		c,
		CreateOrderReq{
			Id:         orderId,
			PromptId:   r.PromptId,
			SellerId:   r.SellerId,
			BuyerId:    r.BuyerId,
			Price:      r.Price,
			CreateTime: time.Now(),
		},
		r.SellerUserId,
		r.PromptTitle,
		model.BillChannelBalance,
		model.BillChannelBalance,
	)
	go IncreaseSellerSellAmount(c, r.SellerId)

	go PromptSoldMsgNotice(c, r.PromptTitle, r.SellerUserId, r.BuyerId)

	return BalancePayResp{
		OrderId: orderId,
	}, nil
}

func genOrderAndBill(c *gin.Context, createOrderReq CreateOrderReq,
	sellerUserId int, promptTitle string, incomeChannel, outcomeChannel int) {
	_, e := createOrderReq.CreateOrder(c)
	if e != nil {
		utils.Log().Error(c.FullPath(), "createOrderReq.CreateOrder error: %s", e.Err.Error())
	}

	// 生成账单
	bill := model.Bill{
		UserId:     sellerUserId,
		Type:       model.BillTypeIncome,
		Amount:     createOrderReq.Price,
		Channel:    incomeChannel,
		Remark:     fmt.Sprintf("售出 Prompt - %s", promptTitle),
		CreateTime: time.Now(),
	}
	if _, e := AddBill(c, bill); e != nil {
		utils.Log().Error(c.FullPath(), "Add seller bill error: %s", e.Err.Error())
	}

	bill.UserId = createOrderReq.BuyerId
	bill.Type = model.BillTypeOutcome
	bill.Channel = outcomeChannel
	bill.Remark = fmt.Sprintf("购买 Prompt - %s", promptTitle)
	if _, e := AddBill(c, bill); e != nil {
		utils.Log().Error(c.FullPath(), "Add buyer bill error: %s", e.Err.Error())
	}
}
