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
		utils.Log().Error(c.FullPath(), "未找到该用户的钱包")
		return model.Wallet{}, errs.NewErrs(errs.ErrRecordNotFound, errors.New("未找到该用户的钱包"))
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

// CalculateBalanceAndIncome 计算余额和收入，余额增加，收入增加
func CalculateBalanceAndIncome(c *gin.Context, userId int, increaseBalance float64, increaseIncome float64) (bool, *errs.Errs) {
	wallet, err := FindWalletByUserId(c, userId)
	if err != nil {
		return false, err
	}
	wallet.Balance += increaseBalance
	wallet.WalletIncome += increaseIncome
	if err := model.DB.Save(&wallet).Error; err != nil {
		utils.Log().Error(c.FullPath(), "更新余额失败，errMsg: %s", err.Error())
		return false, errs.NewErrs(errs.ErrDBError, errors.New("更新余额失败"))
	}
	return true, nil
}

// CalculateBalanceAndOutcome 计算余额和支出，余额减少，支出增加
func CalculateBalanceAndOutcome(c *gin.Context, userId int, decreaseBalance float64, increaseOutcome float64) (bool, *errs.Errs) {
	wallet, err := FindWalletByUserId(c, userId)
	if err != nil {
		return false, err
	}
	wallet.Balance -= decreaseBalance
	wallet.WalletOutcome += increaseOutcome
	if err := model.DB.Save(&wallet).Error; err != nil {
		utils.Log().Error(c.FullPath(), "更新余额失败，errMsg: %s", err.Error())
		return false, errs.NewErrs(errs.ErrDBError, errors.New("更新余额失败"))
	}
	return true, nil
}
