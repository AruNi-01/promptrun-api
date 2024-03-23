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
	SortBy        string      `json:"sortBy"`
}

// PromptListResp 获取提示词列表响应
type PromptListResp struct {
	Paginate *utils.Page    `json:"paginate"`
	Prompts  []model.Prompt `json:"prompts"`
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

	if query.Scopes(
		utils.Paginate(r.Paginate)).
		Where("audit_status = ?", model.AuditStatusPass).
		Where("publish_status = ?", model.PublishStatusOn).
		Find(&prompts).
		Count(&r.Paginate.Rows).Error != nil {
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
