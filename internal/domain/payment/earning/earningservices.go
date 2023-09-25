package earning

// TODO
const DisbursementRequestTreshold = 4500

func TotalEarningsAmount(earnings []Earning) float64 {
	sum := 0.0
	for _, e := range earnings {
		sum += e.Amount
	}
	return sum
}

func CanDisburseEarnings(earnings []Earning) bool {
	sum := 0.0
	for _, e := range earnings {
		sum += e.Amount
	}
	return sum > DisbursementRequestTreshold
}

func DailyRatePerImpression(totalDailyFund float64, totalImpressionCount int) float64 {
	//TODO set appropriate amonut for this
	const fixedRatePerImpression = 0.1

	if (float64(totalImpressionCount) * fixedRatePerImpression) <= totalDailyFund {
		return fixedRatePerImpression
	}
	infaltedRate := totalDailyFund / float64(totalImpressionCount)
	return infaltedRate
}

func PercentOfDisbursementTreshold(earning float64) float64 {
	p := (earning / DisbursementRequestTreshold) * 100
	if p > 100 {
		p = 100
	}
	return p
}
