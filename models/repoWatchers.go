package models

import "github.com/jinzhu/gorm"

// RepoWatchers followers
type RepoWatchers struct {
	gorm.Model
	UserID  uint64
	Version string
	RepoID  uint64
}

// Create ..
func (f *RepoWatchers) Create() (err error) {
	var isExist bool
	if isExist, err = f.IsExist(); err != nil {
		return
	} else if isExist {
		err = f.Update()
	} else {
		err = engine.Create(f).Error
	}
	return
}

// IsExist ..
func (f *RepoWatchers) IsExist() (isExist bool, err error) {
	var count int
	if err = engine.Model(new(RepoWatchers)).
		Where(RepoWatchers{UserID: f.UserID, RepoID: f.RepoID}).
		Count(&count).Error; err != nil {
		return
	}
	if count != 0 {
		isExist = true
		return
	}
	return
}

// Update ..
func (f *RepoWatchers) Update() (err error) {
	return
}
