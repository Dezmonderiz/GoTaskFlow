package repository

import (
	"context"
	"database/sql"
	"errors"

	"GoTaskFlow/internal/model"
)

type PostgresTaskRepository struct {
	db *sql.DB
}

func NewPostgresTaskRepository(db *sql.DB) *PostgresTaskRepository {
	return &PostgresTaskRepository{db: db}
}

func (r *PostgresTaskRepository) Create(ctx context.Context, task model.Task) (model.Task, error) {
	const query = `
		INSERT INTO tasks (title, description, status)
		VALUES ($1, $2, $3)
		RETURNING id, title, COALESCE(description, ''), status, created_at, updated_at
	`

	var created model.Task
	err := r.db.QueryRowContext(ctx, query, task.Title, task.Description, model.StatusTodo).Scan(
		&created.ID,
		&created.Title,
		&created.Description,
		&created.Status,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return model.Task{}, err
	}

	return created, nil
}

func (r *PostgresTaskRepository) List(ctx context.Context) ([]model.Task, error) {
	const query = `
		SELECT id, title, COALESCE(description, ''), status, created_at, updated_at
		FROM tasks
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]model.Task, 0)
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *PostgresTaskRepository) GetByID(ctx context.Context, id int64) (model.Task, error) {
	const query = `
		SELECT id, title, COALESCE(description, ''), status, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	return r.scanTask(ctx, query, id)
}

func (r *PostgresTaskRepository) UpdateStatus(ctx context.Context, id int64, status model.TaskStatus) (model.Task, error) {
	const query = `
		UPDATE tasks
		SET status = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING id, title, COALESCE(description, ''), status, created_at, updated_at
	`

	return r.scanTask(ctx, query, id, status)
}

func (r *PostgresTaskRepository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM tasks WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func (r *PostgresTaskRepository) Stats(ctx context.Context) (model.TaskStats, error) {
	const query = `
		SELECT
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE status = 'todo') AS todo,
			COUNT(*) FILTER (WHERE status = 'done') AS done,
			COUNT(*) FILTER (WHERE status = 'in_progress') AS in_progress
		FROM tasks
	`

	var stats model.TaskStats
	err := r.db.QueryRowContext(ctx, query).Scan(
		&stats.Total,
		&stats.Todo,
		&stats.Done,
		&stats.InProgress,
	)
	if err != nil {
		return model.TaskStats{}, err
	}

	return stats, nil
}

func (r *PostgresTaskRepository) scanTask(ctx context.Context, query string, args ...any) (model.Task, error) {
	var task model.Task
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Task{}, ErrTaskNotFound
	}
	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}
