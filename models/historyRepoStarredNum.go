package models

import "github.com/jinzhu/gorm"

// HistoryRepoStarredNum ..
type HistoryRepoStarredNum struct {
	gorm.Model
	UserID     uint64
	RepoID     uint64
	StarredNum uint64
}

// Create ..
func (s *HistoryRepoStarredNum) Create() error {
	return engine.Create(s).Error
}
