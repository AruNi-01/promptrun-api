package service

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"promptrun-api/common/errs"
	"promptrun-api/model"
	"promptrun-api/utils"
	"slices"
	"time"
)

// OrderListBySellerUserIdReq 根据卖家 userId 获取订单列表请求
type OrderListBySellerUserIdReq struct {
	Paginate     *utils.Page `json:"paginate"`
	SellerUserId int         `json:"sellerUserId"`
}

type CreateOrderReq struct {
	Id         int64     `json:"id"`
	PromptId   int       `json:"promptId"`
	SellerId   int       `json:"sellerId"`
	BuyerId    int       `json:"buyerId"`
	Price      float64   `json:"price"`
	CreateTime time.Time `json:"createTime"`
}

type OrderListAttachFullInfoRsp struct {
	OrderListAttachFullInfo []OrderAttachInfo `json:"orderListAttachFullInfo"`
	Rows                    int               `json:"rows"`
}

type OrderAttachInfo struct {
	Order         model.Order       `json:"order"`
	Buyer         model.User        `json:"buyer"`
	Prompt        model.Prompt      `json:"prompt"`
	Model         model.Model       `json:"model"`
	PromptImgList []model.PromptImg `json:"promptImgList"`
}

// ChartsRsp 根据卖家 userId 获取图表数据响应
type ChartsRsp struct {
	PublishPromptCount int     `json:"publishPromptCount"`
	SellPromptCount    int     `json:"sellPromptCount"`
	EarnMoney          float64 `json:"earnMoney"`
	BecomeSellerDay    int     `json:"becomeSellerDay"`

	SellMoneyEveryMonth []SellMoneyVo      `json:"sellMoneyEveryMonth"`
	SellCountEveryMonth []SellCountVo      `json:"sellCountEveryMonth"`
	SellModelRatio      []SellModelRatioVo `json:"sellModelRatio"`
}
type SellMoneyVo struct {
	Month     string  `json:"month"`
	SellMoney float64 `json:"sellMoney"`
}
type SellCountVo struct {
	Month     string `json:"month"`
	SellCount int    `json:"sellCount"`
}
type SellModelRatioVo struct {
	ModelName string `json:"name"`
	SellCount int    `json:"sellCount"`
}

type OrderListAttachPromptDetailRsp struct {
	Order        model.Order        `json:"order"`
	Prompt       model.Prompt       `json:"prompt"`
	PromptDetail model.PromptDetail `json:"promptDetail"`
}

func FindOrderListByBuyerId(c *gin.Context, buyerId int) ([]model.Order, *errs.Errs) {
	var orders []model.Order
	if model.DB.Where("buyer_id = ?", buyerId).Find(&orders).Error != nil {
		utils.Log().Error(c.FullPath(), "DB 获取订单列表失败")
		return nil, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取订单列表失败"))
	}
	return orders, nil
}

// FindOrderListAttachFullInfoBySellerUserId 根据卖家 userId 获取订单列表，附带详细信息
func (r *OrderListBySellerUserIdReq) FindOrderListAttachFullInfoBySellerUserId(c *gin.Context) (OrderListAttachFullInfoRsp, *errs.Errs) {
	// 1. 根据卖家 userId 获取 sellerId
	seller, e := FindSellerByUserId(c, r.SellerUserId)
	if e != nil {
		return OrderListAttachFullInfoRsp{}, e
	}

	// 2. 根据 sellerId 获取订单列表，分页查询，计算总行数
	var orders []model.Order
	query := model.DB.Model(model.Order{}).
		Where("seller_id = ?", seller.Id).
		Count(&r.Paginate.Rows).
		Order("create_time DESC")
	if query.Scopes(utils.Paginate(r.Paginate)).Find(&orders).Error != nil {
		utils.Log().Error(c.FullPath(), "DB 获取订单列表失败")
		return OrderListAttachFullInfoRsp{}, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取订单列表失败"))
	}

	// 3. 根据订单列表信息查询订单相关的详细信息
	var orderListAttachFullInfo []OrderAttachInfo
	for _, order := range orders {
		// 3.1 获取订单对应的买家信息
		buyer, e := FindUserById(c, order.BuyerId)
		if e != nil {
			buyer = model.User{}
		}

		// 3.2 获取订单对应的提示词信息
		prompt, e := FindPromptById(c, order.PromptId)
		if e != nil {
			prompt = model.Prompt{}
		}

		// 3.3 获取订单对应的提示词模型
		model2, e := FindModelById(c, prompt.ModelId)
		if e != nil {
			model2 = model.Model{}
		}

		// 3.4 获取订单对应的提示词图片列表
		promptImgList, e := FindPromptImgListByPromptId(c, order.PromptId)
		if e != nil {
			promptImgList = nil
		}

		orderListAttachFullInfo = append(orderListAttachFullInfo, OrderAttachInfo{
			Order:         order,
			Buyer:         buyer,
			Prompt:        prompt,
			Model:         model2,
			PromptImgList: promptImgList,
		})
	}

	return OrderListAttachFullInfoRsp{
		OrderListAttachFullInfo: orderListAttachFullInfo,
		Rows:                    r.Paginate.Rows,
	}, nil
}

