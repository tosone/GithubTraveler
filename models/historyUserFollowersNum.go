package models

import "github.com/jinzhu/gorm"

type HistoryUserFollowersNum struct {
	gorm.Model
	UserID       uint64
	FollowersNum uint64
}

func (s *HistoryUserFollowersNum) Create() error {
	return engine.Create(s).Error
}
