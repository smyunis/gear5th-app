package events

import (
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/identityinteractors"
)

type EventHandlerRegistrar struct {
	appEventDispatcher          application.EventDispatcher
	verificationEmailInteractor identityinteractors.VerificationEmailInteractor
}

func NewEventHandlerRegistrar(
	appEventDispatcher application.EventDispatcher,
	verificationEmailInteractor identityinteractors.VerificationEmailInteractor) EventHandlerRegistrar {
	return EventHandlerRegistrar{
		appEventDispatcher,
		verificationEmailInteractor,
	}
}

func (r *EventHandlerRegistrar) RegisterEventHandlers() error {

	r.appEventDispatcher.AddHandler("user/signedup", r.verificationEmailInteractor.HandleUserSignedUpEvent)

	return nil
}