// FindChartsFullInfoBySellerUserId 根据卖家 userId 获取图表数据
func FindChartsFullInfoBySellerUserId(c *gin.Context, sellerUserId int) (ChartsRsp, *errs.Errs) {
	// 1. 根据卖家 userId 获取 seller
	seller, e := FindSellerByUserId(c, sellerUserId)
	if e != nil {
		return ChartsRsp{}, e
	}

	// 2. 获取每月销售金额
	sellMoneyEveryMonth, _ := findSellMoneyEveryMonth(c, seller.Id)

	// 3. 获取每月销售数量
	sellCountEveryMonth, _ := findSellCountEveryMonth(c, seller.Id)

	// 4. 获取每个模型的销售比例
	sellModelRatio, _ := findSellModelRatio(c, seller.Id)

	// 5. 获取卖家的发布提示词数量
	publishPromptCount, _ := FindPromptCountBySellerId(c, seller.Id)

	// 6. 获取卖家的销售提示词数量
	sellPromptCount, _ := findSellPromptCountBySellerId(c, seller.Id)

	// 7. 获取卖家的总收益
	earnMoney, _ := findEarnMoneyBySellerId(c, seller.Id)

	// 8. 获取卖家成为卖家的天数
	becomeSellerDay, _ := FindBecomeSellerDayBySellerId(c, seller.Id)

	return ChartsRsp{
		SellMoneyEveryMonth: sellMoneyEveryMonth,
		SellCountEveryMonth: sellCountEveryMonth,
		SellModelRatio:      sellModelRatio,

		PublishPromptCount: publishPromptCount,
		SellPromptCount:    sellPromptCount,
		EarnMoney:          earnMoney,
		BecomeSellerDay:    becomeSellerDay,
	}, nil
}

func findEarnMoneyBySellerId(c *gin.Context, sellerId int) (float64, *errs.Errs) {
	var earnMoney float64
	if e := model.DB.Model(model.Order{}).Where("seller_id = ?", sellerId).Select("SUM(price)").Row().Scan(&earnMoney); e != nil {
		utils.Log().Error(c.FullPath(), "DB 获取卖家总收益失败")
		return 0, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取卖家总收益失败"))
	}
	return earnMoney, nil
}

func findSellPromptCountBySellerId(c *gin.Context, sellerId int) (int, *errs.Errs) {
	var count int
	if e := model.DB.Model(model.Order{}).Where("seller_id = ?", sellerId).Count(&count).Error; e != nil {
		utils.Log().Error(c.FullPath(), "DB 获取卖家销售提示词数量失败")
		return 0, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取卖家销售提示词数量失败"))
	}
	return count, nil
}

func findSellModelRatio(c *gin.Context, sellId int) ([]SellModelRatioVo, *errs.Errs) {
	var sellModelRatio []SellModelRatioVo

	// 获取模型 id 和销售数量
	rows, e := model.DB.
		Raw("SELECT prompt_id FROM `order` WHERE seller_id = ?", sellId).Rows()
	if e != nil {
		utils.Log().Error(c.FullPath(), "DB 获取每个模型的销售比例失败")
		return nil, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取每个模型的销售比例失败"))
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			utils.Log().Error(c.FullPath(), "关闭 rows 失败")
		}
	}(rows)

	modelCountMap := map[string]int{}
	for rows.Next() {
		var promptId int
		if e := rows.Scan(&promptId); e != nil {
			utils.Log().Error(c.FullPath(), "DB 获取每个模型的销售比例失败")
			return nil, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取每个模型的销售比例失败"))
		}

		prompt, e2 := FindPromptById(c, promptId)
		if e2 != nil {
			return nil, e2
		}

		model2, e3 := FindModelById(c, prompt.ModelId)
		if e3 != nil {
			return nil, e3
		}

		if _, ok := modelCountMap[model2.Name]; ok {
			modelCountMap[model2.Name] += 1
		} else {
			modelCountMap[model2.Name] = 1
		}
	}

	for modelName, count := range modelCountMap {
		sellModelRatio = append(sellModelRatio, SellModelRatioVo{
			ModelName: modelName,
			SellCount: count,
		})
	}

	return sellModelRatio, nil
}

