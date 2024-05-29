package model

type OrderLineItems struct {
	Id	uint `gorm:"primary_key;auto_increment"`
	SkuCode  string
	Price    float64
	Quantity int
	OrderId  uint`gorm:"column:order_id;not null"`
}
