package paymentprofile

import (
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type PaymentProfileRepository interface {
	shared.EntityRepository[PaymentProfile]
}

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
	ID            shared.ID
	UserID        shared.ID
	PaymentMethod PaymentMethod
	Account       string
	FullName      string
	PhoneNumber   string
}

func NewPaymentProfile(userID shared.ID, paymentMethod PaymentMethod, account, fullName, phoneNumber string) PaymentProfile {
	return PaymentProfile{
		ID:            shared.NewID(),
		UserID:        userID,
		PaymentMethod: paymentMethod,
		Account:       account,
		FullName:      fullName,
		PhoneNumber:   phoneNumber,
	}
}
