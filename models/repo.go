package models

import "github.com/jinzhu/gorm"

type Repo struct {
	gorm.Model
	RepoID          uint
	Name            string
	StargazersCount int
}

func (repo *Repo) Create() (err error) {
	err = engine.Create(repo).Error
	return
}
