package main

import (
	"promptrun-api/configs"
	"promptrun-api/model"
	"promptrun-api/routers"
)

func main() {
	configs.Init()
	defer model.CloseDB()

	router := routers.SetupRouter()

	if err := router.Run(":8080"); err != nil {
		return
	}
}
