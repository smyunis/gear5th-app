package identityinteractors

import (
	"context"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type UserProfile struct {
	UserID               string
	AuthenticationMethod user.AuthenticationMethod
	Fullname             user.PersonName
	PhoneNumber          user.PhoneNumber
	Email                user.Email
}

type UserAccountInteractor struct {
	userRepository        user.UserRepository
	managedUserRepository user.ManagedUserRepository
}

func NewUserAccountInteractor(
	userRepository user.UserRepository,
	managedUserRepository user.ManagedUserRepository) UserAccountInteractor {
	return UserAccountInteractor{
		userRepository,
		managedUserRepository,
	}
}

func (i *UserAccountInteractor) User(userID shared.ID) (p UserProfile, err error) {
	u, err := i.userRepository.Get(context.Background(), userID)
	if err != nil {
		return UserProfile{}, nil
	}
	p.UserID = userID.String()
	p.AuthenticationMethod = u.AuthenticationMethod
	p.PhoneNumber = u.PhoneNumber
	p.Email = u.Email

	if u.AuthenticationMethod == user.Managed {
		mu, err := i.managedUserRepository.Get(context.Background(), userID)
		if err != nil {
			return UserProfile{}, nil
		}
		p.Fullname = mu.Name()
	}

	return p, nil
}

func (i *UserAccountInteractor) SetUser(userID shared.ID, userProfile UserProfile) error {
	u, err := i.userRepository.Get(context.Background(), userID)
	if err != nil {
		return err
	}
	u.PhoneNumber = userProfile.PhoneNumber
	i.userRepository.Save(context.Background(), u)

	if u.AuthenticationMethod == user.Managed {
		mu, err := i.managedUserRepository.Get(context.Background(), userID)
		if err != nil {
			return err
		}
		mu.SetName(userProfile.Fullname)
		i.managedUserRepository.Save(context.Background(), mu)
	}

	return nil
}

func (i *UserAccountInteractor) ChangePassword(userID shared.ID, currentPassword string, newPassword string) error {
	mu, err := i.managedUserRepository.Get(context.Background(), userID)
	if err != nil {
		return err
	}
	if !mu.IsPasswordCorrect(currentPassword) {
		return  application.ErrAuthorization
	}

	err = mu.SetPassword(newPassword)
	if err != nil {
		return err
	}

	i.managedUserRepository.Save(context.Background(), mu)
	return nil
}
