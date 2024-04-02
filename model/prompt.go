package model

import "time"

// Prompt 审核状态
const (
	AuditStatusNotAudit = 0 // 未审核（审核失败）
	AuditStatusAuditing = 1 // 审核中
	AuditStatusPass     = 2 // 审核通过
)

// Prompt 上下架状态
const (
	PublishStatusOff = 0 // 下架（未发布）
	PublishStatusOn  = 1 // 上架（已发布）
)

type Prompt struct {
	Id            int       `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	SellerId      int       `gorm:"column:seller_id;type:int(11);comment:卖家 id，逻辑关联到卖家表" json:"seller_id"`
	ModelId       int       `gorm:"column:model_id;type:int(11);comment:模型 id，逻辑关联模型表。" json:"model_id"`
	CategoryType  int       `gorm:"column:category_type;type:int(11);comment:分类类型，具体枚举看文档" json:"category_type"`
	Title         string    `gorm:"column:title;type:varchar(255);comment:提示词标题" json:"title"`
	Intro         string    `gorm:"column:intro;type:varchar(255);comment:提示词介绍" json:"intro"`
	Price         float64   `gorm:"column:price;type:decimal(10,2);comment:价格" json:"price"`
	Rating        float64   `gorm:"column:rating;type:float(2,1);comment:评分，1.0-5.0" json:"rating"`
	Score         float64   `gorm:"column:score;type:double;comment:分数，热度排行使用" json:"score"`
	SellAmount    int       `gorm:"column:sell_amount;type:int(11);comment:销量" json:"sell_amount"`
	BrowseAmount  int       `gorm:"column:browse_amount;type:int(11);comment:浏览数量" json:"browse_amount"`
	LikeAmount    int       `gorm:"column:like_amount;type:int(11);comment:提示词被喜欢数量" json:"like_amount"`
	PublishStatus int       `gorm:"column:publish_status;type:int(11);comment:上架状态，0: 下架，1：上架" json:"publish_status"`
	AuditStatus   int       `gorm:"column:audit_status;type:int(11);comment:审核状态，0：审核失败，1：审核中，2：审核通过" json:"audit_status"`
	CreateTime    time.Time `gorm:"column:create_time;type:datetime" json:"create_time"`
}
