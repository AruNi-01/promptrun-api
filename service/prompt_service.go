package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"promptrun-api/common/errs"
	"promptrun-api/model"
	"promptrun-api/utils"
)

// SortBy 前端传入的排序方式
const (
	SortByHot        = "hot"         // 按照热度
	SortByTime       = "time"        // 按照时间
	SortBySellAmount = "sell_amount" // 按照销量
)

// PromptListReq 获取提示词列表请求
type PromptListReq struct {
	Paginate      *utils.Page `json:"paginate"`
	ModelId       *int        `json:"modelId,omitempty"`
	CategoryTypes []int       `json:"categoryTypes,omitempty"`
	SortBy        string      `json:"sortBy,omitempty"`

	PublishStatus []int `json:"publishStatus,omitempty"`
	AuditStatus   []int `json:"auditStatus,omitempty"`

	UserId *int `json:"userId,omitempty"` // 用户 ID，用于获取某个用户（seller）发布的提示词
}

// PromptMasterImgListReq 获取主图列表请求
type PromptMasterImgListReq struct {
	PromptIds []int `json:"promptIds" binding:"required"`
}

// PromptListByBuyerIdReq 根据买家 id 获取提示词列表请求
type PromptListByBuyerIdReq struct {
	Paginate *utils.Page `json:"paginate"`
	BuyerId  int         `json:"buyerId"`
}

// PromptListResp 获取提示词列表响应
type PromptListResp struct {
	Prompts []model.Prompt `json:"prompts"`
	Rows    int            `json:"rows"`
}

type PromptDetailResp struct {
	Prompt        model.Prompt      `json:"prompt"`
	Seller        model.Seller      `json:"seller"`
	Model         model.Model       `json:"model"`
	PromptImgList []model.PromptImg `json:"promptImgList"`
}

// PromptList 获取提示词列表
func (r *PromptListReq) PromptList(c *gin.Context) ([]model.Prompt, *errs.Errs) {
	var prompts []model.Prompt
	query := model.DB.Model(model.Prompt{})

	if r.ModelId != nil {
		query = query.Where("model_id = ?", *r.ModelId)
	}
	if r.CategoryTypes != nil {
		query = query.Where("category_type IN (?)", r.CategoryTypes)
	}
	switch r.SortBy {
	case SortByHot:
		query = query.Order("score DESC")
	case SortBySellAmount:
		query = query.Order("sell_amount DESC")
	case SortByTime:
		query = query.Order("create_time DESC")
	default:
		query = query.Order("create_time DESC")
	}

	if r.PublishStatus != nil {
		query = query.Where("publish_status IN (?)", r.PublishStatus)
	}
	if r.AuditStatus != nil {
		query = query.Where("audit_status IN (?)", r.AuditStatus)
	}

	if r.UserId != nil {
		seller, e := FindSellerByUserId(*r.UserId)
		if e != nil {
			return nil, e
		}
		query = query.Where("seller_id = ?", seller.Id)
	}

	// 注意：Count 要放在分页查询之前，否则会导致 count 为空
	query.Count(&r.Paginate.Rows)

	if query.Scopes(
		utils.Paginate(r.Paginate)).
		Find(&prompts).
		Error != nil {
		return nil, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取提示词列表失败"))
	}

	return prompts, nil
}

// FindPromptFullInfoById 根据 ID 查找提示词的详细信息
func FindPromptFullInfoById(c *gin.Context, promptId int) (PromptDetailResp, *errs.Errs) {
	prompt, e := FindPromptById(promptId)
	if e != nil {
		return PromptDetailResp{}, e
	}
	seller, e := FindSellerById(prompt.SellerId)
	if e != nil {
		return PromptDetailResp{}, e
	}
	promptModel, e := FindModelById(prompt.ModelId)
	if e != nil {
		return PromptDetailResp{}, e
	}
	promptImgList, e := FindPromptImgListByPromptId(c, promptId)
	if e != nil {
		return PromptDetailResp{}, e
	}
	return PromptDetailResp{
		Prompt:        prompt,
		Seller:        seller,
		Model:         promptModel,
		PromptImgList: promptImgList,
	}, nil

}

// FindPromptById 根据 ID 查找提示词
func FindPromptById(promptId int) (model.Prompt, *errs.Errs) {
	var prompt model.Prompt
	if model.DB.First(&prompt, promptId).RecordNotFound() {
		return model.Prompt{}, errs.NewErrs(errs.ErrRecordNotFound, errors.New("提示词不存在"))
	}
	return prompt, nil
}

func FindPromptImgListByPromptId(c *gin.Context, id int) ([]model.PromptImg, *errs.Errs) {
	var promptImgList []model.PromptImg
	if model.DB.Where("prompt_id = ?", id).Find(&promptImgList).Error != nil {
		utils.Log().Error(c.FullPath(), "DB 获取提示词图片列表失败")
		return nil, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取提示词图片列表失败"))
	}
	return promptImgList, nil
}

func FindPromptMasterImgByPromptId(promptId int) (model.PromptImg, *errs.Errs) {
	var promptImg model.PromptImg
	if model.DB.Where("prompt_id = ?", promptId).Where("is_master = ?", 1).First(&promptImg).RecordNotFound() {
		return model.PromptImg{}, errs.NewErrs(errs.ErrRecordNotFound, errors.New("未找到主图"))
	}
	return promptImg, nil
}

func (r *PromptMasterImgListReq) FindMasterImgListByPromptIds(c *gin.Context) ([]model.PromptImg, *errs.Errs) {
	var promptImgList []model.PromptImg
	if model.DB.Where("prompt_id IN (?)", r.PromptIds).Where("is_master = ?", 1).Find(&promptImgList).Error != nil {
		return nil, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取主图列表失败"))
	}
	return promptImgList, nil
}

func (r PromptListByBuyerIdReq) FindListByBuyerId(c *gin.Context) ([]model.Prompt, *errs.Errs) {
	orderList, err := FindOrderListByBuyerId(c, r.BuyerId)
	if err != nil {
		return nil, err
	}

	var promptIds = func() []int {
		var ids []int
		for _, order := range orderList {
			ids = append(ids, order.PromptId)
		}
		return ids
	}()

	var prompts []model.Prompt
	query := model.DB.
		Model(model.Prompt{}).
		Where("id in (?)", promptIds).
		Order("create_time DESC").
		Count(&r.Paginate.Rows)

	if query.Scopes(utils.Paginate(r.Paginate)).Find(&prompts).Error != nil {
		return nil, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取提示词列表失败"))
	}
	return prompts, nil
}
