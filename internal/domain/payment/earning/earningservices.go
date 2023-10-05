package earning


// TODO
const DisbursementRequestTreshold = 2
const FixedRatePerImpression = 0.1


func TotalEarningsAmount(totalDailyFund float64, totalImpressionCount int, impressionCount int) float64 {
	rpi := DailyRatePerImpression(totalDailyFund, totalImpressionCount) 
	return rpi * float64(impressionCount)
}

func CanDisburseEarnings(currentBalance float64) bool {
	return currentBalance > DisbursementRequestTreshold
}

func DailyRatePerImpression(totalDailyFund float64, totalImpressionCount int) float64 {
	//TODO set appropriate amonut for this

	if (float64(totalImpressionCount) * FixedRatePerImpression) <= totalDailyFund {
		return FixedRatePerImpression
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
