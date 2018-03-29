package models

import (
	"github.com/jinzhu/gorm"
)

// Repo ..
type Repo struct {
	gorm.Model
	UserID      uint64
	RepoID      uint64
	Name        string
	Homepage    string
	Language    string
	Size        uint64
	Licence     string
	Description string
	Stargazers  string
	Watchers    string
}

// IsEmpty ..
func (repo *Repo) IsEmpty() (isEmpty bool, err error) {
	var count int
	if err = engine.Model(new(Repo)).Count(&count).Error; err != nil {
		return
	}
	if count != 0 {
		isEmpty = false
	}
	isEmpty = true
	return
}

// Create ..
func (repo *Repo) Create() (err error) {
	var isExist bool
	if isExist, err = repo.IsExist(); err != nil {
		return
	} else if isExist {
		err = repo.Update()
	} else {
		err = engine.Create(repo).Error
	}
	return
}

// FindByID ..
func (repo *Repo) FindByID(id uint) (u *Repo, err error) {
	u = new(Repo)
	err = engine.Find(u, id).Error
	return
}

// IsExist ..
func (repo *Repo) IsExist() (isExist bool, err error) {
	var count int
	if err = engine.Model(new(Repo)).Where(Repo{UserID: repo.UserID, RepoID: repo.RepoID}).Count(&count).Error; err != nil {
		return
	}
	if count != 0 {
		isExist = true
		return
	}
	return
}

// Update ..
func (repo *Repo) Update() (err error) {
	return engine.Model(new(Repo)).Where(Repo{UserID: repo.UserID, RepoID: repo.RepoID}).Updates(repo).Error
}
