package app

type CreatePaymentRequest struct {
	InvoiceID     int     `json:"invoice_id"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
}
