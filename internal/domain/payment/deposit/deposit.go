package deposit

import (
	"math"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type DepositRepository interface {
	shared.EntityRepository[Deposit]
	DailyDisposits(day time.Time) ([]Deposit, error)
}

type Deposit struct {
	ID           shared.ID
	Events       shared.Events
	AdvertiserID shared.ID
	Amount       float64
	DepositTime  time.Time
	Start        time.Time
	End          time.Time
}

func NewDeposit(advertiserID shared.ID, amount float64, start time.Time, end time.Time) Deposit {
	if end.Before(start) {
		end = start.Add(168 * time.Hour) // 1 week
	}
	d := Deposit{
		ID:           shared.NewID(),
		Events:       make(shared.Events),
		AdvertiserID: advertiserID,
		DepositTime:  time.Now(),
		Start:        start,
		End:          end,
		Amount:       amount,
	}
	d.Events.Emit("deposit/made",d)
	return d
}

func (d *Deposit) DailyFundContributionAmount(day time.Time) float64 {
	if d.Start.After(day) || d.End.Before(day) {
		return 0
	}
	durationDays := math.Floor(d.End.Sub(d.Start).Hours() / 24)
	if durationDays == 0 {
		return d.Amount
	}

	return d.Amount / durationDays
}
