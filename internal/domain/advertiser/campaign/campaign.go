package campaign

import (
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type CampaignRepository interface {
	shared.EntityRepository[Campaign]
}

type Campaign struct {
	ID           shared.ID
	Events       shared.Events
	Name         string
	AdvertiserID shared.ID
	Start        time.Time
	End          time.Time
	RunDuration  time.Duration
}

func NewCampaign(name string, advertiserID shared.ID, start, end time.Time) Campaign {
	if end.Before(start) {
		end = start.Add(1 * 30 * 24 * time.Hour) // 1 month after start
	}
	return Campaign{
		ID:           shared.NewID(),
		Events:       make(shared.Events),
		Name:         name,
		AdvertiserID: advertiserID,
		Start:        start,
		End:          end,
		RunDuration:  end.Sub(start),
	}
}

func (c *Campaign) Quit() {
	c.RunDuration = time.Now().Sub(c.Start)
	c.Events.Emit("campaign/quitted", c)
}

func (c *Campaign) IsQuitted() bool {
	return c.End.Sub(c.Start) > c.RunDuration
}
