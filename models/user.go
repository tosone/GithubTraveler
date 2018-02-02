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
	isEmpty = true
	return
}

func (user *User) Create() (err error) {
	var isExist bool
	if isExist, err = user.IsExist(); err != nil {
		return
	} else if isExist {
		err = user.Update()
	} else {
		err = engine.Create(user).Error
	}
	return
}

func (user *User) FindByID(id uint) (u *User, err error) {
	u = new(User)
	err = engine.Find(u, id).Error
	return
}

func (user *User) IsExist() (isExist bool, err error) {
	var count int
	if err = engine.Model(new(User)).Where(User{Login: user.Login}).Count(&count).Error; err != nil {
		return
	}
	if count != 0 {
		isExist = true
		return
	}
	return
}

func (user *User) Update() (err error) {
	err = engine.Model(new(Repo)).Where(User{Login: user.Login}).Updates(user).Error
	return
}
