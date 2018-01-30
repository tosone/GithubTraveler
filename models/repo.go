package models

import "github.com/jinzhu/gorm"

type Repo struct {
	gorm.Model
	RepoID          string
	Name            string
	StargazersCount int
	Raw             []byte
}
