package repo

import (
	"context"
	"github.com/Enthreeka/auth-lab3/internal/apperror"
	"github.com/Enthreeka/auth-lab3/internal/entity"
	"github.com/Enthreeka/auth-lab3/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

type tokenRepoPG struct {
	*postgres.Postgres
}

func NewTokenRepoPG(pg *postgres.Postgres) TokenRepository {
	return &tokenRepoPG{
		pg,
	}
}

func (t *tokenRepoPG) GetToken(ctx context.Context, token string) (*entity.Token, error) {
	query := `SELECT  token,expires_at,user_id FROM session WHERE token = $1`

	sessionData := new(entity.Token)

	err := t.Pool.QueryRow(ctx, query, token).Scan(
		&sessionData.Token,
		&sessionData.ExpiresAt,
		&sessionData.UserID,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperror.ErrNoFoundRows
		}
		return nil, err
	}

	return sessionData, nil
}

func (t *tokenRepoPG) CreateToken(ctx context.Context, token *entity.Token) error {
	query := `INSERT INTO session (token, user_id, expires_at) VALUES ($1,$2,$3)`

	_, err := t.Pool.Exec(ctx, query, token.Token, token.UserID, token.ExpiresAt)

	return err
}
