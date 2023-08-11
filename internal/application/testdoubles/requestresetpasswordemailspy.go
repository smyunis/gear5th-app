package testdoubles

import (
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
)

// This is not how spies are meant to be. This is a very bad implmentation

var (
	sendMailCalled  bool
	sendMailArgUser user.User
)

type RequestResetPasswordEmailSpy struct{}

func (r RequestResetPasswordEmailSpy) SendMail(u user.User, tok string) error {
	sendMailCalled = true
	sendMailArgUser = u
	return nil
}

func RequestResetPasswordEmailSpyGet() (user.User, bool) {
	return sendMailArgUser, sendMailCalled
}

func RequestResetPasswordEmailSpyReset() {
	sendMailArgUser = user.User{}
	sendMailCalled = false
}

type RequestResetPasswordEmailStub struct{}

func (r RequestResetPasswordEmailStub) SendMail(u user.User,token string) error {
	return nil
}
