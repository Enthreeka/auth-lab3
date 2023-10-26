package repo

import (
	"context"
	"github.com/Enthreeka/auth-lab3/internal/apperror"
	"github.com/Enthreeka/auth-lab3/internal/entity"
	"github.com/Enthreeka/auth-lab3/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

type userRepoPG struct {
	*postgres.Postgres
}

func NewUserRepoPG(pg *postgres.Postgres) UserRepository {
	return &userRepoPG{
		pg,
	}
}

func (u *userRepoPG) GetUserByLogin(ctx context.Context, login string) (*entity.User, error) {
	query := `SELECT "user".id,"user".password,"user".login, role.id, role.role
        FROM "user"
        JOIN role ON "user".role_id = role.id
        WHERE "user".login = $1`
	user := new(entity.User)

	err := u.Pool.QueryRow(ctx, query, login).Scan(
		&user.ID,
		&user.Password,
		&user.Login,
		&user.Role.ID,
		&user.Role.Role,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperror.ErrNoFoundRows
		}
		return nil, err
	}

	return user, nil
}

func (u *userRepoPG) CreateUser(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO "user" (id,login,password) VALUES ($1,$2,$3)`

	_, err := u.Pool.Exec(ctx, query, user.ID, user.Login, user.Password)

	return err
}

func (u *userRepoPG) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	query := `SELECT "user".id,"user".password,"user".login, role.id, role.role
        FROM "user"
        JOIN role ON "user".role_id = role.id
        WHERE "user".id = $1`

	user := new(entity.User)

	err := u.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Password,
		&user.Login,
		&user.Role.ID,
		&user.Role.Role,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperror.ErrNoFoundRows
		}
		return nil, err
	}

	return user, nil
}
