package repo

import (
	"context"
	"github.com/Enthreeka/auth-lab3/internal/entity"
)

type UserRepository interface {
	GetUserByLogin(ctx context.Context, login string) (*entity.User, error)
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) error
}

type TokenRepository interface {
	GetToken(ctx context.Context, token string) (*entity.Token, error)
	CreateToken(ctx context.Context, token *entity.Token) error
}
