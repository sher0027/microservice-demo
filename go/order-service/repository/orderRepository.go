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

func (repo *OrderRepository) SaveOrderWithItems(order *model.Order, items []model.OrderLineItems) error {
	tx := repo.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	log.Printf("Saving order: %+v", order)
	if err := tx.Table("t_orders").Save(order).Error; err != nil {
		log.Printf("Error saving order: %v", err)
		tx.Rollback()
		return err
	}

	for i := range items {
		items[i].OrderId = order.Id
	}

	for _, item := range items {
		log.Printf("Saving order line item: %+v", item)
		if err := tx.Table("t_order_line_items").Save(&item).Error; err != nil {
			log.Printf("Error saving order line item: %v", err)
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction: %v", err)
		tx.Rollback()
		return err
	}

	log.Printf("Order and items saved successfully")
	return nil
}