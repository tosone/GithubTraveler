package models

import "github.com/jinzhu/gorm"

// HistoryRepoIssuesNum ..
type HistoryRepoIssuesNum struct {
	gorm.Model
	UserID    uint64
	RepoID    uint64
	IssuesNum uint64
}

// Create ..
func (s *HistoryRepoIssuesNum) Create() error {
	return engine.Create(s).Error
}
