package repos

import (
	"database/sql"
	"microservice/app/core"
	"microservice/domain"
)

type TaskDBRepo struct {
	log core.Logger
	db  *sql.DB
}

func NewTaskDBRepo(log core.Logger, db *sql.DB) *TaskDBRepo {
	return &TaskDBRepo{
		db:  db,
		log: log,
	}
}

func (r *TaskDBRepo) All() ([]*domain.Task, error) {
	var tasks []*domain.Task

	query := `SELECT id, name, created_at, updated_at FROM tasks ORDER BY id;`
	raws, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	for raws.Next() {
		task := &domain.Task{}
		err = raws.Scan(&task.Id, &task.Name, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskDBRepo) FetchById(id int64) (*domain.Task, error) {
	task := &domain.Task{
		Id: id,
	}
	query := `SELECT name, created_at, updated_at FROM tasks WHERE id=$1`
	err := r.db.QueryRow(query, id).Scan(&task.Name, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return task, err
	}
	return task, nil
}

func (r *TaskDBRepo) Insert(task *domain.Task) error {
	var id int64
	query := "INSERT INTO tasks (name) VALUES ($1) returning id"
	err := r.db.QueryRow(query, task.Name).Scan(&id)
	if err != nil {
		return err
	}
	task.Id = id
	return nil
}

func (r *TaskDBRepo) Update(task *domain.Task) error {
	query := "UPDATE tasks SET name=$2, updated_at=now() WHERE id=$1"
	_, err := r.db.Exec(query, task.Id, task.Name)
	if err != nil {
		return err
	}
	return nil
}
