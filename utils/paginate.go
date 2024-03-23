package utils

import "github.com/jinzhu/gorm"

type Page struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Rows     int `json:"rows,omitempty"`
}

// Paginate 分页
func Paginate(req *Page) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if req.Page == 0 {
			req.Page = 1
		}
		if req.PageSize <= 0 {
			req.PageSize = 10
		} else if req.PageSize > 100 {
			req.PageSize = 100
		}
		offset := (req.Page - 1) * req.PageSize
		return db.Offset(offset).Limit(req.PageSize)
	}
}
