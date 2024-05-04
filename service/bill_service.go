package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"promptrun-api/common/errs"
	"promptrun-api/model"
	"promptrun-api/utils"
)

// BillListByUserIdReq 根据 userId 获取账单列表请求
type BillListByUserIdReq struct {
	Paginate *utils.Page `json:"paginate"`
	UserId   int         `json:"userId"`
}

type BillListRsp struct {
	BillList []model.Bill `json:"billList"`
	Rows     int          `json:"rows"`
}

func AddBill(c *gin.Context, bill model.Bill) (bool, *errs.Errs) {
	if err := model.DB.Create(&bill).Error; err != nil {
		utils.Log().Error(c.FullPath(), "创建账单失败，errMsg: %s", err.Error())
		return false, errs.NewErrs(errs.ErrDBError, errors.New("创建账单失败"))
	}
	return true, nil
}

func (r *BillListByUserIdReq) FindBillListByUserId(c *gin.Context) (BillListRsp, *errs.Errs) {
	var billList []model.Bill
	query := model.DB.Model(model.Bill{}).
		Where("user_id = ?", r.UserId).
		Count(&r.Paginate.Rows).
		Order("create_time DESC")
	if query.Scopes(utils.Paginate(r.Paginate)).Find(&billList).Error != nil {
		utils.Log().Error(c.FullPath(), "DB 获取账单列表失败")
		return BillListRsp{}, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取账单列表失败"))
	}

	return BillListRsp{
		BillList: billList,
		Rows:     r.Paginate.Rows,
	}, nil

}
