package events

import (
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/adsinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/identityinteractors"
)

type EventHandlerRegistrar struct {
	appEventDispatcher          application.EventDispatcher
	adsPool                     adsinteractors.AdsPool
	verificationEmailInteractor identityinteractors.VerificationEmailInteractor
}

func NewEventHandlerRegistrar(
	appEventDispatcher application.EventDispatcher,
	adsPool adsinteractors.AdsPool,
	verificationEmailInteractor identityinteractors.VerificationEmailInteractor) EventHandlerRegistrar {
	return EventHandlerRegistrar{
		appEventDispatcher,
		adsPool,
		verificationEmailInteractor,
	}
}

func (r *EventHandlerRegistrar) RegisterEventHandlers() error {

	r.appEventDispatcher.AddHandler("user/signedup", r.verificationEmailInteractor.HandleUserSignedUpEvent)
	r.appEventDispatcher.AddHandler("campaign/adpieceadded", r.adsPool.OnNewAdPiece)

	return nil
}
