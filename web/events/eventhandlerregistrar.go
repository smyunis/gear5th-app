package events

import (
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/adsinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/identityinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/paymentinteractors"
)

type EventHandlerRegistrar struct {
	appEventDispatcher          application.EventDispatcher
	adsPool                     adsinteractors.AdsPool
	adsInteractor               adsinteractors.AdsInteractor
	verificationEmailInteractor identityinteractors.VerificationEmailInteractor
	disbursementInteractor      paymentinteractors.DisbursementInteractor
	depositInteractor           paymentinteractors.DepositInteractor
	earningInteractor           paymentinteractors.EarningInteractor
}

func NewEventHandlerRegistrar(
	appEventDispatcher application.EventDispatcher,
	adsPool adsinteractors.AdsPool,
	adsInteractor adsinteractors.AdsInteractor,
	depositInteractor paymentinteractors.DepositInteractor,
	earningInteractor paymentinteractors.EarningInteractor,
	disbursementInteractor paymentinteractors.DisbursementInteractor,
	verificationEmailInteractor identityinteractors.VerificationEmailInteractor) EventHandlerRegistrar {
	return EventHandlerRegistrar{
		appEventDispatcher,
		adsPool,
		adsInteractor,
		verificationEmailInteractor,
		disbursementInteractor,
		depositInteractor,
		earningInteractor,
	}
}

func (r *EventHandlerRegistrar) RegisterEventHandlers() error {

	r.appEventDispatcher.AddHandler("user/signedup", r.verificationEmailInteractor.HandleUserSignedUpEvent)
	r.appEventDispatcher.AddHandler("campaign/adpieceadded", r.adsPool.OnNewAdPiece)
	r.appEventDispatcher.AddHandler("deposit/made", r.depositInteractor.OnNewDeposit)

	r.appEventDispatcher.AddHandler("impression/made", r.adsInteractor.IncrementImpressionCount)

	r.appEventDispatcher.AddHandler("disbursement/requested", r.disbursementInteractor.OnRequestDisbursement)

	return nil
}
