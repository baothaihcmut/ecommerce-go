package repositories

import (
	"context"
	"database/sql"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/persistence/sqlc/sqlc"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/admin"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/enums"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/outbound"
	"github.com/google/uuid"
)

type PostgresAdminRepo struct {
	queries *sqlc.Queries
	conn    *sql.DB
}

func (p *PostgresAdminRepo) Save(ctx context.Context, admin *admin.Admin) error {
	var insert bool
	_, err := p.queries.FindAdminByCriteria(ctx, sqlc.FindAdminByCriteriaParams{
		Criteria: "id",
		Value: sql.NullString{
			String: uuid.UUID(admin.Id).String(),
			Valid:  true,
		},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			insert = true
		} else {
			return err
		}
	}
	insert = false

	if insert {
		err := p.queries.InsertAdmin(ctx, sqlc.InsertAdminParams{
			ID: uuid.NullUUID{
				UUID:  uuid.UUID(admin.Id),
				Valid: true,
			},

			FirstName: sql.NullString{
				String: admin.FirstName,
				Valid:  true,
			},
			LastName: sql.NullString{
				String: admin.LastName,
				Valid:  true,
			},
			Email: sql.NullString{
				String: string(admin.Email),
				Valid:  true,
			},
			Password: sql.NullString{
				String: string(admin.Password),
				Valid:  true,
			},
			PhoneNumber: sql.NullString{
				String: string(admin.PhoneNumber),
				Valid:  true,
			},
			CurrentRefreshToken: sql.NullString{
				String: admin.CurrentRefreshToken.Value,
				Valid:  true,
			},
		})
		if err != nil {
			return err
		}
	} else {
		err := p.queries.UpdateAdmin(ctx, sqlc.UpdateAdminParams{
			ID: uuid.NullUUID{
				UUID:  uuid.UUID(admin.Id),
				Valid: true,
			},

			FirstName: sql.NullString{
				String: admin.FirstName,
				Valid:  true,
			},
			LastName: sql.NullString{
				String: admin.LastName,
				Valid:  true,
			},
			Email: sql.NullString{
				String: string(admin.Email),
				Valid:  true,
			},
			Password: sql.NullString{
				String: string(admin.Password),
				Valid:  true,
			},
			PhoneNumber: sql.NullString{
				String: string(admin.PhoneNumber),
				Valid:  true,
			},
			CurrentRefreshToken: sql.NullString{
				String: admin.CurrentRefreshToken.Value,
				Valid:  true,
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PostgresAdminRepo) FindByEmail(ctx context.Context, email valueobject.Email) (*admin.Admin, error) {
	res, err := p.queries.FindAdminByCriteria(ctx, sqlc.FindAdminByCriteriaParams{
		Criteria: "email",
		Value: sql.NullString{
			String: string(email),
			Valid:  true,
		},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &admin.Admin{
		Id:                  valueobject.UserId(res.ID),
		Email:               valueobject.Email(res.Email),
		Password:            valueobject.Password(res.Password),
		FirstName:           res.FirstName,
		LastName:            res.LastName,
		PhoneNumber:         valueobject.PhoneNumber(res.PhoneNumber),
		CurrentRefreshToken: valueobject.Token{Value: res.CurrentRefreshToken.String, TokenType: enums.REFRESH_TOKEN},
	}, nil
}
func NewPostgresAdminRepo(db *sql.DB) outbound.AdminRepository {
	return &PostgresAdminRepo{
		conn:    db,
		queries: sqlc.New(db),
	}
}
