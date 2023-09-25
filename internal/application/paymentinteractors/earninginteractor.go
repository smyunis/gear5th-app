package paymentinteractors

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/adsinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/domain/ads/impression"
	"gitlab.com/gear5th/gear5th-app/internal/domain/payment/deposit"
	"gitlab.com/gear5th/gear5th-app/internal/domain/payment/earning"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/publisher"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/site"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type EarningInteractor struct {
	earningRepository    earning.EarningRepository
	depositRepository    deposit.DepositRepository
	publisherRepository  publisher.PublisherRepository
	siteRepository       site.SiteRepository
	impressionRepository impression.ImpressionRepository
	cacheStore           application.KeyValueStore
	eventDispatcher      application.EventDispatcher
	logger               application.Logger
}

func NewEarningInteractor(
	earningRepository earning.EarningRepository,
	depositRepository deposit.DepositRepository,
	publisherRepository publisher.PublisherRepository,
	siteRepository site.SiteRepository,
	impressionRepository impression.ImpressionRepository,
	cacheStore application.KeyValueStore,
	eventDispatcher application.EventDispatcher,
	logger application.Logger) EarningInteractor {
	return EarningInteractor{
		earningRepository,
		depositRepository,
		publisherRepository,
		siteRepository,
		impressionRepository,
		cacheStore,
		eventDispatcher,
		logger,
	}
}

func (i *EarningInteractor) CurrentBalance(publisherID shared.ID) (float64, error) {
	p, err := i.publisherRepository.Get(context.Background(), publisherID)
	if err != nil {
		return 0.0, err
	}
	balanceCacheKey := fmt.Sprintf("balance:%s:%s", publisherID.String(), time.Now().Format("20060102"))
	b, err := i.cacheStore.Get(balanceCacheKey)
	bal, parseErr := strconv.ParseFloat(b, 64)
	if err != nil || parseErr != nil {
		earnings, err := i.earningRepository.EarningsForPublisher(publisherID, p.LastDisbursement, time.Now())
		if err != nil {
			return 0.0, err
		}
		bal = earning.TotalEarningsAmount(earnings)
		fs := strconv.FormatFloat(bal, 'f', 2, 64)
		i.cacheStore.Save(balanceCacheKey, fs, 2*time.Hour)
	}
	return bal, nil
}

func (i *EarningInteractor) Earnings(publisherID shared.ID, start time.Time, end time.Time) (float64, error) {
	earningCacheKey := fmt.Sprintf("earning:%s:%s-%s", publisherID.String(), start.Format("20060102"), end.Format("20060102"))
	b, err := i.cacheStore.Get(earningCacheKey)
	e, parseErr := strconv.ParseFloat(b, 64)
	if err != nil || parseErr != nil {
		earnings, err := i.earningRepository.EarningsForPublisher(publisherID, start, end)
		if err != nil {
			return 0.0, err
		}
		e = earning.TotalEarningsAmount(earnings)
		fs := strconv.FormatFloat(e, 'f', 2, 64)
		i.cacheStore.Save(earningCacheKey, fs, 24*time.Hour)
	}
	return e, nil
}

func (i *EarningInteractor) CanRequestDisbursement(publisherID shared.ID) bool {
	bal, err := i.CurrentBalance(publisherID)
	if err != nil {
		return false
	}
	return bal > earning.DisbursementRequestTreshold
}

func (i *EarningInteractor) OnImpression(newImpression any) {
	imp := newImpression.(impression.Impression)

	if !i.canSiteMonetize(imp.OriginSiteID) {
		return
	}

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

func (i *EarningInteractor) canSiteMonetize(siteID shared.ID) bool {
	siteCanMonetizeCacheKey := fmt.Sprintf("site:%s:canMonetize", siteID.String())
	sm, err := i.cacheStore.Get(siteCanMonetizeCacheKey)
	siteCanMonetize, parseErr := strconv.ParseBool(sm)
	if err != nil || parseErr != nil {
		s, err := i.siteRepository.Get(context.Background(), siteID)
		if err != nil {
			return false
		}
		siteCanMonetize = s.CanMonetize()
		i.cacheStore.Save(siteCanMonetizeCacheKey, strconv.FormatBool(siteCanMonetize), 12*time.Hour)
	}
	return siteCanMonetize
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
