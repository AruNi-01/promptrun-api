package configs

import (
	"github.com/joho/godotenv"
	"os"
	"promptrun-api/model"
	"promptrun-api/utils"
)

// Init 从 .env 中加载配置初始化 app
func Init() {
	// 从本地 .env 文件中读取配置到 os 的环境变量中
	if err := godotenv.Load(); err != nil {
		utils.Log().Panic("", "load local env fail", err)
	}

	utils.BuildLogger(os.Getenv("LOG_LEVEL"))
	if err := model.InitDB(os.Getenv("MySQL_DSN")); err != nil {
		utils.Log().Panic("", "connect database fail", err)
	}
}
