package usecase

import (
	"context"
	"github.com/Enthreeka/auth-lab3/internal/entity"
)

type UserUsecase interface {
	GetUserByLogin(ctx context.Context, login string, password string) (*entity.User, error)
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User, userID string) error
}

type TokenUsecase interface {
	CreateToken(ctx context.Context, tokenSession string, userID string) error
	GetToken(ctx context.Context, token string) (*entity.Token, error)
}
