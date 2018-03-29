package models

import "github.com/jinzhu/gorm"

// HistoryUserReposNum ..
type HistoryUserReposNum struct {
	gorm.Model
	UserID   uint64
	ReposNum uint64
}

// Create ..
func (s *HistoryUserReposNum) Create() error {
	return engine.Create(s).Error
}
