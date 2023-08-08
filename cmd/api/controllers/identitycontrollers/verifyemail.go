package identitycontrollers

import "gitlab.com/gear5th/gear5th-api/cmd/api/controllers"

// Clickable link - GET method
// /identity/user/{userId}/verify-email?tmp-token={guid}
// - stores tmp-token in mem cache when email was sent
// - fetch now then verify if id did not expire
// - Redirect to a page that informs email was verifed, which inturn prompts to go to sign in page

type VerifyEmailController struct {
	controllers.Controller
}
