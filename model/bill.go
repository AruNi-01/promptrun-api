package model

import "time"

const (
	BillTypeOutcome = iota
	BillTypeIncome
)

const (
	BillChannelWxPay = iota
	BillChannelAliPay
	BillChannelBalance
	BillChannelActivity
)

type Bill struct {
	Id         int       `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	UserId     int       `gorm:"column:user_id;type:int(11)" json:"user_id"`
	Type       int       `gorm:"column:type;type:smallint(6);comment:账单类型，0：支出，1：收入" json:"type"`
	Amount     float64   `gorm:"column:amount;type:decimal(10,2);comment:账单金额" json:"amount"`
	Channel    int       `gorm:"column:channel;type:smallint(6);comment:账单通道/渠道（支出/收入的渠道），0：微信，1：支付宝，2：余额、3：活动" json:"channel"`
	Remark     string    `gorm:"column:remark;type:varchar(255);comment:备注" json:"remark"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime" json:"create_time"`
}
