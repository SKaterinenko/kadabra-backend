package users_postgres

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"kadabra/internal/core"
	"kadabra/internal/core/config"
	user_model "kadabra/internal/features/users/model"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUsersPostgres(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Register(ctx context.Context, user *user_model.User) (*user_model.User, error) {
	query, args, err := config.Psql.
		Insert("users").
		Columns(
			"first_name",
			"last_name",
			"email",
			"birth_date",
			"phone_number",
			"gender",
			"password_hash",
		).
		Values(
			user.FirstName,
			user.LastName,
			user.Email,
			user.BirthDate,
			user.PhoneNumber,
			user.Gender,
			user.PasswordHash,
		).
		Suffix("RETURNING id, first_name, last_name, email, birth_date, phone_number, gender, password_hash, created_at, updated_at").
		ToSql()

	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	newUser, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[user_model.User])
	if err != nil {
		return nil, core.RowsError(err)
	}

	return newUser, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user_model.User, error) {
	query, args, err := config.Psql.
		Select(
			"id",
			"first_name",
			"last_name",
			"email",
			"birth_date",
			"phone_number",
			"gender",
			"password_hash",
			"created_at",
			"updated_at",
		).
		From("users").
		Where(sq.Eq{"email": email}).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[user_model.User])
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*user_model.User, error) {
	query, args, err := config.Psql.
		Select(
			"id",
			"first_name",
			"last_name",
			"email",
			"birth_date",
			"phone_number",
			"gender",
			"password_hash",
			"created_at",
			"updated_at",
		).
		From("users").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[user_model.User])
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*user_model.User, error) {
	query, args, err := config.Psql.
		Select(
			"id",
			"first_name",
			"last_name",
			"email",
			"birth_date",
			"phone_number",
			"gender",
			"password_hash",
			"created_at",
			"updated_at",
		).
		From("users").
		Where(sq.Eq{"phone_number": phoneNumber}).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[user_model.User])
	if err != nil {
		return nil, err
	}

	return user, nil
}
