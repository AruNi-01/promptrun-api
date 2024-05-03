package model

import (
	"errors"
	"fmt"
	"promptrun-api/utils"
	"time"
)

const (
	Money0  = 0.00
	Money5  = 5.00
	Money8  = 8.00
	Money10 = 10.00
	Money15 = 15.00
)

type Wallet struct {
	UserId        int       `gorm:"column:user_id;type:int(11);primary_key" json:"user_id"`
	WalletIncome  float64   `gorm:"column:wallet_income;type:decimal(10,2);comment:钱包总收入" json:"wallet_income"`
	WalletOutcome float64   `gorm:"column:wallet_outcome;type:decimal(10,2);comment:钱包总支出" json:"wallet_outcome"`
	Balance       float64   `gorm:"column:balance;type:decimal(10,2);comment:余额" json:"balance"`
	CreateTime    time.Time `gorm:"column:create_time;type:datetime" json:"create_time"`
}

func (w *Wallet) Create() error {
	if err := DB.Create(w).Error; err != nil {
		utils.Log().Error("", "创建钱包失败，errMsg: %s", err.Error())
		return errors.New(fmt.Sprintf("创建钱包失败, errMsg: %s", err.Error()))
	}
	return nil
}
