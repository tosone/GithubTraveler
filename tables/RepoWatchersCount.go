package tables

import "github.com/jinzhu/gorm"

// RepoWatchersCount ..
type RepoWatchersCount struct {
	gorm.Model
	RepoID int64
	Num    int
}

// Upsert find or create new record
func (s *RepoWatchersCount) Upsert() (err error) {
	if err = engine.Model(new(RepoWatchersCount)).Where(RepoWatchersCount{
		RepoID: s.RepoID,
		Num:    s.Num,
	}).First(s).Error; err == gorm.ErrRecordNotFound {
		err = engine.Create(s).Error
	}
	return
}
