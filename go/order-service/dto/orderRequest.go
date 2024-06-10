package dto

type OrderRequest struct {
	OrderLineItems []OrderLineItemsDto `json:"orderLineItemsDtoList"`
}
