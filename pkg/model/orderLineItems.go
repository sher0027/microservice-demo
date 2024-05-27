package model

import (
	"github.com/jinzhu/gorm"
)

type OrderLineItems struct {
	gorm.Model
	OrderID  uint `gorm:"index"`
	SkuCode  string
	Price    float64
	Quantity int
}
