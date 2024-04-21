package lantu_pay

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"promptrun-api/utils"
	"strconv"
	"time"
)

const (
	lantuPayUrl           = "https://api.ltzf.cn/api/wxpay/native"
	lantuPayOrderQueryUrl = "https://api.ltzf.cn/api/wxpay/get_pay_order"
)

const (
	PayTimeExpire2Minutes  = "2m"
	PayTimeExpire5Minutes  = "5m"
	PayTimeExpire10Minutes = "10m"
	PayTimeExpire15Minutes = "15m"
	PayTimeExpire30Minutes = "30m"

	PayTimeExpire1Hours = "1h"
	PayTimeExpire2Hours = "2h"
)

const (
	PayReturnSuccess = 0
	PayReturnFail    = 1
)

const (
	PayOrderQuerySuccess = 0
	PayOrderQueryFail    = 1
)

const (
	PayStatusNotPay = 0
	PayStatusPaid   = 1
)

type LantuWxPayReq struct {
	OrderId    string // 订单 id，从业务传入，跟 order 表的 id 保持一致
	TotalFee   string // 支付金额
	GoodsDesc  string // 商品描述
	Attach     string // 附加数据，在支付通知中原样返回，可作为自定义参数使用。
	TimeExpire string // 订单失效时间，枚举值。m：分钟 h：小时 取值范围：1m～2h（接口请求后开始计算时间）
}

type LantuWxPayQueryOrderReq struct {
	OrderId string // 订单 id，从业务传入，跟 order 表的 id 保持一致
}

type LantuWxPayResp struct {
	Code      int     `json:"code"` //返回状态，枚举值。0：成功 1：失败
	Data      PayData `json:"data"`
	Msg       string  `json:"msg"`
	RequestID string  `json:"request_id"`
}
type PayData struct {
	CodeURL   string `json:"code_url"`   // 微信原生支付链接，此 URL 用于生成支付二维码，然后提供给用户扫码支付。
	QRCodeURL string `json:"QRcode_url"` // 蓝兔支付生成的二维码链接地址
}

type LantuWxPayQueryOrderResp struct {
	Code      int            `json:"code"` //返回状态，枚举值。0：成功 1：失败
	Data      OrderQueryData `json:"data"`
	Msg       string         `json:"msg"`
	RequestID string         `json:"request_id"`
}
type OrderQueryData struct {
	AddTime     string `json:"add_time"`     // 下单时间
	MchId       string `json:"mch_id"`       // 商户号
	OrderNo     string `json:"order_no"`     // 蓝兔系统订单号
	OutTradeNo  string `json:"out_trade_no"` // 商户订单号
	PayNo       string `json:"pay_no"`       // 微信支付订单号
	Body        string `json:"body"`         // 商品描述
	TotalFee    string `json:"total_fee"`    // 支付金额
	TradeType   string `json:"trade_type"`   // 交易类型，JSAPI：公众号支付，NATIVE：扫码支付，APP：APP 支付，H5：H5 支付，MINIPROGRAM：小程序支付
	SuccessTime string `json:"success_time"` // 支付成功时间
	Attach      string `json:"attach"`       // 附加数据，在支付接口中填写的数据，可作为自定义参数使用。
	OpenId      string `json:"open_id"`      // 支付者信息，在此商户下的唯一标识
	PayStatus   int    `json:"pay_status"`   // 支付状态，0：未支付 1：已支付
}

// LantuWxPay 蓝兔微信支付接口，返回支付二维码
func (r *LantuWxPayReq) LantuWxPay() (LantuWxPayResp, error) {
	// 调用蓝兔支付接口，获取支付二维码
	resp, err := http.PostForm(lantuPayUrl, r.getLantuWxPayApiReq())
	if err != nil {
		return LantuWxPayResp{}, err
	}

	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			utils.Log().Error("", "close resp body io fail, err: %s", err.Error())
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LantuWxPayResp{}, err
	}

	var lantuWxPayResp LantuWxPayResp
	if err := json.Unmarshal(body, &lantuWxPayResp); err != nil {
		return LantuWxPayResp{}, err
	}

	if lantuWxPayResp.Code != PayReturnSuccess {
		return LantuWxPayResp{}, errors.New(lantuWxPayResp.Msg)
	}

	return lantuWxPayResp, nil
}

