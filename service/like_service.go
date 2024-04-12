package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"promptrun-api/common/errs"
	"promptrun-api/model"
	"promptrun-api/utils"
	"time"
)

var (
	LatestLikeTime time.Time // 最近一次点赞的时间，防止短时间内多次点赞
)

const (
	LikeInterval = 1 * time.Second // 点赞的时间间隔，1s
)

type LikeReq struct {
	UserId   int `json:"userId"`
	PromptId int `json:"promptId"`
	SellerId int `json:"sellerId"`
}

func IsLike(c *gin.Context, userId, promptId int) (bool, *errs.Errs) {
	var count int
	if model.DB.Model(model.Likes{}).Where("user_id = ? AND prompt_id = ?", userId, promptId).Count(&count).Error != nil {
		utils.Log().Error(c.FullPath(), "DB 查询是否点赞错误")
		return false, errs.NewErrs(errs.ErrDBError, errors.New("DB 查询是否点赞错误"))
	}
	return count > 0, nil
}

func (r *LikeReq) Like(c *gin.Context) (err *errs.Errs) {
	if LatestLikeTime.Add(LikeInterval).After(time.Now()) {
		utils.Log().Warning(c.FullPath(), "点赞时间间隔过短")
		return errs.NewErrs(errs.ErrLikeIntervalTooShort, errors.New("点赞时间间隔过短"))
	}
	LatestLikeTime = time.Now()

	var count, addOrSubAmount int
	if model.DB.Model(model.Likes{}).Where("user_id = ? AND prompt_id = ?", r.UserId, r.PromptId).Count(&count).Error != nil {
		utils.Log().Error(c.FullPath(), "DB 查询是否点赞错误")
		return errs.NewErrs(errs.ErrDBError, errors.New("DB 查询是否点赞错误"))
	}

	// 已经点赞则取消点赞
	if count > 0 {
		if model.DB.Where("user_id = ? AND prompt_id = ?", r.UserId, r.PromptId).Delete(&model.Likes{}).Error != nil {
			utils.Log().Error(c.FullPath(), "取消点赞失败")
			return errs.NewErrs(errs.ErrDBError, errors.New("取消点赞失败"))
		}
		addOrSubAmount = -1
	} else { // 未点赞则点赞
		likes := model.Likes{
			UserId:     r.UserId,
			PromptId:   r.PromptId,
			SellerId:   r.SellerId,
			CreateTime: time.Now(),
		}
		if err := model.DB.Create(&likes).Error; err != nil {
			utils.Log().Error(c.FullPath(), "点赞失败")
			return errs.NewErrs(errs.ErrDBError, errors.New("点赞失败"))
		}
		addOrSubAmount = 1
	}

	if model.DB.Model(&model.Prompt{}).Where("id = ?", r.PromptId).UpdateColumn("like_amount", gorm.Expr("like_amount + ?", addOrSubAmount)).Error != nil {
		utils.Log().Error(c.FullPath(), "更新点赞数失败")
		return errs.NewErrs(errs.ErrDBError, errors.New("更新点赞数失败"))
	}
	if model.DB.Model(&model.Seller{}).Where("id = ?", r.SellerId).UpdateColumn("like_amount", gorm.Expr("like_amount + ?", addOrSubAmount)).Error != nil {
		utils.Log().Error(c.FullPath(), "更新卖家获点赞数失败")
		return errs.NewErrs(errs.ErrDBError, errors.New("更新卖家获点赞数失败"))
	}
	return nil

}
