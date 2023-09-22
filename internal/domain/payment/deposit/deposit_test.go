package deposit_test

import (
	"testing"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/domain/payment/deposit"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

func TestDailyFuncContribution(t *testing.T) {

	st := time.Now()
	en := st.Add(2 * time.Hour)
	d := deposit.NewDeposit(shared.NewID(),4500, st, en)

	df := d.DailyFundContributionAmount(st.Add(1 * time.Hour))
	if df == 0 {
		t.Fatal(df)
	}
	t.Log(df)
}
