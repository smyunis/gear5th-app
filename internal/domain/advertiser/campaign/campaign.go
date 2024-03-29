package campaign

import (
	"net/url"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/adpiece"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type CampaignRepository interface {
	shared.EntityRepository[Campaign]
	CampaignsForAdvertiser(advertiserID shared.ID) ([]Campaign, error)
	RunningCampaigns() ([]Campaign, error)
}

type Campaign struct {
	ID               shared.ID
	Events           shared.Events
	Name             string
	AdvertiserUserID shared.ID
	Start            time.Time
	End              time.Time
	RunDuration      time.Duration
	IsPreempted      bool
}

func NewCampaign(name string, advertiserID shared.ID, start, end time.Time) Campaign {
	if end.Before(start) {
		end = start.Add(1 * 30 * 24 * time.Hour) // 1 month after start
	}
	return Campaign{
		ID:               shared.NewID(),
		Events:           make(shared.Events),
		Name:             name,
		AdvertiserUserID: advertiserID,
		Start:            start,
		End:              end,
		RunDuration:      end.Sub(start),
	}
}

func (c *Campaign) Quit() {
	c.IsPreempted = true
	c.RunDuration = time.Now().Sub(c.Start)
	c.Events.Emit("campaign/quitted", c)
}

func (c *Campaign) IsQuitted() bool {
	return c.End.Sub(c.Start) > c.RunDuration && c.IsPreempted
}

func (c *Campaign) AddAdPiece(slot adslot.AdSlotType, ref *url.URL, resource string, resourceMIME string) adpiece.AdPiece {
	ad := adpiece.NewAdPiece(c.ID, slot, ref, resource, resourceMIME)
	c.Events.Emit("campaign/adpieceadded", ad)
	return ad
}

func (c *Campaign) IsRunning() bool {
	now := time.Now()
	return (!c.IsQuitted()) && c.Start.Before(now) && c.End.After(now)
}
