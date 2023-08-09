package identityemail

import (
	"fmt"

	"gitlab.com/gear5th/gear5th-api/internal/application"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
)

func init() {
	application.ApplicationEventDispatcher.AddHandler("user.signedup", sendEmailVerification)
}

func sendEmailVerification(event any) {
	signedUpUser := event.(user.UserCreatedEvent)
	if !signedUpUser.IsEmailVerified {
		//TODO send email with link to verify email
		fmt.Printf("Sending Verification Email ... for %s \n", signedUpUser.Email.String())
	}

}
