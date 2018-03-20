package models

import "github.com/jinzhu/gorm"

type HistoryUserGistNum struct {
	gorm.Model
	UserID  uint64
	GistNum uint64
}

func (s *HistoryUserGistNum) Create() error {
	return engine.Create(s).Error
}
