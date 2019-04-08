package tables

import "github.com/jinzhu/gorm"

// RepoStargazersCount ..
type RepoStargazersCount struct {
	gorm.Model
	RepoID int64
	Num    int
}

// Upsert find or create new record
func (s *RepoStargazersCount) Upsert() (err error) {
	if err = engine.Model(new(RepoStargazersCount)).Where(RepoStargazersCount{
		RepoID: s.RepoID,
		Num:    s.Num,
	}).First(s).Error; err == gorm.ErrRecordNotFound {
		err = engine.Create(s).Error
	}
	return
}
