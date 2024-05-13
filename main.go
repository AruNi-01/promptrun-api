package main

import (
	"promptrun-api/configs"
	"promptrun-api/model"
	"promptrun-api/routers"
	"promptrun-api/service"
)

func main() {
	configs.Init()
	defer model.CloseDB()

	service.InitConsumer()

	router := routers.SetupRouter()

	if err := router.Run(":8080"); err != nil {
		return
	}
}
