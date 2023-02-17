package interactors

import (
	"github.com/pkg/errors"
	"microservice_clean_design/app"
	"microservice_clean_design/domain"
)

type TasksInteractor struct {
	log       app.Logger
	tasksRepo domain.TasksRepository
}

func NewTasksUCase(log app.Logger, tasksRepo domain.TasksRepository) domain.TasksInteractor {
	return &TasksInteractor{
		log:       log,
		tasksRepo: tasksRepo,
	}
}

func (i *TasksInteractor) GetAllTasks() (domain.GetAllTasksResponse, error) {
	tasks, err := i.tasksRepo.All()
	if err != nil {
		return domain.GetAllTasksResponse{}, errors.Wrap(err, "Cannot fetch all tasks")
	}

	return domain.GetAllTasksResponse{
		StatusCode: domain.Success,
		Tasks:      tasks,
	}, nil
}

func (i *TasksInteractor) CreateTask(name string) (domain.CreateTaskResponse, error) {
	task := &domain.Task{
		Name: name,
	}
	err := i.tasksRepo.Insert(task)
	if err != nil {
		return domain.CreateTaskResponse{}, errors.Wrap(err, "Cannot insert task")
	}

	return domain.CreateTaskResponse{
		StatusCode: domain.Success,
		Id:         task.Id,
	}, nil
}
