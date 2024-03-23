package service

import (
	"errors"
	"promptrun-api/common/errs"
	"promptrun-api/model"
)

func FindModelById(id int) (model.Model, *errs.Errs) {
	var promptModel model.Model
	if model.DB.First(&promptModel, id).RecordNotFound() {
		return model.Model{}, errs.NewErrs(errs.ErrRecordNotFound, errors.New("未找到该模型"))
	}
	return promptModel, nil
}
