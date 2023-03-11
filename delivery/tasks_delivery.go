package delivery

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"microservice_clean_design/app"
	"microservice_clean_design/app/core"
	"microservice_clean_design/domain"
	pb "microservice_clean_design/pkg/pb/api"
)

type TasksDeliveryService struct {
	pb.TasksServiceServer
	log        core.Logger
	tasksUCase domain.TasksInteractor
}

func NewTasksDeliveryService(log core.Logger, tasksUCase domain.TasksInteractor) *TasksDeliveryService {
	return &TasksDeliveryService{
		log:        log,
		tasksUCase: tasksUCase,
	}
}

func (d *TasksDeliveryService) Init() error {
	app.InitGRPCService(pb.RegisterTasksServiceServer, pb.TasksServiceServer(d))
	return nil
}

func (d *TasksDeliveryService) AllTasks(ctx context.Context, r *pb.AllTasksRequest) (*pb.AllTasksResponse, error) {

	uCaseRes, err := d.tasksUCase.GetAllTasks()
	if err != nil {
		return nil, err
	}

	response := &pb.AllTasksResponse{
		Status: &pb.StatusResponse{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
	}

	if uCaseRes.StatusCode == domain.Success {

		var tasks []*pb.Task
		for _, task := range uCaseRes.Tasks {
			tasks = append(tasks, &pb.Task{
				Id:        task.Id,
				Name:      task.Name,
				CreatedAt: timestamppb.New(task.CreatedAt),
				UpdatedAt: timestamppb.New(task.UpdatedAt),
			})
		}
		response.Tasks = tasks
	}

	return response, nil
}

func (d *TasksDeliveryService) CreateTask(ctx context.Context, r *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {

	uCaseRes, err := d.tasksUCase.CreateTask(r.Name)
	if err != nil {
		return nil, err
	}

	response := &pb.CreateTaskResponse{
		Status: &pb.StatusResponse{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
		Id: uCaseRes.Id,
	}

	return response, nil
}
