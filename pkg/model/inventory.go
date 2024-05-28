package model

import (
	"github.com/jinzhu/gorm"
)

type Inventory struct {
	gorm.Model
	SkuCode  string `gorm:"unique_index;not null"`
	Quantity int    `gorm:"not null"`
}