package model

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	CommonUserType = 0
	SellerUserType = 1
)

type User struct {
	Id         int       `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	Email      string    `gorm:"column:email;type:varchar(50)" json:"email"`
	Password   string    `gorm:"column:password;type:varchar(50)" json:"password"`
	Nickname   string    `gorm:"column:nickname;type:varchar(50)" json:"nickname"`
	HeaderUrl  string    `gorm:"column:header_url;type:varchar(255)" json:"header_url"`
	Type       int       `gorm:"column:type;type:int(11);comment:用户类型，0: 买家，1: 卖家。做个字段冗余，方便查询" json:"type"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime" json:"create_time"`
}

// SetPassword 设置密码（使用 bcrypt 加密）
func (u *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// CheckPassword 校验密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
