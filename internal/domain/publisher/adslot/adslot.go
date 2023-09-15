package adslot

import (
	"fmt"

	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

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

func (t AdSlotType) String() string {
	switch t {
	case Horizontal:
		return "Horizontal"
	case Vertical:
		return "Vertical"
	case Box:
		return "Box"
	default:
		return ""
	}
}

func AdSlotTypeFromString(slotType string) AdSlotType {
	switch slotType {
	case "horizontal":
		return Horizontal
	case "vertical":
		return Vertical
	case "box":
		return Box
	default:
		return 0
	}
}

type Dimentions struct {
	Width  int
	Height int
}

func (d Dimentions) String() string {
	return fmt.Sprintf("%dx%d", d.Width, d.Height)
}

func (t AdSlotType) Dimentions() Dimentions {
	switch t {
	case Horizontal:
		return Dimentions{
			Width:  728,
			Height: 90,
		}
	case Vertical:
		return Dimentions{
			Width:  160,
			Height: 600,
		}
	case Box:
		return Dimentions{
			Width:  300,
			Height: 250,
		}
	default:
		return Dimentions{
			Width:  0,
			Height: 0,
		}
	}
}

type AdSlot struct {
	id            shared.ID
	siteID        shared.ID
	name          string
	slotType      AdSlotType
	isDeactivated bool
	events        shared.Events
}

func NewAdSlot(siteID shared.ID, name string, slotType AdSlotType) AdSlot {
	return AdSlot{
		id:       shared.NewID(),
		siteID:   siteID,
		name:     name,
		slotType: slotType,
		events:   make(shared.Events),
	}
}

func ReconstituteAdSlot(
	id shared.ID,
	siteID shared.ID,
	name string,
	slotType AdSlotType,
	isDeactivated bool) AdSlot {
	return AdSlot{
		id,
		siteID,
		name,
		slotType,
		isDeactivated,
		make(shared.Events),
	}
}

func (s *AdSlot) ID() shared.ID {
	return s.id
}

func (s *AdSlot) Name() string {
	return s.name
}

func (s *AdSlot) SetName(name string) {
	s.name = name
}

func (s *AdSlot) AdSlotType() AdSlotType {
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

func (s *AdSlot) DomainEvents() shared.Events {
	return s.events
}
