package models

import (
	"github.com/jinzhu/gorm"
)

type Repo struct {
	gorm.Model
	RepoID          uint
	Name            string
	StargazersCount int
}

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

func (repo *Repo) IsExist() (isExist bool, err error) {
	var count int
	if err = engine.Model(new(Repo)).Where(Repo{RepoID: repo.RepoID}).Count(&count).Error; err != nil {
		return
	}
	if count != 0 {
		isExist = true
		return
	}
	return
}

func (repo *Repo) Update() (err error) {
	err = engine.Model(new(Repo)).Where(Repo{RepoID: repo.RepoID}).Updates(repo).Error
	return
}
