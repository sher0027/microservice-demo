package order

type InventoryResponse struct {
	SkuCode string `json:"skuCode"`
	InStock bool   `json:"inStock"`
}
