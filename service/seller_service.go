package service

import (
	"errors"
	"promptrun-api/common/errs"
	"promptrun-api/model"
)

func FindSellerById(id int) (model.Seller, *errs.Errs) {
	var seller model.Seller
	if model.DB.First(&seller, id).RecordNotFound() {
		return model.Seller{}, errs.NewErrs(errs.ErrRecordNotFound, errors.New("未找到该卖家"))
	}
	return seller, nil
}
