package emailverification

import (
	"fmt"

	"gitlab.com/gear5th/gear5th-api/internal/application"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
)

func init() {
	application.ApplicationDomainEventDispatcher.AddHandler("user.created", SendEmailVerification)
}

func SendEmailVerification(event any) {
	event = event.(user.UserCreatedEvent)
	fmt.Printf("Sending Verification Email ... for %s\n", event)
}
