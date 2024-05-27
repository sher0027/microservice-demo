package repository

import (
	"go-microservice-demo/pkg/model"

	"github.com/jinzhu/gorm"
)

type OrderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (repo *OrderRepository) Save(order *model.Order) error {
	return repo.DB.Save(order).Error
}
