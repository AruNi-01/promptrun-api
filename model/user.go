package model

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"promptrun-api/third_party"
	"promptrun-api/utils"
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

// BeforeUpdate 更新前的钩子函数，若更新了头像，则需要删除 OSS 中的原头像
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if u.HeaderUrl != "" {
		var oldUser User
		tx.First(&oldUser, u.Id)
		if oldUser.HeaderUrl != "" {
			utils.Log().Info("", "UserUpdate 更新了头像，执行删除 OSS 中的头像：%s", oldUser.HeaderUrl)
			err = third_party.DeleteOSSImgByUrl(oldUser.HeaderUrl)
		}
	}
	return
}
