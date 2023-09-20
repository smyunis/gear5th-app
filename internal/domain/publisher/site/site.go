package site

import (
	"errors"
	"net/url"
	"slices"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

var ErrDuplicateAdSlot = errors.New("duplicate adslot exists")

type SiteRepository interface {
	shared.EntityRepository[Site]
	ActiveSitesForPublisher(publisherID shared.ID) ([]Site, error)
}

type MonetizationStatus struct {
	CanMonetize bool
	Time        time.Time
}

type Site struct {
	ID                        shared.ID
	Events                    shared.Events
	PublisherID               shared.ID
	URL                       url.URL
	IsVerified                bool
	IsDeactivated             bool
	MonetiaztionStatusHistory []MonetizationStatus
}

func NewSite(publisherId shared.ID, url url.URL) Site {
	return Site{
		ID:                        shared.NewID(),
		PublisherID:               publisherId,
		URL:                       url,
		Events:                    make(shared.Events),
		MonetiaztionStatusHistory: []MonetizationStatus{{CanMonetize: true, Time: time.Now()}},
	}
}

func ReconstituteSite(
	id shared.ID,
	publisherId shared.ID,
	url url.URL,
	isVerified bool,
	isDeactivated bool,
	monetiaztionStatusHistory []MonetizationStatus) Site {
	return Site{
		id,
		make(shared.Events),
		publisherId,
		url,
		isVerified,
		isDeactivated,
		monetiaztionStatusHistory,
	}
}

func (s *Site) Verify() {
	s.IsVerified = true
}

func (s *Site) SiteDomain() string {
	return s.URL.Hostname()
}

func (s *Site) AddAdSlot(name string, adSlotType adslot.AdSlotType) adslot.AdSlot {
	return adslot.NewAdSlot(s.ID, name, adSlotType)
}

func (s *Site) Deactivate() {
	s.IsDeactivated = true
	s.Events.Emit("site/deactivated", *s)
}

func (s *Site) MonetizationStatusHistory() []MonetizationStatus {
	history := make([]MonetizationStatus, len(s.MonetiaztionStatusHistory))
	copy(history, s.MonetiaztionStatusHistory)
	return history
}

func (s *Site) CanServeAdPiece() bool {
	return s.IsVerified && !s.IsDeactivated
}

func (s *Site) Demonetize() {
	if s.lastMonietizationStatus() {
		s.MonetiaztionStatusHistory = append(s.MonetiaztionStatusHistory, MonetizationStatus{
			CanMonetize: false,
			Time:        time.Now(),
		})
	}
}

func (s *Site) AllowMonetization() {
	if !s.lastMonietizationStatus() {
		s.MonetiaztionStatusHistory = append(s.MonetiaztionStatusHistory, MonetizationStatus{
			CanMonetize: true,
			Time:        time.Now(),
		})
	}
}

func (s *Site) DemonetizeForTimePeriod(period time.Duration) {
	s.Demonetize()
	s.MonetiaztionStatusHistory = append(s.MonetiaztionStatusHistory, MonetizationStatus{
		CanMonetize: true,
		Time:        time.Now().Add(period),
	})
}

func (s *Site) CanMonetize() bool {
	return s.CanServeAdPiece() && s.canMonetizeCurrently()
}

func (s *Site) canMonetizeCurrently() bool {

	upperHistoryIndex := slices.IndexFunc(s.MonetiaztionStatusHistory, func(ms MonetizationStatus) bool {
		return ms.Time.After(time.Now())
	})
	if upperHistoryIndex == -1 || upperHistoryIndex == 0 {
		return s.lastMonietizationStatus()
	}
	return s.MonetiaztionStatusHistory[upperHistoryIndex-1].CanMonetize
}

func (s *Site) lastMonietizationStatus() bool {
	return s.MonetiaztionStatusHistory[len(s.MonetiaztionStatusHistory)-1].CanMonetize
}