// LantuWxPayQueryOrder 蓝兔支付订单查询接口
func (r *LantuWxPayQueryOrderReq) LantuWxPayQueryOrder() (LantuWxPayQueryOrderResp, error) {
	// 调用蓝兔支付订单查询接口，获取订单信息
	resp, err := http.PostForm(lantuPayOrderQueryUrl, r.getLantuWxPayQueryOrderApiReq())
	if err != nil {
		return LantuWxPayQueryOrderResp{}, err
	}

	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			utils.Log().Error("", "close resp body io fail, err: %s", err.Error())
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LantuWxPayQueryOrderResp{}, err
	}

	var lantuWxPayQueryOrderResp LantuWxPayQueryOrderResp
	if err := json.Unmarshal(body, &lantuWxPayQueryOrderResp); err != nil {
		return LantuWxPayQueryOrderResp{}, err
	}

	if lantuWxPayQueryOrderResp.Code != PayOrderQuerySuccess {
		return LantuWxPayQueryOrderResp{}, errors.New(lantuWxPayQueryOrderResp.Msg)
	}

	return lantuWxPayQueryOrderResp, nil
}

// 获取蓝兔支付 api 的请求参数
func (r *LantuWxPayReq) getLantuWxPayApiReq() url.Values {
	opts := map[string]string{
		"mch_id":       os.Getenv("MCH_ID"),
		"out_trade_no": r.OrderId,
		"total_fee":    r.TotalFee,
		"body":         r.GoodsDesc,
		"timestamp":    strconv.FormatInt(time.Now().Unix(), 10),
		"notify_url":   os.Getenv("NOTIFY_URL"),
		"attach":       r.Attach,
		"time_expire":  r.TimeExpire,
		"sign":         "",
	}
	sign := genSign(r.getLantuWxPaySignParams(opts))
	opts["sign"] = sign

	req := url.Values{}
	for key, value := range opts {
		req.Add(key, value)
	}
	return req
}

// 获取蓝兔支付订单查询 api 的请求参数
func (r *LantuWxPayQueryOrderReq) getLantuWxPayQueryOrderApiReq() url.Values {
	opts := map[string]string{
		"mch_id":       os.Getenv("MCH_ID"),
		"out_trade_no": r.OrderId,
		"timestamp":    strconv.FormatInt(time.Now().Unix(), 10),
		"sign":         "",
	}
	sign := genSign(r.getLantuWxPayQueryOrderSignParams(opts))
	opts["sign"] = sign

	req := url.Values{}
	for key, value := range opts {
		req.Add(key, value)
	}
	return req
}

// 获取 LantuWxPaySign 签名所需的参数，只有必填参数才参与签名
func (r *LantuWxPayReq) getLantuWxPaySignParams(params map[string]string) map[string]string {
	return map[string]string{
		"mch_id":       params["mch_id"],
		"out_trade_no": params["out_trade_no"],
		"total_fee":    params["total_fee"],
		"body":         params["body"],
		"timestamp":    params["timestamp"],
		"notify_url":   params["notify_url"],
	}
}

// 获取 LantuWxPayQueryOrderSignParams 签名所需的参数，只有必填参数才参与签名
func (r *LantuWxPayQueryOrderReq) getLantuWxPayQueryOrderSignParams(params map[string]string) map[string]string {
	return map[string]string{
		"mch_id":       params["mch_id"],
		"out_trade_no": params["out_trade_no"],
		"timestamp":    params["timestamp"],
	}
}
