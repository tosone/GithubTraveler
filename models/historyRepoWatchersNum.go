package models

import "github.com/jinzhu/gorm"

// HistoryRepoWatchersNum ..
type HistoryRepoWatchersNum struct {
	gorm.Model
	UserID      uint64
	RepoID      uint64
	WatchersNum uint64
}

// Create ..
func (s *HistoryRepoWatchersNum) Create() error {
	return engine.Create(s).Error
}
