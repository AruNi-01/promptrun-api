package model

import "time"

type PromptDetail struct {
	Id            int       `gorm:"column:id;type:int(11);primary_key" json:"id"`
	PromptId      int       `gorm:"column:prompt_id;type:int(11)" json:"prompt_id"`
	Content       string    `gorm:"column:content;type:text;comment:Prompt 具体内容" json:"content"`
	MediaType     int       `gorm:"column:media_type;type:int(11)" json:"media_type"`
	UseSuggestion string    `gorm:"column:use_suggestion;type:text;comment:使用建议" json:"use_suggestion"`
	CreateTime    time.Time `gorm:"column:create_time;type:datetime" json:"create_time"`
}
