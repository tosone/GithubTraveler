package tables

import "github.com/jinzhu/gorm"

type Topic struct {
	gorm.Model
	Content string `gorm:"not null;unique"`
}

// Upsert ..
func (topic *Topic) Upsert() (err error) {
	if err = engine.Model(new(Topic)).Where(Topic{Content: topic.Content}).First(topic).Error; err == gorm.ErrRecordNotFound {
		err = engine.Create(topic).Error
		return
	}
	return
}
