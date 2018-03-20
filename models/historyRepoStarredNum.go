package models

import "github.com/jinzhu/gorm"

type HistoryRepoStarredNum struct {
	gorm.Model
	UserID     uint64
	RepoID     uint64
	StarredNum uint64
}

func (s *HistoryRepoStarredNum) Create() error {
	return engine.Create(s).Error
}
