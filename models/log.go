package models

import "github.com/jinzhu/gorm"

type Log struct {
	gorm.Model
	Method   string
	Url      string
	Response []byte
	ErrMsg   []byte
}

func (log *Log) Create() (err error) {
	err = engine.Create(log).Error
	return
}
