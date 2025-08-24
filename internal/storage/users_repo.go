package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/HDBOOMONE12/TaskManager/internal/entity"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *entity.User) error {
	row := r.db.QueryRowContext(ctx,
		"INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id, created_at, updated_at",
		user.Username, user.Email,
	)
	return row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepo) GetAll(ctx context.Context) ([]entity.User, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, username, email, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (entity.User, error) {
	row := r.db.QueryRowContext(ctx,
		"SELECT id, username, email, created_at, updated_at FROM users WHERE id = $1",
		id,
	)

	var user entity.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, sql.ErrNoRows
		}
		return entity.User{}, err
	}
	return user, nil
}

func (r *UserRepo) Update(ctx context.Context, id int64, name, email string) (entity.User, error) {
	row := r.db.QueryRowContext(ctx,
		"UPDATE users SET username = $1, email = $2, updated_at = now() WHERE id = $3 RETURNING id, username, email, created_at, updated_at",
		name, email, id,
	)
	var user entity.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, sql.ErrNoRows
		}
		return entity.User{}, err
	}
	return user, nil
}

func (r *UserRepo) Patch(ctx context.Context, id int64, name, email *string) (entity.User, error) {
	if name == nil && email == nil {
		return entity.User{}, errors.New("nothing to update")
	}

	query := "UPDATE users SET "
	params := []interface{}{}
	idx := 1

	if name != nil {
		query += fmt.Sprintf("username = $%d", idx)
		params = append(params, *name)
		idx++
	}

	if email != nil {
		if len(params) > 0 {
			query += ", "
		}
		query += fmt.Sprintf("email = $%d", idx)
		params = append(params, *email)
		idx++
	}

	query += fmt.Sprintf(", updated_at = now() WHERE id = $%d RETURNING id, username, email, created_at, updated_at", idx)
	params = append(params, id)

	row := r.db.QueryRowContext(ctx, query, params...)

	var user entity.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r *UserRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
