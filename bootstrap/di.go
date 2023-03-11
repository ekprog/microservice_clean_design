package bootstrap

import (
	"go.uber.org/dig"
	"microservice_clean_design/app"
	"microservice_clean_design/delivery"
	"microservice_clean_design/domain"
	"microservice_clean_design/interactors"
	"microservice_clean_design/repos"
)

func initDependencies(di *dig.Container) error {

	di.Provide(repos.NewTaskDBRepo, dig.As(new(domain.TasksRepository)))

	di.Provide(interactors.NewTasksUCase, dig.As(new(domain.TasksInteractor)))

	di.Provide(delivery.NewTasksDeliveryService)

	// DELIVERY
	deliveryInit := func(tasks *delivery.TasksDeliveryService) error {
		if err := app.InitDelivery(tasks); err == nil {
			return err
		}
		return nil
	}

	err := di.Invoke(deliveryInit)
	if err != nil {
		return err
	}

	return nil
}
