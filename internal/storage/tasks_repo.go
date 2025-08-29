//go:generate mockgen -source=tasks_repo.go -destination=../mocks/mock_task_repo.go -package=mocks
package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/HDBOOMONE12/TaskManager/internal/entity"
	"time"
)

type TaskRepository interface {
	Create(ctx context.Context, task *entity.Task) error
	GetByID(ctx context.Context, id int64) (entity.Task, error)
	GetByUserID(ctx context.Context, userID int64) ([]entity.Task, error)
	Update(ctx context.Context, task *entity.Task) (entity.Task, error)
	Patch(ctx context.Context, uid, tid int64, title, desc, status *string, priority *int, dueAtProvided bool, dueAt *time.Time) (entity.Task, error)
	Delete(ctx context.Context, id int64) error
}

type TaskRepo struct {
	db *sql.DB
}

func NewTaskRepo(db *sql.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) Create(ctx context.Context, task *entity.Task) error {
	query := `
		INSERT INTO tasks (user_id, title, description, status, due_date, priority)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at;
	`

	row := r.db.QueryRowContext(ctx, query,
		task.UserID,
		task.Title,
		task.Description,
		task.Status,
		task.DueAt,
		task.Priority,
	)

	err := row.Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *TaskRepo) GetByID(ctx context.Context, id int64) (entity.Task, error) {
	query := `
		SELECT id, user_id, title, description, status, due_date, priority, created_at, updated_at
		FROM tasks
		WHERE id = $1;
	`

	var task entity.Task

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.DueAt,
		&task.Priority,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return entity.Task{}, err
	}

	return task, err
}

func (r *TaskRepo) GetByUserID(ctx context.Context, userID int64) ([]entity.Task, error) {
	query := `
		SELECT id, user_id, title, description, status, due_date, priority, created_at, updated_at
		FROM tasks
		WHERE user_id = $1
		ORDER BY created_at DESC;
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []entity.Task

	for rows.Next() {
		var task entity.Task
		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.DueAt,
			&task.Priority,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	query := `
		UPDATE tasks
		SET status = $1, updated_at = now()
		WHERE id = $2;
	`
	res, err := r.db.ExecContext(ctx, query, status, id)
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

func (r *TaskRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM tasks WHERE id = $1;`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *TaskRepo) Update(ctx context.Context, t *entity.Task) (entity.Task, error) {
	query := `
		UPDATE tasks
		SET title = $1,
		    description = $2,
		    status = $3,
		    due_date = $4,
		    priority = $5,
		    updated_at = now()
		WHERE id = $6 AND user_id = $7
		RETURNING id, user_id, title, description, status, due_date, priority, created_at, updated_at;
	`

	var out entity.Task
	err := r.db.QueryRowContext(ctx, query,
		t.Title,
		t.Description,
		t.Status,
		t.DueAt,
		t.Priority,
		t.ID,
		t.UserID,
	).Scan(
		&out.ID,
		&out.UserID,
		&out.Title,
		&out.Description,
		&out.Status,
		&out.DueAt,
		&out.Priority,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Task{}, sql.ErrNoRows
		}
		return entity.Task{}, err
	}
	return out, nil
}

func (r *TaskRepo) Patch(
	ctx context.Context,
	uid, tid int64,
	title, desc, status *string,
	priority *int,
	dueAtProvided bool,
	dueAt *time.Time,
) (entity.Task, error) {

	if title == nil && desc == nil && status == nil && priority == nil && !dueAtProvided {
		return entity.Task{}, errors.New("nothing to update")
	}

	query := "UPDATE tasks SET "
	args := []interface{}{}
	idx := 1

	if title != nil {
		query += fmt.Sprintf("title = $%d", idx)
		args = append(args, *title)
		idx++
	}
	if desc != nil {
		if len(args) > 0 {
			query += ", "
		}
		query += fmt.Sprintf("description = $%d", idx)
		args = append(args, *desc)
		idx++
	}
	if status != nil {
		if len(args) > 0 {
			query += ", "
		}
		query += fmt.Sprintf("status = $%d", idx)
		args = append(args, *status)
		idx++
	}
	if priority != nil {
		if len(args) > 0 {
			query += ", "
		}
		query += fmt.Sprintf("priority = $%d", idx)
		args = append(args, *priority)
		idx++
	}
	if dueAtProvided {
		if len(args) > 0 {
			query += ", "
		}
		query += fmt.Sprintf("due_date = $%d", idx)
		args = append(args, dueAt)
		idx++
	}

	query += fmt.Sprintf(", updated_at = now() WHERE id = $%d AND user_id = $%d ", idx, idx+1)
	args = append(args, tid, uid)

	query += "RETURNING id, user_id, title, description, status, due_date, priority, created_at, updated_at;"

	var out entity.Task
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&out.ID,
		&out.UserID,
		&out.Title,
		&out.Description,
		&out.Status,
		&out.DueAt,
		&out.Priority,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Task{}, sql.ErrNoRows
		}
		return entity.Task{}, err
	}
	return out, nil
}
