package models

import "github.com/jinzhu/gorm"

type HistoryUserFollowingNum struct {
	gorm.Model
	UserID       uint64
	FollowingNum uint64
}

func (s *HistoryUserFollowingNum) Create() error {
	return engine.Create(s).Error
}
