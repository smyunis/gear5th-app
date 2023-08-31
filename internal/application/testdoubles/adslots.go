package testdoubles

import (
	"context"

	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type AdSlotRepositoryStub struct{}

// ActiveAdSlotsForSite implements adslot.AdSlotRepository.
func (AdSlotRepositoryStub) ActiveAdSlotsForSite(siteID shared.ID) ([]adslot.AdSlot, error) {
	panic("unimplemented")
}

// Get implements adslot.AdSlotRepository.
func (AdSlotRepositoryStub) Get(ctx context.Context, id shared.ID) (adslot.AdSlot, error) {
	return adslot.ReconstituteAdSlot(id, StubID, "my-adslot-x", adslot.Vertical, false), nil
}

// Save implements adslot.AdSlotRepository.
func (AdSlotRepositoryStub) Save(ctx context.Context, e adslot.AdSlot) error {
	return nil
}

var x adslot.AdSlotRepository = AdSlotRepositoryStub{}
