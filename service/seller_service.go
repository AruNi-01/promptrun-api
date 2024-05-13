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

type SellerUpdateReq struct {
	UserId          int    `json:"userId"`
	Nickname        string `json:"nickname"`
	Email           string `json:"email"`
	HeaderImgBase64 string `json:"headerImgBase64,omitempty"`

	Intro string `json:"intro"`
}

func FindSellerById(c *gin.Context, id int) (model.Seller, *errs.Errs) {
	var seller model.Seller
	if model.DB.First(&seller, id).RecordNotFound() {
		utils.Log().Error(c.FullPath(), "未找到该卖家")
		return model.Seller{}, errs.NewErrs(errs.ErrRecordNotFound, errors.New("未找到该卖家"))
	}
	return seller, nil
}

func FindSellerByUserId(c *gin.Context, userId int) (model.Seller, *errs.Errs) {
	var seller model.Seller
	if model.DB.Where("user_id = ?", userId).First(&seller).RecordNotFound() {
		utils.Log().Error(c.FullPath(), "未找到该卖家")
		return model.Seller{}, errs.NewErrs(errs.ErrRecordNotFound, errors.New("未找到该卖家"))
	}
	return seller, nil
}

func (r *SellerUpdateReq) UpdateSeller(c *gin.Context) (bool, *errs.Errs) {
	updateUserReq := UserUpdateReq{
		UserId:          r.UserId,
		Nickname:        r.Nickname,
		Email:           r.Email,
		HeaderImgBase64: r.HeaderImgBase64,
	}
	flag, e := updateUserReq.UpdateUser(c)
	if !flag || e != nil {
		return false, e
	}

	seller, e := FindSellerByUserId(c, r.UserId)
	if e != nil {
		return false, e
	}

	if err := model.DB.Model(&seller).Where("id = ?", seller.Id).Update("intro", r.Intro).Error; err != nil {
		utils.Log().Error(c.FullPath(), "DB 更新卖家信息失败，errMsg: %s", err.Error())
		return false, errs.NewErrs(errs.ErrDBError, errors.New("DB 更新卖家信息失败"))
	}
	return true, nil
}

func FindBecomeSellerDayBySellerId(c *gin.Context, sellerId int) (int, *errs.Errs) {
	var seller model.Seller
	if model.DB.First(&seller, sellerId).RecordNotFound() {
		utils.Log().Error(c.FullPath(), "未找到该卖家")
		return 0, errs.NewErrs(errs.ErrRecordNotFound, errors.New("未找到该卖家"))
	}
	return int(time.Now().Sub(seller.CreateTime).Hours() / 24), nil
}

func IncreaseSellerSellAmount(sellerId int) (bool, *errs.Errs) {
	if err := model.DB.Model(&model.Seller{}).Where("id = ?", sellerId).UpdateColumn("sell_amount", gorm.Expr("sell_amount + ?", 1)).Error; err != nil {
		utils.Log().Error("", "DB 更新卖家销售量失败，errMsg: %s", err.Error())
		return false, errs.NewErrs(errs.ErrDBError, errors.New("DB 更新卖家销售量失败"))
	}
	return true, nil
}
