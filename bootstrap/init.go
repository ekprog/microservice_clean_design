package bootstrap

import (
	"database/sql"
	"microservice_clean_design/app"
	"microservice_clean_design/delivery"
	"microservice_clean_design/interactors"
	"microservice_clean_design/repos"
)

func injectDependencies(db *sql.DB, logger app.Logger) error {

	// DI Auto example
	//diObj := dig.New()

	// logger
	//if err := diObj.Provide(func() app.Logger {
	//	return logger
	//}); err != nil {
	//	return err
	//}
	//
	//// db
	//if err := diObj.Provide(func() *sql.DB {
	//	return db
	//}); err != nil {
	//	return err
	//}

	//
	//err = provide(diObj,
	//	repos.NewUsersRepo,
	//	repos.NewUserTokensRepo,
	//	interactors.NewTasksUCase,
	//	interactors.NewTasksUCase,
	//	grpc_delivery.NewAuthDeliveryService)
	//if err != nil {
	//	return err
	//}

	// DI Manual
	tasksRepo := repos.NewTaskDBRepo(logger, db)
	tasksUCase := interactors.NewTasksUCase(logger, tasksRepo)

	// Delivery Init
	tasksDelivery := delivery.NewTasksDeliveryService(logger, tasksUCase)

	err := app.InitDelivery(tasksDelivery)
	if err != nil {
		return err
	}
	return nil
}

//func provide(diObj *dig.Container, list ...interface{}) error {
//	for _, p := range list {
//		if err := diObj.Provide(p); err != nil {
//			return err
//		}
//	}
//	return nil
//}
