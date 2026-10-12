package main

type QuotationItemRequest struct {
	ServiceID int `json:"service_id"`
	Quantity  int `json:"quantity"`
}

type CreateQuotationRequest struct {
	CustomerID int                    `json:"customer_id"`
	VehicleID  int                    `json:"vehicle_id"`
	AdvanceRequested float64          `json:"advance_requested"`
	Items      []QuotationItemRequest `json:"items"`
}
