package model

type Inventory struct {
	Id	uint	`gorm:"primary_key"`
	SkuCode  string `gorm:"column:sku_code;not null"`
	Quantity int    `gorm:"not null"`
}