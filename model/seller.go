package model

import "time"

// SellerStatus 卖家状态
const (
	SellerStatusDisable = 0
	SellerStatusEnable  = 1
)

type Seller struct {
	Id         int       `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	UserId     int       `gorm:"column:user_id;type:int(11);comment:逻辑外键，关联到用户 id;NOT NULL" json:"user_id"`
	Rating     float64   `gorm:"column:rating;type:float(2,1);comment:卖家总评分，1.0-5.0" json:"rating"`
	Status     int       `gorm:"column:status;type:int(11);comment:卖家状态，0: 禁用，1: 启用" json:"status"`
	Intro      string    `gorm:"column:intro;type:varchar(255);comment:卖家简介" json:"intro"`
	SellAmount int       `gorm:"column:sell_amount;type:int(11);comment:销量" json:"sell_amount"`
	LikeAmount int       `gorm:"column:like_amount;type:int(11);comment:提示词被喜欢数量" json:"like_amount"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime;comment:成为卖家时间" json:"create_time"`
}
