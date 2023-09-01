package adslotinteractors

import (
	"context"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/authorization"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/site"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/siteadslotservices"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type AdSlotInteractor struct {
	siteRepository   site.SiteRepository
	userRepository   user.UserRepository
	adslotRepository adslot.AdSlotRepository
	eventDispatcher  application.EventDispatcher
}

func NewAdSlotInteractor(siteRepository site.SiteRepository,
	userRepository user.UserRepository,
	adslotRepository adslot.AdSlotRepository,
	eventDispatcher application.EventDispatcher) AdSlotInteractor {
	return AdSlotInteractor{
		siteRepository,
		userRepository,
		adslotRepository,
		eventDispatcher,
	}
}

func (i *AdSlotInteractor) CreateAdSlotForSite(actorUserID shared.ID, siteID shared.ID, name string, adslotType adslot.AdSlotType) error {
	u, err := i.userRepository.Get(context.Background(), actorUserID)
	if err != nil {
		return err
	}

	s, err := i.siteRepository.Get(context.Background(), siteID)

	if err != nil {
		return err
	}

	if !authorization.CanModifySite(u, s) {
		return application.ErrAuthorization
	}

	slot := adslot.NewAdSlot(siteID, name, adslotType)

	err = i.adslotRepository.Save(context.Background(), slot)
	if err != nil {
		return err
	}

	i.eventDispatcher.DispatchAsync(slot.DomainEvents())

	return nil
}

func (i *AdSlotInteractor) ChangeAdSlotName(actorUserID shared.ID, adSlotID shared.ID, newName string) error {
	u, err := i.userRepository.Get(context.Background(), actorUserID)
	if err != nil {
		return err
	}

	slot, err := i.adslotRepository.Get(context.Background(), adSlotID)
	if err != nil {
		return err
	}

	s, err := i.siteRepository.Get(context.Background(), slot.SiteID())
	if err != nil {
		return err
	}

	if !authorization.CanModifyAdSlot(u, s, slot) {
		return application.ErrAuthorization
	}

	slot.SetName(newName)
	err = i.adslotRepository.Save(context.Background(), slot)
	if err != nil {
		return err
	}

	i.eventDispatcher.DispatchAsync(slot.DomainEvents())

	return nil

}

func (i *AdSlotInteractor) GetIntegrationHTMLSnippet(adSlotID shared.ID) (string, error) {
	slot, err := i.adslotRepository.Get(context.Background(), adSlotID)
	if err != nil {
		return "", err
	}
	s, err := i.siteRepository.Get(context.Background(), slot.SiteID())
	if err != nil {
		return "", err
	}

	return siteadslotservices.GenerateIntegrationHTMLSnippet(s, slot)
}

type SiteAdSlots map[string][]adslot.AdSlot

func (i *AdSlotInteractor) ActiveAdSlotsForPublisher(publisherID shared.ID) (SiteAdSlots, error) {
	siteAdSlots := make(SiteAdSlots)
	activeSites, err := i.siteRepository.ActiveSitesForPublisher(publisherID)
	if err != nil {
		return siteAdSlots, err
	}
	for _, activeSite := range activeSites {
		slots, err := i.adslotRepository.ActiveAdSlotsForSite(activeSite.ID())
		if err != nil {
			return siteAdSlots, err
		}
		siteAdSlots[activeSite.SiteDomain()] = slots
	}
	return siteAdSlots, nil
}

func (i *AdSlotInteractor) DeactivateAdSlot(actorUserID shared.ID, adSlotID shared.ID) error {
	u, err := i.userRepository.Get(context.Background(), actorUserID)
	if err != nil {
		return err
	}

	slot, err := i.adslotRepository.Get(context.Background(), adSlotID)
	if err != nil {
		return err
	}

	s, err := i.siteRepository.Get(context.Background(), slot.SiteID())
	if err != nil {
		return err
	}

	if !authorization.CanModifyAdSlot(u, s, slot) {
		return application.ErrAuthorization
	}

	slot.Deactivate()

	err = i.adslotRepository.Save(context.Background(), slot)
	if err != nil {
		return err
	}
	i.eventDispatcher.DispatchAsync(slot.DomainEvents())
	return nil
}

func (i *AdSlotInteractor) AdSlot(adSlotID shared.ID) (adslot.AdSlot, error) {
	return i.adslotRepository.Get(context.Background(), adSlotID)
}
