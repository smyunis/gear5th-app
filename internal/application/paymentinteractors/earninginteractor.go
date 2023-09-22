package paymentinteractors

import (
	"context"
	"strconv"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/adsinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/domain/ads/impression"
	"gitlab.com/gear5th/gear5th-app/internal/domain/payment/deposit"
	"gitlab.com/gear5th/gear5th-app/internal/domain/payment/earning"
)

type EarningInteractor struct {
	earningRepository    earning.EarningRepository
	depositRepository    deposit.DepositRepository
	impressionRepository impression.ImpressionRepository
	cacheStore           application.KeyValueStore
	eventDispatcher      application.EventDispatcher
	logger               application.Logger
}

func NewEarningInteractor(
	earningRepository earning.EarningRepository,
	depositRepository deposit.DepositRepository,
	impressionRepository impression.ImpressionRepository,
	cacheStore application.KeyValueStore,
	eventDispatcher application.EventDispatcher,
	logger application.Logger) EarningInteractor {
	return EarningInteractor{
		earningRepository,
		depositRepository,
		impressionRepository,
		cacheStore,
		eventDispatcher,
		logger,
	}
}

func (i *EarningInteractor) OnImpression(newImpression any) {
	imp := newImpression.(impression.Impression)

	totalFund, err := i.totalDailyFund()
	if err != nil {
		return
	}

	totalImpressions, err := i.totalImpressionCount()
	if err != nil {
		return
	}

	dailyRate := earning.DailyRatePerImpression(totalFund, totalImpressions)
	impressionEarning := earning.NewEarning(imp.OriginPublisherID, earning.Impression, dailyRate, imp.AdPieceID, imp.OriginAdSlotID, imp.OriginSiteID)
	err = i.earningRepository.Save(context.Background(), impressionEarning)
	if err != nil {
		i.logger.Error("earning/save", err)
		return
	}

	i.eventDispatcher.DispatchAsync(impressionEarning.Events)
}

func (i *EarningInteractor) totalImpressionCount() (int, error) {
	today := time.Now()
	c, err := i.cacheStore.Get(adsinteractors.DailyImpressionCountCacheKey(today))
	totalImpressions, parseErr := strconv.Atoi(c)
	if err != nil || parseErr != nil {
		totalImpressions, err = i.impressionRepository.DailyImpressionCount(today)
		if err != nil {
			i.logger.Error("impressions/dailyimpressioncount", err)
			return 0, err
		}
		cs := strconv.Itoa(totalImpressions)
		i.cacheStore.Save(adsinteractors.DailyImpressionCountCacheKey(today), cs, 24*time.Hour)
	}
	return totalImpressions, nil
}

func (i *EarningInteractor) totalDailyFund() (float64, error) {
	today := time.Now()
	tf, err := i.cacheStore.Get(DailyDepositedFundCacheKey(today))
	totalFund, parseErr := strconv.ParseFloat(tf, 64)
	if err != nil || parseErr != nil {
		deposits, err := i.depositRepository.DailyDisposits(today)
		if err != nil {
			i.logger.Error("deposit/new/dailydeposits", err)
			return 0.0, err
		}
		totalFund = deposit.TotalDailyFund(today, deposits)
		fs := strconv.FormatFloat(totalFund, 'f', 2, 64)
		err = i.cacheStore.Save(DailyDepositedFundCacheKey(today), fs, 24*time.Hour)
		if err != nil {
			i.logger.Error("deposit/dailyfund/cachesave", err)
		}
	}
	return totalFund, nil
}
