package models

import "github.com/jinzhu/gorm"

// IssueComments ..
type IssueComments struct {
	gorm.Model
	UserID uint64
	RepoID uint64
	Number uint64
	Body   []byte
}

// Create ..
func (s *IssueComments) Create() error {
	return engine.Create(s).Error
}
