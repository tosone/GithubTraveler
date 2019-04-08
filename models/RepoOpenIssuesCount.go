package models

import "github.com/jinzhu/gorm"

// RepoOpenIssuesCount ..
type RepoOpenIssuesCount struct {
	gorm.Model
	RepoID int64
	Num    int
}

// Upsert find or create new record
func (s *RepoOpenIssuesCount) Upsert() (err error) {
	if err = engine.Model(new(RepoOpenIssuesCount)).Where(RepoOpenIssuesCount{
		RepoID: s.RepoID,
		Num:    s.Num,
	}).First(s).Error; err == gorm.ErrRecordNotFound {
		err = engine.Create(s).Error
	}
	return
}
