package models

import (
	"github.com/jinzhu/gorm"
)

// UserFollowers followers
type UserFollowers struct {
	gorm.Model
	UserID         uint64
	Version        string
	FollowerUserID uint64
}

// Create ..
func (f *UserFollowers) Create() (err error) {
	var isExist bool
	if isExist, err = f.IsExist(); err != nil {
		return
	} else if isExist {
		err = f.Update()
	} else {
		err = engine.Create(f).Error
	}
	return
}

// IsExist ..
func (f *UserFollowers) IsExist() (isExist bool, err error) {
	var count int
	if err = engine.Model(new(UserFollowers)).
		Where(UserFollowers{UserID: f.UserID, FollowerUserID: f.FollowerUserID}).
		Count(&count).Error; err != nil {
		return
	}
	if count != 0 {
		isExist = true
		return
	}
	return
}

// Update ..
func (f *UserFollowers) Update() (err error) {
	return
}
