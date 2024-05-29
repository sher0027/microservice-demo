package model

type Order struct {
	Id	uint	`gorm:"primary_key;auto_increment"`
	OrderNumber	string	`gorm:"column:order_number;unique_index"`
}