func findSellCountEveryMonth(c *gin.Context, sellId int) ([]SellCountVo, *errs.Errs) {
	var sellCountEveryMonth []SellCountVo
	rows, e := model.DB.
		Raw("SELECT DATE_FORMAT(create_time, '%y-%m') AS month, COUNT(*) AS sellCount FROM `order` WHERE seller_id = ? GROUP BY month ORDER BY month DESC LIMIT 12", sellId).Rows()
	if e != nil {
		utils.Log().Error(c.FullPath(), "DB 获取每月销售数量失败")
		return nil, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取每月销售数量失败"))
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			utils.Log().Error(c.FullPath(), "关闭 rows 失败")
		}
	}(rows)

	for rows.Next() {
		var sellCountVo SellCountVo
		if e := rows.Scan(&sellCountVo.Month, &sellCountVo.SellCount); e != nil {
			utils.Log().Error(c.FullPath(), "DB 获取每月销售数量失败")
			return nil, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取每月销售数量失败"))
		}
		sellCountEveryMonth = append(sellCountEveryMonth, sellCountVo)
	}
	slices.Reverse(sellCountEveryMonth)
	return sellCountEveryMonth, nil
}

func findSellMoneyEveryMonth(c *gin.Context, sellId int) ([]SellMoneyVo, *errs.Errs) {
	var sellMoneyEveryMonth []SellMoneyVo
	rows, e := model.DB.
		Raw("SELECT DATE_FORMAT(create_time, '%y-%m') AS month, SUM(price) AS sellMoney FROM `order` WHERE seller_id = ? GROUP BY month ORDER BY month DESC LIMIT 12", sellId).Rows()
	if e != nil {
		utils.Log().Error(c.FullPath(), "DB 获取每月销售金额失败")
		return nil, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取每月销售金额失败"))
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			utils.Log().Error(c.FullPath(), "关闭 rows 失败")
		}
	}(rows)

	for rows.Next() {
		var sellMoneyVo SellMoneyVo
		if e := rows.Scan(&sellMoneyVo.Month, &sellMoneyVo.SellMoney); e != nil {
			utils.Log().Error(c.FullPath(), "DB 获取每月销售金额失败")
			return nil, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取每月销售金额失败"))
		}
		sellMoneyEveryMonth = append(sellMoneyEveryMonth, sellMoneyVo)
	}
	slices.Reverse(sellMoneyEveryMonth)
	return sellMoneyEveryMonth, nil
}

func FindOrderListAttachPromptDetailById(c *gin.Context, orderId int64) (OrderListAttachPromptDetailRsp, *errs.Errs) {
	order, e := FindOrderById(c, orderId)
	if e != nil {
		return OrderListAttachPromptDetailRsp{}, e
	}

	prompt, e := FindPromptById(c, order.PromptId)
	if e != nil {
		return OrderListAttachPromptDetailRsp{}, e
	}

	promptDetail, e := FindPromptDetailByPromptId(c, prompt.Id)
	if e != nil {
		return OrderListAttachPromptDetailRsp{}, e
	}

	return OrderListAttachPromptDetailRsp{
		Order:        order,
		Prompt:       prompt,
		PromptDetail: promptDetail,
	}, nil
}

func FindOrderById(c *gin.Context, orderId int64) (model.Order, *errs.Errs) {
	var order model.Order
	if err := model.DB.Where("id = ?", orderId).First(&order).Error; err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.Log().Error(c.FullPath(), "DB 订单不存在")
			return model.Order{}, errs.NewErrs(errs.ErrRecordNotFound, errors.New("DB 订单不存在"))
		}
		utils.Log().Error(c.FullPath(), "DB 获取订单失败")
		return model.Order{}, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取订单失败"))
	}
	return order, nil
}

func OrderRatingById(c *gin.Context, orderId int64, rating float64) (model.Order, *errs.Errs) {
	var order model.Order
	if model.DB.Where("id = ?", orderId).First(&order).Error != nil {
		utils.Log().Error(c.FullPath(), "DB 获取订单失败")
		return model.Order{}, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取订单失败"))
	}
	order.IsRating = model.OrderRatingDone
	order.Rating = rating

	if e := model.DB.Save(&order).Error; e != nil {
		utils.Log().Error(c.FullPath(), "DB 评价订单失败")
		return model.Order{}, errs.NewErrs(errs.ErrDBError, errors.New("DB 评价订单失败"))
	}

	// TODO：异步插入订单评分表，后续在低峰期使用定时任务扫描评分表，计算卖家和 Prompt 平均评分
	go func() {
		if _, e := AddOrderRating(c, order, rating); e != nil {
			utils.Log().Error(c.FullPath(), "异步插入订单评分表失败")
		}
	}()

	go OrderRatingMsgNotice(c, order)

	return order, nil
}

func (r *CreateOrderReq) CreateOrder(c *gin.Context) (bool, *errs.Errs) {
	order := model.Order{
		Id:         r.Id,
		PromptId:   r.PromptId,
		SellerId:   r.SellerId,
		BuyerId:    r.BuyerId,
		Price:      r.Price,
		CreateTime: r.CreateTime,
	}

	if e := model.DB.Create(&order).Error; e != nil {
		utils.Log().Error(c.FullPath(), "DB 创建订单失败, err: %s", e.Error())
		return false, errs.NewErrs(errs.ErrDBError, errors.New("DB 创建订单失败"))
	}
	return true, nil
}
