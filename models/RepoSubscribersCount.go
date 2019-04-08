package models

import "github.com/jinzhu/gorm"

// RepoSubscribersCount ..
type RepoSubscribersCount struct {
	gorm.Model
	RepoID int64
	Num    int
}

// Create find or create new record
func (s *RepoSubscribersCount) Create() (err error) {
	if err = engine.Model(new(RepoSubscribersCount)).Where(RepoSubscribersCount{RepoID: s.RepoID, Num: s.Num}).First(s).Error; err == gorm.ErrRecordNotFound {
		if err = engine.Create(s).Error; err != nil {
			return
		}
	} else if err != nil {
		return
	}
	return
}
