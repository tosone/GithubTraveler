package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	UserID string
	Login  string
	Type   string
	Raw    []byte
}
