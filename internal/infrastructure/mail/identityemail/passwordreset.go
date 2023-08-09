package identityemail

import (
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
)

type RequestPassordResetEmailService struct{}

func (r RequestPassordResetEmailService) SendMail(u user.User) error {
	//TODO send mail here
	return nil
}
