package advertiserinteractors

import (
	"context"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/advertiser"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type AdvertiserSignUpUnitOfWork interface {
	Save(ctx context.Context, u user.User, a advertiser.Advertiser) error
}

type AdvertiserInteractor struct {
	advertiserRepository advertiser.AdvertiserRepository
	userRepository       user.UserRepository
	advertiserSignUpUoW  AdvertiserSignUpUnitOfWork
	eventDispatcher      application.EventDispatcher
}

func NewAdvertiserInteractor(advertiserRepository advertiser.AdvertiserRepository,
	userRepository user.UserRepository,
	advertiserSignUpUoW AdvertiserSignUpUnitOfWork,
	eventDispatcher application.EventDispatcher) AdvertiserInteractor {
	return AdvertiserInteractor{
		advertiserRepository,
		userRepository,
		advertiserSignUpUoW,
		eventDispatcher,
	}
}

func (i *AdvertiserInteractor) SignUpAdvertiser(email user.Email, ph user.PhoneNumber, name, note string) error {
	u := user.NewUser(email)
	u.PhoneNumber = ph
	a := u.SignUpAdvertiser(name)
	a.Note = note

	err := i.advertiserSignUpUoW.Save(context.Background(), u, a)
	if err != nil {
		return err
	}

	i.eventDispatcher.DispatchAsync(u.Events)
	return nil
}

func (i *AdvertiserInteractor) Adveritsers() ([]advertiser.Advertiser, error) {
	return i.advertiserRepository.Advertisers()
}

type AdveritserDetails struct {
	User       user.User
	Advertiser advertiser.Advertiser
}

func (i *AdvertiserInteractor) Advertiser(advertiserID shared.ID) (AdveritserDetails, error) {
	a, err := i.advertiserRepository.Get(context.Background(), advertiserID)
	if err != nil {
		return AdveritserDetails{}, err
	}
	u, err := i.userRepository.Get(context.Background(), advertiserID)
	if err != nil {
		return AdveritserDetails{}, err
	}
	return AdveritserDetails{u, a}, nil
}
