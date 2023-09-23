package disbursement

import "gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"

type PaymentMethod int

const (
	_ PaymentMethod = iota
	PhysicalCollection
	Cheque
	CommercialBankOfEthiopia
	BankOfAbyssinia
	Telebirr
)

type PaymentProfile struct {
	PaymentMethod PaymentMethod
	Account       string
	FullName      string
	PhoneNumber   user.PhoneNumber
}
