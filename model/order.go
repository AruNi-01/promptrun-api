package model

import "time"

const (
	// OrderRatingNotYet 买家个人是否已进行评分，0：否
	OrderRatingNotYet = 0
	// OrderRatingDone 买家个人是否已进行评分，1：是
	OrderRatingDone = 1
)

type Order struct {
	Id         int64     `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT" json:"id"`
	PromptId   int       `gorm:"column:prompt_id;type:int(11);comment:购买的提示词 id" json:"prompt_id"`
	SellerId   int       `gorm:"column:seller_id;type:int(11);comment:提示词所属的卖家 id" json:"seller_id"`
	BuyerId    int       `gorm:"column:buyer_id;type:int(11);comment:购买提示词的买家 id（等于 user_id）" json:"buyer_id"`
	Price      float64   `gorm:"column:price;type:decimal(10,2);comment:买入的价格" json:"price"`
	IsRating   int       `gorm:"column:is_rating;type:int(11);comment:买家个人是否已进行评分，0：否，1：是" json:"is_rating"`
	Rating     float64   `gorm:"column:rating;type:float(2,1);comment:买家个人对该提示词的评分，1.0-5.0" json:"rating"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime" json:"create_time"`
}
