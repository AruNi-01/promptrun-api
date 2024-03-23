package model

// ModelMediaType 模型媒体类型
const (
	ModelMediaTypeText  = 0
	ModelMediaTypeImage = 1
	ModelMediaTypeVideo = 2
)

type Model struct {
	Id        int    `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	Name      string `gorm:"column:name;type:varchar(50);comment:模型名称" json:"name"`
	MediaType int    `gorm:"column:media_type;type:int(11);comment:媒体类型，0：文本、1：图片、2：视频" json:"media_type"`
}
