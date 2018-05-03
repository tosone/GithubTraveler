package models

import "github.com/jinzhu/gorm"

// RepoIssues issues
type RepoIssues struct {
	gorm.Model
	UserID   uint64
	RepoID   uint64
	Number   uint64
	Comments uint64
	Title    string
	Body     []byte
}

// Create ..
func (f *RepoIssues) Create() (err error) {
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

// FindByID ..
func (f *RepoIssues) FindByID(id uint) (u *RepoIssues, err error) {
	u = new(RepoIssues)
	err = engine.Find(u, id).Error
	return
}

// IsExist ..
func (f *RepoIssues) IsExist() (isExist bool, err error) {
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
func (f *RepoIssues) Update() (err error) {
	return
}
