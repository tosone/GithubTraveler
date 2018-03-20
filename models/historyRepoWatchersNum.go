package models

import "github.com/jinzhu/gorm"

type HistoryRepoWatchersNum struct {
	gorm.Model
	UserID      uint64
	RepoID      uint64
	WatchersNum uint64
}

func (s *HistoryRepoWatchersNum) Create() error {
	return engine.Create(s).Error
}
