package models

import "github.com/jinzhu/gorm"

type Log struct {
	gorm.Model
	Type     string
	Method   string
	Url      string
	Response []byte `gorm:"size:65535"`
	ErrMsg   []byte
}

func (log *Log) Create() (err error) {
	err = engine.Create(log).Error
	return
}
