package apis

import (
	"testing"

	"garago/internal/database/postgres"
)

func TestPaymentSummary(t *testing.T) {
	invoice := postgres.Invoice{TotalAmount: 10000}
	cases := []struct {
		name       string
		payments   []postgres.Payment
		wantPaid   float64
		wantRemain float64
		wantStatus string
	}{
		{"unpaid", nil, 0, 10000, InvoiceUnpaid},
		{"partial", []postgres.Payment{{Amount: 2500}}, 2500, 7500, InvoicePartial},
		{"paid", []postgres.Payment{{Amount: 2500}, {Amount: 7500}}, 10000, 0, InvoicePaid},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			paid, remaining, status := paymentSummary(invoice.TotalAmount, tc.payments)
			if paid != tc.wantPaid || remaining != tc.wantRemain || status != tc.wantStatus {
				t.Fatalf("got paid=%v remaining=%v status=%s", paid, remaining, status)
			}
		})
	}
}

func TestValidPaymentMethod(t *testing.T) {
	if !validPaymentMethod(PaymentUPI) || validPaymentMethod("Cheque") {
		t.Fatal("payment method validation returned an unexpected result")
	}
}
