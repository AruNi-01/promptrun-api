package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"promptrun-api/common/errs"
	"promptrun-api/model"
	"promptrun-api/third_party"
	"promptrun-api/utils"
	"strconv"
	"time"
)

const (
	HeaderImgOSSPrefix = "header_img/"
)

type UserUpdateReq struct {
	UserId          int    `json:"userId"`
	Nickname        string `json:"nickname"`
	Email           string `json:"email"`
	HeaderImgBase64 string `json:"headerImgBase64,omitempty"`
}

func FindUserById(c *gin.Context, id int) (model.User, *errs.Errs) {
	var user model.User
	if model.DB.First(&user, id).RecordNotFound() {
		utils.Log().Error(c.FullPath(), "未找到该用户")
		return model.User{}, errs.NewErrs(errs.ErrUserNotExist, errors.New("用户不存在"))
	}
	return user, nil
}

func (r *UserUpdateReq) UpdateUser(c *gin.Context) (bool, *errs.Errs) {
	var user model.User

	user.Id = r.UserId
	user.Nickname = r.Nickname
	user.Email = r.Email

	// 头像 base64 图片上传到 OSS，返回图片地址
	if r.HeaderImgBase64 != "" {
		objectName := HeaderImgOSSPrefix + strconv.Itoa(r.UserId) + "-" + time.Now().Format("2006-01-02_150405")
		headerUrl, err := third_party.UploadBase64ImgToOSS(objectName, r.HeaderImgBase64)
		if err != nil {
			utils.Log().Error(c.FullPath(), "OSS 上传头像失败，errMsg: %s", err.Error())
			return false, errs.NewErrs(errs.ErrUploadImgToOSS, errors.New("OSS 上传头像失败"))
		}
		user.HeaderUrl = headerUrl
	}

	if err := model.DB.Model(&user).Updates(user).Error; err != nil {
		utils.Log().Error(c.FullPath(), "DB 更新用户信息失败，errMsg: %s", err.Error())
		return false, errs.NewErrs(errs.ErrDBError, errors.New("DB 更新用户信息失败"))
	}

	return true, nil
}
