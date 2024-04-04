package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"promptrun-api/common/errs"
	"promptrun-api/model"
	"promptrun-api/utils"
)

// OrderListBySellerUserIdReq 根据卖家 userId 获取订单列表请求
type OrderListBySellerUserIdReq struct {
	Paginate     *utils.Page `json:"paginate"`
	SellerUserId int         `json:"sellerUserId"`
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

func FindOrderListByBuyerId(c *gin.Context, buyerId int) ([]model.Order, *errs.Errs) {
	var orders []model.Order
	if model.DB.Where("buyer_id = ?", buyerId).Find(&orders).Error != nil {
		utils.Log().Error(c.FullPath(), "DB 获取订单列表失败")
		return nil, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取订单列表失败"))
	}
	return orders, nil
}

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
