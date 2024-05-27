package order

type OrderRequest struct {
	OrderLineItems []OrderLineItemsDto `json:"orderLineItems"`
}
