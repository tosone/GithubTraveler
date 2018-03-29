package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// User ..
type User struct {
	gorm.Model
	UserID    uint64
	Login     string
	Type      string
	Location  string
	Email     string
	Followers string
	Following string
}

// IsEmpty ..
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

// Create ..
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

// FindByID ..
func (user *User) FindByID(id uint) (u *User, err error) {
	u = new(User)
	err = engine.Find(u, id).Error
	return
}

// FindByUserID ..
func (user *User) FindByUserID(id uint64) (u *User, err error) {
	if id == 0 {
		err = errors.New("cannot find record with user ID 0")
		return
	}
	u = new(User)
	err = engine.Find(u, User{UserID: id}).Error
	return
}

// IsExist ..
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

// Update ..
func (user *User) Update() (err error) {
	return engine.Model(new(User)).Where(User{Login: user.Login}).Updates(user).Error
}
