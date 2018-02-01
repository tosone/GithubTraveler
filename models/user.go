package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	UserID string
	Login  string
	Type   string
}

func (user *User) IsEmpty() (isEmpty bool, err error) {
	var count int
	if err = engine.Model(new(User)).Count(&count).Error; err != nil {
		return
	}
	if count != 0 {
		isEmpty = false
	}
	return
}

func (user *User) Create() (err error) {
	err = engine.Create(user).Error
	return
}

func (user *User) FindByID(id uint) (u *User, err error) {
	u = new(User)
	err = engine.Find(u, id).Error
	return
}
