package model

// Img 是否为主图
const (
	PromptImgIsMaster  = 1
	PromptImgNotMaster = 0
)

type PromptImg struct {
	Id       int    `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	PromptId int    `gorm:"column:prompt_id;type:int(11);comment:提示词 id，逻辑关联到提示词表" json:"prompt_id"`
	ImgUrl   string `gorm:"column:img_url;type:varchar(255)" json:"img_url"`
	IsMaster int    `gorm:"column:is_master;type:int(11);comment:是否主图，0：非主图，1：主图" json:"is_master"`
}
