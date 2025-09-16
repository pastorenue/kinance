package common

type Currency string

const (
	EUR Currency = "EUR"
	USD Currency = "USD"
	GBP Currency = "GBP"
	JPY Currency = "JPY"
	CHF Currency = "CHF"
	NGN Currency = "NGN"
	// Add more currencies as needed
)

type PaymentMethod string

const (
	Cash         PaymentMethod = "cash"
	Card         PaymentMethod = "card"
	BankTransfer PaymentMethod = "bank_transfer"
)
