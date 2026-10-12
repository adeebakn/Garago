package app

type CreateInvoiceRequest struct {
	WorkOrderID int `json:"work_order_id"`
}

//Work Order already tells us which quotation it belongs to, and the quotation already has the total amount and customer.
