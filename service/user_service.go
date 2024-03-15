package service

import "promptrun-api/model"

func FindUserById(id int) model.User {
	var user model.User
	model.DB.First(&user, id)
	return user
}
