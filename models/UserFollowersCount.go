package models

import (
	"github.com/jinzhu/gorm"
)

// UserFollowersCount followers
type UserFollowersCount struct {
	gorm.Model
	UserID int64
	Num    int
}

// Upsert find or create new record
func (s *UserFollowersCount) Upsert() (err error) {
	if err = engine.Model(new(UserFollowersCount)).Where(UserFollowersCount{
		UserID: s.UserID,
		Num:    s.Num,
	}).First(s).Error; err == gorm.ErrRecordNotFound {
		err = engine.Create(s).Error
	}
	return
}
