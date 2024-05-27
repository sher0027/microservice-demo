package model

import (
	"github.com/jinzhu/gorm"
)

type Order struct {
	gorm.Model
	OrderNumber    string           `gorm:"unique_index"`
	OrderLineItems []OrderLineItems `gorm:"foreignkey:OrderID"`
}
