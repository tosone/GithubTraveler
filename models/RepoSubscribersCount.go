package models

import "github.com/jinzhu/gorm"

// RepoSubscribersCount ..
type RepoSubscribersCount struct {
	gorm.Model
	RepoID int64
	Num    int
}

// Upsert find or create new record
func (s *RepoSubscribersCount) Upsert() (err error) {
	if err = engine.Model(new(RepoSubscribersCount)).Where(RepoSubscribersCount{
		RepoID: s.RepoID,
		Num:    s.Num,
	}).First(s).Error; err == gorm.ErrRecordNotFound {
		err = engine.Create(s).Error
	}
	return
}
