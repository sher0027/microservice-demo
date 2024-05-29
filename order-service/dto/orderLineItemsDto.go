package dto

type OrderLineItemsDto struct {
	SkuCode  string  `json:"skuCode"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}
