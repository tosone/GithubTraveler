package tables

import "github.com/jinzhu/gorm"

// RepoForksCount ..
type RepoForksCount struct {
	gorm.Model
	RepoID int64
	Num    int
}

// Upsert find or create new record
func (s *RepoForksCount) Upsert() (err error) {
	if err = engine.Model(new(RepoForksCount)).Where(RepoForksCount{
		RepoID: s.RepoID,
		Num:    s.Num,
	}).First(s).Error; err == gorm.ErrRecordNotFound {
		err = engine.Create(s).Error
	}
	return
}
