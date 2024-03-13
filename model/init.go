package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(dsn string) (err error) {
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		return
	}

	// 日志模式: 设置为 true 打印详细日志，默认只打印 error 日志
	DB.LogMode(true)
	// 表名影射时使用单数
	DB.SingularTable(true)
	// 连接池设置
	DB.DB().SetMaxIdleConns(30)
	DB.DB().SetMaxOpenConns(100)
	DB.DB().SetConnMaxLifetime(time.Minute * 5)

	// 测试连通性再返回，ping 得通返回 nil，否则返回 error
	return DB.DB().Ping()
}

func CloseDB() {
	DB.Close()
}
