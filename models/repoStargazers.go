package models

import "github.com/jinzhu/gorm"

// RepoStargazers followers
type RepoStargazers struct {
	gorm.Model
	UserID  uint64
	Version string
	RepoID  uint64
}

// Create ..
func (f *RepoStargazers) Create() (err error) {
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
func (f *RepoStargazers) IsExist() (isExist bool, err error) {
	var count int
	if err = engine.Model(new(RepoStargazers)).
		Where(RepoStargazers{UserID: f.UserID, RepoID: f.RepoID}).
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
func (f *RepoStargazers) Update() (err error) {
	return
}
