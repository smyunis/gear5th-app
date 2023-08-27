package adslot

import "gitlab.com/gear5th/gear5th-app/internal/domain/shared"

type AdSlotRepository interface {
	shared.EntityRepository[AdSlot]
	ActiveAdSlotsForSite(siteID shared.ID) ([]AdSlot, error)
}

type AdSlotType int

const (
	_ AdSlotType = iota
	Horizontal
	Vertical
	Box
)

type AdSlot struct {
	id            shared.ID
	name          string
	siteID        shared.ID
	slotType      AdSlotType
	isDeactivated bool
}

func NewAdSlot(siteID shared.ID, name string, slotType AdSlotType) AdSlot {
	return AdSlot{
		id:       shared.NewID(),
		siteID:   siteID,
		name:     name,
		slotType: slotType,
	}
}

func (s *AdSlot) AdSlotID() shared.ID {
	return s.id
}

func (s *AdSlot) Name() string {
	return s.name
}

func (s *AdSlot) SetName(name string) {
	s.name = name
}

func (s *AdSlot) Type() AdSlotType {
	return s.slotType
}

func (s *AdSlot) SiteID() shared.ID {
	return s.siteID
}

func (s *AdSlot) IsDeactivated() bool {
	return s.isDeactivated
}

func (s *AdSlot) Deactivate() {
	s.isDeactivated = true
}


