package paymentinteractors

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/payment/deposit"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type DepositInteractor struct {
	depositRepository deposit.DepositRepository
	cacheStore        application.KeyValueStore
	eventDispatcher   application.EventDispatcher
	logger            application.Logger
}

func NewDepositInteractor(
	depositRepository deposit.DepositRepository,
	cacheStore application.KeyValueStore,
	eventDispatcher application.EventDispatcher,
	logger application.Logger) DepositInteractor {
	return DepositInteractor{
		depositRepository,
		cacheStore,
		eventDispatcher,
		logger,
	}
}

func (i *DepositInteractor) AcceptDeposit(advertiserID shared.ID, amount float64, start, end time.Time) error {

	d := deposit.NewDeposit(advertiserID, amount, start, end)
	err := i.depositRepository.Save(context.Background(), d)
	if err != nil {
		return err
	}
	i.eventDispatcher.DispatchAsync(d.Events)
	return nil
}

func (i *DepositInteractor) OnNewDeposit(newDeposit any) {
	dep := newDeposit.(deposit.Deposit)

	err := i.updateTotalDailyFund(dep)
	if err != nil {
		return
	}
}

func (i *DepositInteractor) TotalDailyFund(day time.Time) (float64, error) {
	tf, err := i.cacheStore.Get(DailyDepositedFundCacheKey(day))
	totalFund, parseErr := strconv.ParseFloat(tf, 64)
	if err != nil || parseErr != nil {
		deposits, err := i.depositRepository.DailyDisposits(day)
		if err != nil {
			i.logger.Error("deposit/new/dailydeposits", err)
			return 0.0, err
		}
		totalFund = deposit.TotalDailyFund(day, deposits)
		fs := strconv.FormatFloat(totalFund, 'f', 2, 64)
		err = i.cacheStore.Save(DailyDepositedFundCacheKey(day), fs, 24*time.Hour)
		if err != nil {
			i.logger.Error("deposit/dailyfund/cachesave", err)
		}
	}
	return totalFund, nil
}

func (i *DepositInteractor) updateTotalDailyFund(dep deposit.Deposit) error {
	today := time.Now()
	tf, err := i.cacheStore.Get(DailyDepositedFundCacheKey(today))
	totalFund, parseErr := strconv.ParseFloat(tf, 64)
	if err != nil || parseErr != nil {
		deposits, err := i.depositRepository.DailyDisposits(today)
		if err != nil {
			i.logger.Error("deposit/new/dailydeposits", err)
			return err
		}
		totalFund = deposit.TotalDailyFund(today, deposits)

	} else {
		totalFund += dep.DailyFundContributionAmount(today)
	}

	fs := strconv.FormatFloat(totalFund, 'f', 2, 64)
	err = i.cacheStore.Save(DailyDepositedFundCacheKey(today), fs, 24*time.Hour)
	if err != nil {
		i.logger.Error("deposit/dailyfund/cachesave", err)
		return err
	}
	return nil
}

func DailyDepositedFundCacheKey(day time.Time) string {
	return fmt.Sprintf("dailyfund:%s", day.Format("20060102"))
}
