package usecase

import (
	"context"
	"errors"
	"github.com/Enthreeka/auth-lab3/internal/apperror"
	"github.com/Enthreeka/auth-lab3/internal/entity"
	"github.com/Enthreeka/auth-lab3/internal/repo"
	"github.com/Enthreeka/auth-lab3/pkg/logger"
	"time"
)

type tokenUsecase struct {
	tokenRepo repo.TokenRepository

	log *logger.Logger
}

func NewTokenUsecase(tokenRepo repo.TokenRepository, log *logger.Logger) TokenUsecase {
	return &tokenUsecase{
		tokenRepo: tokenRepo,
		log:       log,
	}
}

func (t *tokenUsecase) CreateToken(ctx context.Context, tokenSession string, userID string) error {
	token := new(entity.Token)

	timeDurationToken := time.Now().Add(1 * time.Hour)
	token.ExpiresAt = timeDurationToken
	token.Token = tokenSession
	token.UserID = userID

	err := t.tokenRepo.CreateToken(ctx, token)
	if err != nil {
		return apperror.NewAppError(err, "failed to create session token")
	}

	return nil
}

func (t *tokenUsecase) GetToken(ctx context.Context, token string) (*entity.Token, error) {
	tokenEntity, err := t.tokenRepo.GetToken(ctx, token)
	if err != nil {
		if errors.Is(err, apperror.ErrNoFoundRows) {
			return nil, apperror.ErrNoFoundRows
		}
		return nil, apperror.NewAppError(err, "failed to get token")
	}

	return tokenEntity, nil
}
