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
	// earningRepository    earning.EarningRepository
	depositRepository    deposit.DepositRepository
	publisherRepository  publisher.PublisherRepository
	siteRepository       site.SiteRepository
	impressionRepository impression.ImpressionRepository
	cacheStore           application.KeyValueStore
	eventDispatcher      application.EventDispatcher
	logger               application.Logger
}

func NewEarningInteractor(
	// earningRepository earning.EarningRepository,
	depositRepository deposit.DepositRepository,
	publisherRepository publisher.PublisherRepository,
	siteRepository site.SiteRepository,
	impressionRepository impression.ImpressionRepository,
	cacheStore application.KeyValueStore,
	eventDispatcher application.EventDispatcher,
	logger application.Logger) EarningInteractor {
	return EarningInteractor{
		// earningRepository,
		depositRepository,
		publisherRepository,
		siteRepository,
		impressionRepository,
		cacheStore,
		eventDispatcher,
		logger,
	}
}

// TODO when settling settle from beg of day
func (i *EarningInteractor) CurrentBalance(publisherID shared.ID) (float64, error) {
	p, err := i.publisherRepository.Get(context.Background(), publisherID)
	if err != nil {
		return 0.0, err
	}
	balanceCacheKey := fmt.Sprintf("balance:%s:%s", publisherID.String(), time.Now().Format("20060102"))
	b, err := i.cacheStore.Get(balanceCacheKey)
	bal, parseErr := strconv.ParseFloat(b, 64)
	if err != nil || parseErr != nil {

		lastDisbursement, _ := shared.TimeEdges(p.LastDisbursement, time.Now())
		today, _ := shared.TimeEdges(time.Now(), time.Now())

		bal, err = i.Earnings(publisherID, lastDisbursement, today)
		if err != nil {
			return 0.0, err
		}

		fs := strconv.FormatFloat(bal, 'f', 2, 64)
		i.cacheStore.Save(balanceCacheKey, fs, 168*time.Hour)
	}
	return bal, nil
}

func (i *EarningInteractor) Earnings(publisherID shared.ID, start time.Time, end time.Time) (float64, error) {
	earningCacheKey := fmt.Sprintf("earning:%s:%s-%s", publisherID.String(), start.Format("20060102"), end.Format("20060102"))
	b, err := i.cacheStore.Get(earningCacheKey)
	earn, parseErr := strconv.ParseFloat(b, 64)
	if err != nil || parseErr != nil {

		// s, e := shared.TimeEdges(start, end)

		sum := 0.0

		for in := 0; start.AddDate(0, 0, in).Before(end); in++ {
			day := start.AddDate(0, 0, in)
			dailyImpCount, err := i.impressionRepository.ImpressionsCountForPublisher(publisherID,
				day, day.AddDate(0, 0, 1))
			if err != nil {
				return 0.0, err
			}
			dailyFund, err := i.totalDailyFund(day)
			if err != nil {
				return 0.0, err
			}
			totalImpCount, err := i.totalImpressionCount(day)
			if err != nil {
				return 0.0, err
			}
			sum += earning.TotalEarningsAmount(dailyFund, totalImpCount, dailyImpCount)
		}
		earn = sum
		fs := strconv.FormatFloat(earn, 'f', 2, 64)
		i.cacheStore.Save(earningCacheKey, fs, 168*time.Hour)
	}
	return earn, nil
}

func (i *EarningInteractor) CanRequestDisbursement(publisherID shared.ID) bool {
	bal, err := i.CurrentBalance(publisherID)
	if err != nil {
		return false
	}
	return earning.CanDisburseEarnings(bal)
}

func (i *EarningInteractor) totalImpressionCount(day time.Time) (int, error) {
	c, err := i.cacheStore.Get(adsinteractors.DailyImpressionCountCacheKey(day))
	totalImpressions, parseErr := strconv.Atoi(c)
	if err != nil || parseErr != nil {
		totalImpressions, err = i.impressionRepository.DailyImpressionCount(day)
		if err != nil {
			i.logger.Error("impressions/dailyimpressioncount", err)
			return 0, err
		}
		cs := strconv.Itoa(totalImpressions)
		i.cacheStore.Save(adsinteractors.DailyImpressionCountCacheKey(day), cs, 168*time.Hour)
	}
	return totalImpressions, nil
}

func (i *EarningInteractor) totalDailyFund(day time.Time) (float64, error) {
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
		err = i.cacheStore.Save(DailyDepositedFundCacheKey(day), fs, 168*time.Hour)
		if err != nil {
			i.logger.Error("deposit/dailyfund/cachesave", err)
		}
	}
	return totalFund, nil
}
