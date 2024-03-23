package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"promptrun-api/common/errs"
	"promptrun-api/model"
	"promptrun-api/utils"
)

const (
	SortByHot        = "hot"
	SortByTime       = "time"
	SortBySellAmount = "sell_amount"
)

// PromptListReq 获取提示词列表请求
type PromptListReq struct {
	Paginate      *utils.Page `json:"paginate"`
	ModelId       *int        `json:"modelId,omitempty"`
	CategoryTypes []int       `json:"categoryTypes,omitempty"`
	SortBy        string      `json:"sortBy"`
}

type PromptListResp struct {
	Paginate *utils.Page    `json:"paginate"`
	Prompts  []model.Prompt `json:"prompts"`
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
