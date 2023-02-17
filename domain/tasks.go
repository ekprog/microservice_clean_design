package domain

import "time"

type Task struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type TasksRepository interface {
	All() ([]*Task, error)
	FetchById(int64) (*Task, error)
	Insert(task *Task) error
	Update(*Task) error
}

type TasksInteractor interface {
	GetAllTasks() (GetAllTasksResponse, error)
	CreateTask(string) (CreateTaskResponse, error)
}

// Responses (only for UseCase layer)

type GetAllTasksResponse struct {
	StatusCode string
	Tasks      []*Task
}

type CreateTaskResponse struct {
	StatusCode string
	Id         int64
}
