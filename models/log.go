package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

// Log ..
type Log struct {
	gorm.Model
	Type     string
	Mark     string
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

func (log *Log) GetOutDate() bool {
	var l = new(Log)
	err := engine.Find(l, Log{Type: log.Type, Mark: log.Mark})
	if err != nil {
		return true
	}
	if time.Since(l.CreatedAt).Hours() > float64(viper.GetInt("Crawler.ItemOutDate")) {
		return true
	}
	return false
}
