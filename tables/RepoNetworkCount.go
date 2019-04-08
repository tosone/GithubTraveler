package tables

import "github.com/jinzhu/gorm"

// RepoNetworkCount ..
type RepoNetworkCount struct {
	gorm.Model
	RepoID int64
	Num    int
}

// Upsert find or create new record
func (s *RepoNetworkCount) Upsert() (err error) {
	if err = engine.Model(new(RepoNetworkCount)).Where(RepoNetworkCount{
		RepoID: s.RepoID,
		Num:    s.Num,
	}).First(s).Error; err == gorm.ErrRecordNotFound {
		err = engine.Create(s).Error
	}
	return
}
