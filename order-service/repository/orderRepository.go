package repository

import (
	"log"
	"order-service/model"

	"github.com/jinzhu/gorm"
)

type OrderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (repo *OrderRepository) Save(order *model.Order) error {
	log.Printf("Try save: %+v", order)
	return repo.DB.Table("t_orders").Save(order).Error
}
