package service

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"os"
	"promptrun-api/cache"
	"promptrun-api/common/errs"
	"promptrun-api/model"
	"promptrun-api/utils"
	"time"
)

// LoginReq 登录请求
type LoginReq struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// RegisterReq 注册请求
type RegisterReq struct {
	LoginReq
	ConfirmPassword string `form:"confirmPassword" json:"confirmPassword" binding:"required"`
}

// LoginTicket 登录凭证
type LoginTicket struct {
	UserId    int       `json:"userId"`
	Ticket    string    `json:"ticket"`
	ExpiredAt time.Time `json:"expiredAt"`
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

	// 生成登录凭证
	loginTicket := BuildLoginTicket(user.Id)
	// 将登录凭证存入 Redis
	jsonTicket, err := json.Marshal(loginTicket)
	if err != nil {
		utils.Log().Error(c.FullPath(), "json convert error", err)
		return model.User{}, errs.NewErrs(errs.ErrJsonConvertError, err)
	}
	ticketKey := cache.TicketKey(loginTicket.Ticket)
	cache.RedisCli.Set(c, ticketKey, jsonTicket, time.Until(loginTicket.ExpiredAt))

	// 设置 Cookie，后续请求携带 Cookie，实现登录状态保持
	saveCookie(c, loginTicket.Ticket, loginTicket.ExpiredAt)

	return user, nil
}

func saveCookie(c *gin.Context, ticket string, expiredAt time.Time) {
	c.SetCookie("ticket",
		ticket,
		int(time.Until(expiredAt).Seconds()),
		"/",
		os.Getenv("COOKIE_DOMAIN"),
		false,
		false,
	)
}

func emailIsExist(email string) bool {
	count := 0
	model.DB.Model(&model.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

func (r *RegisterReq) valid() *errs.Errs {
	if emailIsExist(r.Email) {
		return errs.NewErrs(errs.ErrEmailExist, errors.New("此邮箱已经注册"))
	}
	if r.Password != r.ConfirmPassword {
		return errs.NewErrs(errs.ErrConfirmPasswordDiff, errors.New("两次输入的密码不相同"))
	}
	return nil
}

// BuildLoginTicket 创建登录凭证
func BuildLoginTicket(userId int) LoginTicket {
	return LoginTicket{
		UserId:    userId,
		Ticket:    utils.GenUUID(),
		ExpiredAt: time.Now().Add(24 * time.Hour),
	}
}

// FindLoginTicket 根据 ticket 查找登录凭证
func FindLoginTicket(c *gin.Context, ticket string) (LoginTicket, *errs.Errs) {
	jsonTicket := cache.RedisCli.Get(c, cache.TicketKey(ticket)).Val()

	var loginTicket LoginTicket
	err := json.Unmarshal([]byte(jsonTicket), &loginTicket)
	if err != nil {
		utils.Log().Error(c.FullPath(), "json convert error", err)
		return LoginTicket{}, errs.NewErrs(errs.ErrJsonConvertError, err)
	}

	return loginTicket, nil
}
