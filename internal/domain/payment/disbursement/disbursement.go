package disbursement

import (
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type DisbursementRepository interface {
	shared.EntityRepository[Disbursement]
	DisbursementsForPublisher(publisherID shared.ID, status DisbursementStatus) ([]Disbursement, error)
}

type DisbursementStatus = int

const (
	_ DisbursementStatus = iota
	Requested
	Confirmed
	Rejected
	Settled
)

type Disbursement struct {
	ID               shared.ID
	Events           shared.Events
	Status           DisbursementStatus
	PublisherID      shared.ID
	PaymentProfile   PaymentProfile
	Amount           float64
	Time             time.Time
	PeriodStart      time.Time
	PeriodEnd        time.Time
	SettlementRemark string
}

func NewDisbursement(pubID shared.ID, paymentProfile PaymentProfile, amount float64, start time.Time, end time.Time) Disbursement {
	d := Disbursement{
		ID:             shared.NewID(),
		Events:         make(shared.Events),
		Status:         Requested,
		PublisherID:    pubID,
		PaymentProfile: paymentProfile,
		Amount:         amount,
		Time:           time.Now(),
		PeriodStart:    start,
		PeriodEnd:      end,
	}
	d.Events.Emit("disbursement/requested", d)
	return d
}

func (d *Disbursement) Settle(settlementRemark string) error {
	if d.Status != Confirmed {
		return shared.ErrInvalidOperation
	}
	d.Status = Settled
	d.Time = time.Now()
	d.SettlementRemark = settlementRemark
	d.Events.Emit("disbursement/settled", d)
	return nil
}

func (d *Disbursement) Confirm() error {
	if d.Status != Requested {
		return shared.ErrInvalidOperation
	}
	d.Status = Confirmed
	d.Time = time.Now()
	d.Events.Emit("disbursement/confirmed", d)
	return nil
}

func (d *Disbursement) Reject() error {
	if d.Status != Requested {
		return shared.ErrInvalidOperation
	}
	d.Status = Rejected
	d.Time = time.Now()
	d.Events.Emit("disbursement/rejected", d)
	return nil
}
