package models

import "github.com/jinzhu/gorm"

// RepoNetworkCount ..
type RepoNetworkCount struct {
	gorm.Model
	RepoID int64
	Num    int
}

// Create find or create new record
func (s *RepoNetworkCount) Create() (err error) {
	if err = engine.Model(new(RepoNetworkCount)).Where(RepoNetworkCount{RepoID: s.RepoID, Num: s.Num}).First(s).Error; err == gorm.ErrRecordNotFound {
		if err = engine.Create(s).Error; err != nil {
			return
		}
	} else if err != nil {
		return
	}
	return
}
