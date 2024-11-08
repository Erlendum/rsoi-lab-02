package reservation

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

const (
	defaultTimeout = 5 * time.Second
)

type repository struct {
	conn *sqlx.DB
}

func NewRepository(conn *sqlx.DB) *repository {
	return &repository{conn: conn}
}

func (r *repository) CreateReservation(ctx context.Context, res *reservation) (int, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Insert("reservation").Columns("reservation_uid", "username", "book_uid", "library_uid", "status", "start_date", "till_date").
		Values(*res.ReservationUid, *res.UserName, *res.BookUid, *res.LibraryUid, *res.Status, *res.StartDate, *res.TillDate)
	query, args, err := builder.Suffix("RETURNING id").ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "failed to build query")
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var id int
	err = r.conn.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "failed to execute query")
	}

	return id, nil
}

func (r *repository) UpdateReservationStatus(ctx context.Context, uid string, username string, status string) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	builder := psql.Update("reservation").Set("status", status).Where(sq.And{sq.Eq{"reservation_uid": uid}, sq.Eq{"username": username}})

	query, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build query")
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	res, err := r.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to execute query")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (r *repository) GetReservation(ctx context.Context, uid string) (reservation, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	builder := psql.Select("reservation_uid", "username", "book_uid", "library_uid", "status", "start_date", "till_date").From("reservation").Where(sq.Eq{"reservation_uid": uid})

	query, args, err := builder.ToSql()
	if err != nil {
		return reservation{}, errors.Wrap(err, "failed to build query")
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	res := reservation{}

	err = r.conn.GetContext(ctx, &res, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return reservation{}, errNotFound
		}
		return reservation{}, errors.Wrap(err, "failed to execute query")
	}

	return res, nil
}

func (r *repository) GetReservations(ctx context.Context, username string, status string) ([]reservation, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	builder := psql.Select("reservation_uid", "username", "book_uid", "library_uid", "status", "start_date", "till_date").From("reservation").Where(sq.And{sq.Eq{"username": username}, sq.Eq{"status": status}})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query")
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	res := make([]reservation, 0)

	err = r.conn.GetContext(ctx, &res, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errNotFound
		}
		return nil, errors.Wrap(err, "failed to execute query")
	}

	return res, nil
}
