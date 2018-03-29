package models

import "github.com/jinzhu/gorm"

// Log ..
type Log struct {
	gorm.Model
	Type     string
	Method   string
	URL      string
	Response []byte `gorm:"size:65535"`
	ErrMsg   []byte
}

// Create ..
func (log *Log) Create() (err error) {
	err = engine.Create(log).Error
	return
}
