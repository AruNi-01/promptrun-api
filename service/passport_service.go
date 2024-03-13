package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"promptrun-api/common/errs"
	"promptrun-api/model"
	"promptrun-api/utils"
	"time"
)

type LoginReq struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type RegisterReq struct {
	LoginReq
	ConfirmPassword string `form:"confirmPassword" json:"confirmPassword" binding:"required"`
}

// Register 用户注册
func (r *RegisterReq) Register(c *gin.Context) (model.User, *errs.Errs) {
	// 1. 表单数据校验
	if err := r.valid(); err != nil {
		return model.User{}, err
	}
	var user = model.User{
		Email:      r.Email,
		Nickname:   r.Email,
		Type:       model.CommonUserType,
		CreateTime: time.Now(),
	}
	// 2. 加密密码
	if err := user.SetPassword(r.Password); err != nil {
		utils.Log().Error(c.FullPath(), "密码加密错误")
		return model.User{}, errs.NewErrs(errs.ErrEncryptError, errors.New("密码加密失败"))
	}
	// 3. 创建用户
	if err := model.DB.Create(&user).Error; err != nil {
		utils.Log().Error(c.FullPath(), "DB 创建用户失败")
		return model.User{}, errs.NewErrs(errs.ErrDBError, errors.New("DB 创建用户失败"))
	}

	return user, nil
}

// Login 用户登录
func (r *LoginReq) Login(c *gin.Context) (model.User, *errs.Errs) {
	var user model.User
	if model.DB.Where("email = ?", r.Email).First(&user).Error != nil {
		return model.User{}, errs.NewErrs(errs.ErrUserNotExist, errors.New("账号不存在"))
	}
	if !user.CheckPassword(r.Password) {
		return model.User{}, errs.NewErrs(errs.ErrWrongPassword, errors.New("密码错误"))
	}
	return user, nil
}

// EmailIsExist 邮箱是否存在
func EmailIsExist(email string) *errs.Errs {
	count := 0
	model.DB.Model(&model.User{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		return errs.NewErrs(errs.ErrEmailExist, errors.New("邮箱已注册"))
	}
	return nil
}

func (r *RegisterReq) valid() *errs.Errs {
	if r.Password != r.ConfirmPassword {
		return errs.NewErrs(errs.ErrConfirmPasswordDiff, errors.New("两次输入的密码不相同"))
	}
	if err := EmailIsExist(r.Email); err != nil {
		return err
	}
	return nil
}
