package repository

import (
    "inventory-service/model"

    "github.com/jinzhu/gorm"
)

type InventoryRepository struct {
    DB *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
    return &InventoryRepository{DB: db}
}

func (repo *InventoryRepository) FindBySkuCodeIn(skuCodes []string) ([]model.Inventory, error) {
	var inventory []model.Inventory
	if err := repo.DB.Table("t_inventory").Where("sku_code IN (?)", skuCodes).Find(&inventory).Error; err != nil {
		return nil, err
	}
	return inventory, nil
}