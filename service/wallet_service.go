package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"promptrun-api/common/errs"
	"promptrun-api/model"
	"promptrun-api/utils"
	"time"
)

func FindWalletByUserId(c *gin.Context, userId int) (wallet model.Wallet, err *errs.Errs) {
	if model.DB.Where("user_id = ?", userId).First(&wallet).RecordNotFound() {
		wallet, err = CreateWallet(c, userId)
	}
	return
}

func CreateWallet(c *gin.Context, userId int) (model.Wallet, *errs.Errs) {
	wallet := model.Wallet{
		UserId:        userId,
		WalletIncome:  model.Money0,
		WalletOutcome: model.Money0,
		Balance:       model.Money10,
		CreateTime:    time.Now(),
	}
	if err := wallet.Create(); err != nil {
		utils.Log().Error(c.FullPath(), "创建钱包失败，errMsg: %s", err.Error())
		return model.Wallet{}, errs.NewErrs(errs.ErrDBError, errors.New("创建钱包失败"))
	}
	return wallet, nil
}
