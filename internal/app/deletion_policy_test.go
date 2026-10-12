package app

import "testing"

func TestCustomerHistoryBlocksDeletion(t *testing.T) {
	tests := []struct {
		name    string
		history customerHistory
		blocked bool
	}{
		{name: "no history", history: customerHistory{}, blocked: false},
		{name: "vehicle", history: customerHistory{vehicles: true}, blocked: true},
		{name: "quotation", history: customerHistory{quotations: true}, blocked: true},
		{name: "work order", history: customerHistory{workOrders: true}, blocked: true},
		{name: "invoice", history: customerHistory{invoices: true}, blocked: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.history.blocksDeletion(); got != test.blocked {
				t.Fatalf("blocksDeletion() = %v, want %v", got, test.blocked)
			}
		})
	}
}
