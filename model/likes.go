package model

import "time"

type Likes struct {
	Id           int       `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	UserId       int       `gorm:"column:user_id;type:int(11);comment:触发喜欢的用户 id" json:"user_id"`
	PromptId     int       `gorm:"column:prompt_id;type:int(11);comment:被喜欢的提示词 id" json:"prompt_id"`
	SellerUserId int       `gorm:"column:seller_user_id;type:int(11);comment:被喜欢的提示词的卖家用户 id" json:"seller_user_id"`
	CreateTime   time.Time `gorm:"column:create_time;type:datetime" json:"create_time"`
}
