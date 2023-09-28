package paymentinteractors

import (
	"context"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/domain/payment/disbursement"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/publisher"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type DisbursementEmailService interface {
	SendRequestDisbursementConfirmation(email user.Email, d disbursement.Disbursement) error
}

type DisbursementInteractor struct {
	disbursementRepository   disbursement.DisbursementRepository
	userRepository           user.UserRepository
	publisherRepository      publisher.PublisherRepository
	earningInteractor        EarningInteractor
	disbursementEmailService DisbursementEmailService
	cacheStore               application.KeyValueStore
	digiSignService          application.DigitalSignatureService
	eventDispatcher          application.EventDispatcher
	logger                   application.Logger
}

func NewDisbursementInteractor(
	disbursementRepository disbursement.DisbursementRepository,
	userRepository user.UserRepository,
	publisherRepository publisher.PublisherRepository,
	earningInteractor EarningInteractor,
	disbursementEmailService DisbursementEmailService,
	cacheStore application.KeyValueStore,
	digiSignService application.DigitalSignatureService,
	eventDispatcher application.EventDispatcher,
	logger application.Logger) DisbursementInteractor {
	return DisbursementInteractor{
		disbursementRepository,
		userRepository,
		publisherRepository,
		earningInteractor,
		disbursementEmailService,
		cacheStore,
		digiSignService,
		eventDispatcher,
		logger,
	}
}

func (i *DisbursementInteractor) RequestDisbursement(publisherID shared.ID,
	paymentProfile disbursement.PaymentProfile) error {

	pub, err := i.publisherRepository.Get(context.Background(), publisherID)
	if err != nil {
		return err
	}

	if i.earningInteractor.CanRequestDisbursement(publisherID) {
		return application.ErrRequirementFailed
	}

	currentBalance, err := i.earningInteractor.CurrentBalance(publisherID)
	if err != nil {
		return err
	}

	d := disbursement.NewDisbursement(publisherID,
		paymentProfile, currentBalance, pub.LastDisbursement, time.Now())

	err = i.disbursementRepository.Save(context.Background(), d)
	if err != nil {
		return err
	}

	i.eventDispatcher.DispatchAsync(d.Events)
	return nil
}

func (i *DisbursementInteractor) ConfirmDisbursement(disbursementID shared.ID, token string) error {

	if !i.digiSignService.Validate(token) {
		return application.ErrAuthorization
	}

	if m, _ := i.digiSignService.GetMessage(token); m != disbursementID.String() {
		return application.ErrAuthorization
	}

	d, err := i.disbursementRepository.Get(context.Background(), disbursementID)
	if err != nil {
		return err
	}
	err = d.Confirm()
	if err != nil {
		return err
	}

	err = i.disbursementRepository.Save(context.Background(), d)
	if err != nil {
		return err
	}

	i.eventDispatcher.DispatchAsync(d.Events)
	return nil
}

func (i *DisbursementInteractor) RejectDisbursement(disbursementID shared.ID, token string) error {

	if !i.digiSignService.Validate(token) {
		return application.ErrAuthorization
	}

	if m, _ := i.digiSignService.GetMessage(token); m != disbursementID.String() {
		return application.ErrAuthorization
	}

	d, err := i.disbursementRepository.Get(context.Background(), disbursementID)
	if err != nil {
		return err
	}
	err = d.Reject()
	if err != nil {
		return err
	}

	err = i.disbursementRepository.Save(context.Background(), d)
	if err != nil {
		return err
	}

	i.eventDispatcher.DispatchAsync(d.Events)
	return nil
}

func (i *DisbursementInteractor) OnRequestDisbursement(disb any) {
	d := disb.(disbursement.Disbursement)

	u, err := i.userRepository.Get(context.Background(), d.PublisherID)
	if err != nil {
		i.logger.Error("disbursement/get/user", err)
		return
	}

	err = i.disbursementEmailService.SendRequestDisbursementConfirmation(u.Email, d)
	if err != nil {
		i.logger.Error("disbursement/sendconfemail", err)
	}
}

func (i *DisbursementInteractor) DisbursementsForPublisher(publisherID shared.ID) ([]disbursement.Disbursement, error) {
	return i.disbursementRepository.DisbursementsForPublisher(publisherID)
}
