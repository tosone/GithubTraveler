package models

import "github.com/jinzhu/gorm"

// HistoryUserFollowersNum ..
type HistoryUserFollowersNum struct {
	gorm.Model
	UserID       uint64
	FollowersNum uint64
}

// Create ..
func (s *HistoryUserFollowersNum) Create() error {
	return engine.Create(s).Error
}
