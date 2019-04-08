package models

import "github.com/jinzhu/gorm"

// UserFollowingCount followers
type UserFollowingCount struct {
	gorm.Model
	UserID int64
	Num    int
}

// Create find or create new record
func (s *UserFollowingCount) Upsert() (err error) {
	if err = engine.Model(new(UserFollowingCount)).Where(UserFollowingCount{
		UserID: s.UserID,
		Num:    s.Num,
	}).First(s).Error; err == gorm.ErrRecordNotFound {
		err = engine.Create(s).Error
	}
	return
}
