package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/Nicholas2012/time-tracker/internal/models"
	goqu "github.com/doug-martin/goqu/v9"
)

type Repository struct {
	db *sql.DB
}

type UserList struct {
	Users []models.User
	Count int
	Pages int
	Page  int
}

type ListOpts struct {
	Page  int
	Limit int
	Name  string // поиск по имени
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (name, surname, patronymic, passport_serie, passport_number) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	row := r.db.QueryRowContext(ctx, query, user.Name, user.Surname, user.Patronymic, user.PassportSerie, user.PassportNumber)
	if err := row.Scan(&user.ID); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetUser(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT name, surname, patronymic, passport_serie, passport_number FROM users WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)
	user := &models.User{ID: id}
	if err := row.Scan(&user.Name, &user.Surname, &user.Patronymic, &user.PassportSerie, &user.PassportNumber); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) DeleteUser(ctx context.Context, user *models.User) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, user.ID)
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

func (r *Repository) UpdateUser(ctx context.Context, user *models.User) error {

	query := `UPDATE users
		SET name = $1,
		surname = $2,
		patronymic = $3,
		passport_serie = $4,
		passport_number = $5
		WHERE id = $6`

	result, err := r.db.ExecContext(ctx, query, user.Name, user.Surname, user.Patronymic, user.PassportSerie, user.PassportNumber, user.ID)
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

func (r *Repository) ListUsers(opts ListOpts) (*UserList, error) {
	qb := goqu.From("users")

	if opts.Name != "" {
		qb = qb.Where(goqu.L("name || ' ' || surname || ' ' || patronymic ILIKE ?", fmt.Sprint("%", opts.Name, "%")))
	}

	// get total count
	countQuery, countArgs, err := qb.Select(goqu.L("COUNT(*)")).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("build count query: %w", err)
	}
	slog.Debug("list count query", "query", countQuery, "repository", "users")
	var count int
	if err := r.db.QueryRow(countQuery, countArgs...).Scan(&count); err != nil {
		return nil, fmt.Errorf("count: %w", err)
	}

	userList := UserList{
		Users: make([]models.User, 0, opts.Limit),
		Count: count,
		Pages: count / opts.Limit,
		Page:  opts.Page,
	}

	// calculate offset and fix pages if needed
	offset := (opts.Page - 1) * opts.Limit
	if offset < 0 {
		offset = 0
		opts.Page = 1
	}
	if count%opts.Limit > 0 {
		userList.Pages++
	}
	if userList.Pages == 0 {
		userList.Page = 0
		return &userList, nil
	}

	// get data
	selectQuery, selectArgs, err := qb.
		Select("id", "name", "surname", "patronymic", "passport_serie", "passport_number").
		Offset(uint(offset)).
		Limit(uint(opts.Limit)).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("build select query: %w", err)
	}

	slog.Debug("list query", "query", countQuery, "args", countArgs, "repository", "users")

	rows, err := r.db.Query(selectQuery, selectArgs...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Debug("db rows close", "err", err, "repository", "users")
		}
	}()

	for rows.Next() {
		var user models.User

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Surname,
			&user.Patronymic,
			&user.PassportSerie,
			&user.PassportNumber,
		)
		if err != nil {
			return nil, err
		}

		userList.Users = append(userList.Users, user)
	}

	return &userList, nil
}

func (r *Repository) CreateTask(ctx context.Context, task *models.Task) error {
	query := `INSERT INTO tasks (user_id, start_time, end_time, minutes) 
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	row := r.db.QueryRowContext(ctx, query, task.UserID, task.Since, task.Until, task.Minutes)
	if err := row.Scan(&task.ID); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetTask(ctx context.Context, userID, taskID int) (*models.Task, error) {
	query := `SELECT id, user_id, start_time, end_time, minutes FROM tasks WHERE user_id = $1 AND id = $2`

	row := r.db.QueryRowContext(ctx, query, userID, taskID)
	task := &models.Task{ID: taskID}
	if err := row.Scan(&task.UserID, &task.Since, &task.Until, &task.Minutes); err != nil {
		return nil, err
	}

	return task, nil
}

func (r *Repository) UpdateTask(ctx context.Context, task *models.Task) error {
	query := `UPDATE tasks SET start_time = $1, end_time = $2, minutes = $3 WHERE id = $4`

	_, err := r.db.ExecContext(ctx, query, task.Since, task.Until, task.Minutes, task.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ListTasks(ctx context.Context, userID int) ([]models.Task, error) {
	query := `SELECT id, user_id, start_time, end_time, minutes FROM tasks WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Debug("db rows close", "err", err, "repository", "tasks")
		}
	}()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task

		err := rows.Scan(&task.ID, &task.UserID, &task.Since, &task.Until, &task.Minutes)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
