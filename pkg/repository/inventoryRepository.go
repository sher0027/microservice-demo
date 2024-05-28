package repository

import (
    "go-microservice-demo/pkg/model"

    "github.com/jinzhu/gorm"
)

type InventoryRepository struct {
    DB *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
    return &InventoryRepository{DB: db}
}

func (repo *InventoryRepository) FindBySkuCodeIn(skuCodes []string) ([]model.Inventory, error) {
    var inventories []model.Inventory
    if err := repo.DB.Where("sku_code IN (?)", skuCodes).Find(&inventories).Error; err != nil {
        return nil, err
    }
    return inventories, nil
}