package dto

type InventoryResponse struct {
	SkuCode string `json:"skuCode"`
	IsInStock bool   `json:"inStock"`
}
