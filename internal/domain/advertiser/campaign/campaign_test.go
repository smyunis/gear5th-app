package campaign_test

import (
	"slices"
	"testing"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/campaign"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

func TestQuitCampaign(t *testing.T) {
	c := newCampagin()
	campRun := 4 * 30 * 24 * time.Hour

	c.Quit()

	if c.RunDuration >= campRun {
		t.FailNow()
	}
}

func TestIsQuitted(t *testing.T) {
	c := newCampagin()

	c.Quit()

	if !c.IsQuitted() {
		t.FailNow()
	}
}

func TestIsNotQuitted(t *testing.T) {
	c := newCampagin()

	if c.IsQuitted() {
		t.FailNow()
	}
}

func TestQuitCampaignEvent(t *testing.T) {
	c := newCampagin()

	c.Quit()

	keys := make([]string, 0, len(c.Events))
	for k := range c.Events {
		keys = append(keys, k)
	}
	
	if !slices.Contains(keys,"campaign/quitted") {
		t.FailNow()
	}
}

func newCampagin() campaign.Campaign {
	start := time.Now()
	campRun := 4 * 30 * 24 * time.Hour
	end := start.Add(campRun) // 4 months
	return campaign.NewCampaign("coca-cola summer", shared.NewID(), start, end)
}
