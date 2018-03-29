package models

import "github.com/jinzhu/gorm"

// HistoryUserFollowingNum ..
type HistoryUserFollowingNum struct {
	gorm.Model
	UserID       uint64
	FollowingNum uint64
}

// Create ..
func (s *HistoryUserFollowingNum) Create() error {
	return engine.Create(s).Error
}
