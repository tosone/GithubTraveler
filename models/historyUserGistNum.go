package models

import "github.com/jinzhu/gorm"

// HistoryUserGistNum ..
type HistoryUserGistNum struct {
	gorm.Model
	UserID  uint64
	GistNum uint64
}

// Create ..
func (s *HistoryUserGistNum) Create() error {
	return engine.Create(s).Error
}
