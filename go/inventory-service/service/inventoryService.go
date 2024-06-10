package service

import (
	"inventory-service/dto"
	"inventory-service/repository"
)

type InventoryService struct {
	InventoryRepo *repository.InventoryRepository
}

func NewInventoryService(inventoryRepo *repository.InventoryRepository) *InventoryService {
	return &InventoryService{
		InventoryRepo: inventoryRepo,
	}
}

func (service *InventoryService) IsInStock(skuCodes []string) ([]dto.InventoryResponse, error) {
	inventories, err := service.InventoryRepo.FindBySkuCodeIn(skuCodes)
	if err != nil {
		return nil, err
	}

	var responses []dto.InventoryResponse
	for _, inventory := range inventories {
		responses = append(responses, dto.InventoryResponse{
			SkuCode:  inventory.SkuCode,
			IsInStock: inventory.Quantity > 0,
		})
	}
	return responses, nil
}
