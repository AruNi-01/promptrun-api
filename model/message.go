package model

import "time"

const (
	MessageTypeActivity = 1
	MessageTypeSell     = 2
	MessageTypeAudit    = 3
	MessageTypeLike     = 4
	MessageTypeWithdraw = 5
)

const (
	MessageIsReadNo  = 0
	MessageIsReadYes = 1
)

type Message struct {
	Id         int       `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	FromUserId int       `gorm:"column:from_user_id;type:int(11)" json:"from_user_id"`
	ToUserId   int       `gorm:"column:to_user_id;type:int(11)" json:"to_user_id"`
	Type       int       `gorm:"column:type;type:smallint(6);comment:消息类型：1-活动，2-售出，3-审核，4-点赞，5-提现" json:"type"`
	Content    string    `gorm:"column:content;type:text" json:"content"`
	IsRead     int       `gorm:"column:is_read;type:smallint(6);comment:是否已读：0-未读，1-已读" json:"is_read"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime" json:"create_time"`
}
