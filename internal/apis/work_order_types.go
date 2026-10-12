package apis

type CreateWorkOrderRequest struct {
	QuotationID int    `json:"quotation_id"`
	CustomerID  int    `json:"customer_id"`
	VehicleID   int    `json:"vehicle_id"`
	DueDate     string `json:"due_date"`
}
