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
	id                        shared.ID
	events                    shared.Events
	publisherId               shared.ID
	url                       url.URL
	isVerified                bool
	isDeactivated             bool
	monetiaztionStatusHistory []MonetizationStatus
}

func NewSite(publisherId shared.ID, url url.URL) Site {
	return Site{
		id:                        shared.NewID(),
		publisherId:               publisherId,
		url:                       url,
		events:                    make(shared.Events),
		monetiaztionStatusHistory: []MonetizationStatus{{CanMonetize: true, Time: time.Now()}},
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

func (s *Site) ID() shared.ID {
	return s.id
}

func (s *Site) Verify() {
	s.isVerified = true
}

func (s *Site) IsVerified() bool {
	return s.isVerified
}

func (s *Site) PublisherId() shared.ID {
	return s.publisherId
}

func (s *Site) URL() url.URL {
	return s.url
}

func (s *Site) AddAdSlot(name string, adSlotType adslot.AdSlotType) adslot.AdSlot {
	return adslot.NewAdSlot(s.id, name, adSlotType)
}

func (s *Site) Deactivate() {
	s.isDeactivated = true
	s.events.Emit("site/deactivated", *s)
}

func (s *Site) IsActive() bool {
	return !s.isDeactivated
}
func (s *Site) DomainEvents() shared.Events {
	return s.events
}

func (s *Site) MonetizationStatusHistory() []MonetizationStatus {
	history := make([]MonetizationStatus, len(s.monetiaztionStatusHistory))
	copy(history, s.monetiaztionStatusHistory)
	return history
}

func (s *Site) CanServeAdPiece() bool {
	return s.IsVerified() && s.IsActive()
}

func (s *Site) Demonetize() {
	if s.lastMonietizationStatus() {
		s.monetiaztionStatusHistory = append(s.monetiaztionStatusHistory, MonetizationStatus{
			CanMonetize: false,
			Time:        time.Now(),
		})
	}
}

func (s *Site) AllowMonetization() {
	if !s.lastMonietizationStatus() {
		s.monetiaztionStatusHistory = append(s.monetiaztionStatusHistory, MonetizationStatus{
			CanMonetize: true,
			Time:        time.Now(),
		})
	}
}

func (s *Site) DemonetizeForTimePeriod(period time.Duration) {
	s.Demonetize()
	s.monetiaztionStatusHistory = append(s.monetiaztionStatusHistory, MonetizationStatus{
		CanMonetize: true,
		Time:        time.Now().Add(period),
	})
}

func (s *Site) CanMonetize() bool {
	return s.CanServeAdPiece() && s.canMonetizeCurrently()
}

func (s *Site) canMonetizeCurrently() bool {

	upperHistoryIndex := slices.IndexFunc(s.monetiaztionStatusHistory, func(ms MonetizationStatus) bool {
		return ms.Time.After(time.Now())
	})
	if upperHistoryIndex == -1 || upperHistoryIndex == 0 {
		return s.lastMonietizationStatus()
	}
	return s.monetiaztionStatusHistory[upperHistoryIndex-1].CanMonetize
}

func (s *Site) lastMonietizationStatus() bool {
	return s.monetiaztionStatusHistory[len(s.monetiaztionStatusHistory)-1].CanMonetize
}
