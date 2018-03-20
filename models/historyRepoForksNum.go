package models

import "github.com/jinzhu/gorm"

type HistoryRepoForksNum struct {
	gorm.Model
	UserID   uint64
	RepoID   uint64
	ForksNum uint64
}

func (s *HistoryRepoForksNum) Create() error {
	return engine.Create(s).Error
}
