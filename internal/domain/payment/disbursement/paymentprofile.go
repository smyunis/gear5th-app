package disbursement

import "gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"

type PaymentProfile struct {
	PaymentMethod string
	Account       string
	FullName      string
	PhoneNumber   user.PhoneNumber
}
