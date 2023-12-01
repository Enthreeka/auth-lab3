package usecase

import (
	"context"
	"github.com/Enthreeka/auth-lab3/internal/apperror"
	"github.com/Enthreeka/auth-lab3/internal/entity"
	"github.com/Enthreeka/auth-lab3/internal/repo"
	"github.com/Enthreeka/auth-lab3/pkg/logger"
)

type userUsecase struct {
	userRepo repo.UserRepository
	log      *logger.Logger
}

func NewUserUsecase(userRepo repo.UserRepository, log *logger.Logger) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
		log:      log,
	}
}

func (u *userUsecase) CreateUser(ctx context.Context, user *entity.User, userID string) error {
	argon := NewPassword("")

	user.ID = userID

	hash, err := argon.GenerateHashFromPassword(user.ID, user.Password)
	if err != nil {
		return err
	}

	user.Password = hash

	err = u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (u *userUsecase) GetUserByLogin(ctx context.Context, login string, password string) (*entity.User, error) {
	user, err := u.userRepo.GetUserByLogin(ctx, login)
	if err != nil {
		if err == apperror.ErrNoFoundRows {
			return nil, apperror.ErrNoFoundRows
		}
		return nil, apperror.NewAppError(err, "failed to get user from postgres")
	}

	argon := NewPassword("")

	err = argon.VerifyPassword(user.Password, user.ID, password)
	if err != nil {
		return nil, apperror.ErrHashPasswordsNotEqual
	}

	return user, nil
}

func (u *userUsecase) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	user, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		if err == apperror.ErrNoFoundRows {
			return nil, apperror.ErrNoFoundRows
		}
		return nil, apperror.NewAppError(err, "failed to get user from postgres")
	}

	return user, nil
}
