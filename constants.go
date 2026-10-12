package main

const (
	// Work Order statuses
	WorkOrderPending    = "Pending"
	WorkOrderInProgress = "In Progress"
	WorkOrderCompleted  = "Completed"
	WorkOrderCancelled  = "Cancelled"

	// Quotation statuses
	QuotationDraft    = "Draft"
	QuotationApproved = "Approved"
	QuotationRejected = "Rejected"

	// Invoice statuses
	InvoiceUnpaid  = "Unpaid"
	InvoicePartial = "Partial"
	InvoicePaid    = "Paid"

	// Payment methods

	PaymentCash         = "Cash"
	PaymentUPI          = "UPI"
	PaymentCard         = "Card"
	PaymentBankTransfer = "Bank Transfer"
)
