package deposit

import "time"

func TotalDailyFund(day time.Time, deposits []Deposit) float64 {
	sum := 0.0
	for _, d := range deposits {
		if d.Start.Before(day) && d.End.After(day) {
			sum += d.DailyFundContributionAmount(day)
		}
	}
	return sum
}
