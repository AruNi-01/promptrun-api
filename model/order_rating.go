package model

import "time"

type OrderRating struct {
	Id         int       `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	OrderId    int       `gorm:"column:order_id;type:int(11)" json:"order_id"`
	PromptId   int       `gorm:"column:prompt_id;type:int(11)" json:"prompt_id"`
	SellerId   int       `gorm:"column:seller_id;type:int(11)" json:"seller_id"`
	Rating     float64   `gorm:"column:rating;type:float(2,1)" json:"rating"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime" json:"create_time"`
}
