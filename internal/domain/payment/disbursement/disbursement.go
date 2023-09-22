package withdrawalrequest

import (
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type DisbursementRepository interface {
	shared.EntityRepository[Disbursement]
	SettledDisbursementsForPublisher(publisherID shared.ID) ([]Disbursement, error)
	RequestedDisbursementsForPublisher(publisherID shared.ID) ([]Disbursement, error)
}

type DisbursementStatus = int

const (
	_ DisbursementStatus = iota
	Requested
	Settled
	Rejected
)

type Disbursement struct {
	ID                  shared.ID
	Events              shared.Events
	Status              DisbursementStatus
	PublisherID         shared.ID
	PaymentProfileID    shared.ID
	Amount              float64
	Time                time.Time
	SettlementReference string
}

func NewDisbursement(pubID shared.ID, paymentProfile shared.ID, amount float64) Disbursement {
	return Disbursement{
		ID:               shared.NewID(),
		Events:           make(shared.Events),
		Status:           Requested,
		PublisherID:      pubID,
		PaymentProfileID: paymentProfile,
		Amount:           amount,
		Time:             time.Now(),
	}
}

func (d *Disbursement) Settle(ref string) error {
	if d.Status == Rejected {
		return shared.ErrInvalidOperation
	}
	d.Status = Settled
	d.Time = time.Now()
	d.SettlementReference = ref
	d.Events.Emit("disbursement/settled", d)
	return nil
}

func (d *Disbursement) Reject() error {
	if d.Status == Settled {
		return shared.ErrInvalidOperation
	}
	d.Status = Rejected
	d.Time = time.Now()
	d.Events.Emit("disbursement/rejected", d)
	return nil
}
