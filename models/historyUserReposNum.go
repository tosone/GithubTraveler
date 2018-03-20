package models

import "github.com/jinzhu/gorm"

type HistoryUserReposNum struct {
	gorm.Model
	UserID   uint64
	ReposNum uint64
}

func (s *HistoryUserReposNum) Create() error {
	return engine.Create(s).Error
}
