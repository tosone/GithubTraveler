package tables

import "github.com/jinzhu/gorm"

var engine *gorm.DB

func Initialize(db *gorm.DB) {
	engine = db
}
