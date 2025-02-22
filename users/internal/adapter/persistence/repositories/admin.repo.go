package repositories

import (
	"context"
	"database/sql"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/tracing"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/persistence/sqlc/sqlc"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/admin"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/outbound/repositories"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type PostgresAdminRepo struct {
	queries *sqlc.Queries
	conn    *sql.DB
	tracer  trace.Tracer
}

func (p *PostgresAdminRepo) Save(ctx context.Context, admin *admin.Admin) (err error) {
	ctx, span := tracing.StartSpan(ctx, p.tracer, "Admin.save: database", nil)
	defer tracing.EndSpan(span, err, nil)
	var insert bool
	_, err = p.queries.FindAdminByCriteria(ctx, sqlc.FindAdminByCriteriaParams{
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
				String: admin.CurrentRefreshToken,
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
				String: admin.CurrentRefreshToken,
				Valid:  true,
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PostgresAdminRepo) FindByEmail(ctx context.Context, email valueobject.Email) (resp *admin.Admin, err error) {
	ctx, span := tracing.StartSpan(ctx, p.tracer, "Admin.FindByEmail: database ", nil)
	defer tracing.EndSpan(span, err, nil)
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
		CurrentRefreshToken: res.CurrentRefreshToken.String,
	}, nil
}
func NewPostgresAdminRepo(db *sql.DB, tracer trace.Tracer) repositories.AdminRepository {
	return &PostgresAdminRepo{
		conn:    db,
		queries: sqlc.New(db),
		tracer:  tracer,
	}
}
