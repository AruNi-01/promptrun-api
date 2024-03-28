package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"promptrun-api/utils"
	"time"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(dsn string) {
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		utils.Log().Panic("", "connect database fail, errMsg: %s", err.Error())
		panic(err)
	}

	// 日志模式: 设置为 true 打印详细日志，默认只打印 error 日志
	db.LogMode(true)
	// 表名影射时使用单数
	db.SingularTable(true)
	// 连接池设置
	db.DB().SetMaxIdleConns(30)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(time.Minute * 5)

	// 测试连通性
	if err := db.DB().Ping(); err != nil {
		utils.Log().Panic("", "ping database fail, errMsg: %s", err.Error())
		panic(err)
	}

	DB = db
}

func CloseDB() {
	DB.Close()
}
