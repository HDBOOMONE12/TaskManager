package storage

import (
	"context"
	"database/sql"
	"errors"
)

type TelegramBindingRepository interface {
	SaveBinding(ctx context.Context, email string, chatID int64) error
	GetChatID(ctx context.Context, email string) (int64, error)
}

type TelegramBindingRepo struct {
	db *sql.DB
}

func NewTelegramBindingRepo(db *sql.DB) *TelegramBindingRepo {
	return &TelegramBindingRepo{
		db: db,
	}
}

func (r *TelegramBindingRepo) SaveBinding(ctx context.Context, email string, chatID int64) error {

	query := `
INSERT INTO telegram_bindings(email, chat_id)
VALUES ($1, $2)
ON CONFLICT (email) DO UPDATE SET chat_id = EXCLUDED.chat_id
`
	_, err := r.db.ExecContext(ctx, query, email, chatID)
	if err != nil {
		return err
	}

	return nil
}

func (r *TelegramBindingRepo) GetChatID(ctx context.Context, email string) (int64, error) {
	var chatID int64

	query := `
SELECT chat_id from telegram_bindings
where email = $1`

	row := r.db.QueryRowContext(ctx, query, email)
	err := row.Scan(&chatID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	return chatID, nil
}
