package bootstrap

import (
	"database/sql"
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"microservice/app"
	"microservice/app/core"
	"microservice/app/job"
	"microservice/app/kafka"
)

func Run(rootPath ...string) error {

	// ENV, etc
	err := app.InitApp(rootPath...)
	if err != nil {
		return errors.Wrap(err, "error while init app")
	}

	// Logger
	logger, err := app.InitLogs(rootPath...)
	if err != nil {
		return errors.Wrap(err, "error while init logs")
	}

	// Storage
	err = app.InitStorage()
	if err != nil {
		return errors.Wrap(err, "error while init storage")
	}

	// Database
	db, err := app.InitDatabase()
	if err != nil {
		return errors.Wrap(err, "error while init db")
	}

	// Migrations
	err = app.RunMigrations(rootPath...)
	if err != nil {
		return errors.Wrap(err, "error while making migrations")
	}

	// gRPC
	_, _, err = app.InitGRPCServer()
	if err != nil {
		return errors.Wrap(err, "cannot init gRPC")
	}

	// DI
	di := dig.New()

	if err = di.Provide(func() *sql.DB {
		return db
	}); err != nil {
		return errors.Wrap(err, "cannot provide db")
	}

	if err = di.Provide(func() core.Logger {
		return logger
	}); err != nil {
		return errors.Wrap(err, "cannot provide logger")
	}

	// CRON
	err = job.Init(logger, di)
	if err != nil {
		return errors.Wrap(err, "cannot init jobs")
	}

	// KAFKA
	err = kafka.InitKafka(logger)
	if err != nil {
		return errors.Wrap(err, "cannot init kafka")
	}

	// CORE
	if err := initDependencies(di); err != nil {
		return errors.Wrap(err, "error while init dependencies")
	}

	//
	//
	// HERE CORE READY FOR WORK...
	//
	//

	// CRON
	if err := initJobs(); err != nil {
		return errors.Wrap(err, "error while init jobs")
	}

	if err := job.Start(); err != nil {
		return errors.Wrap(err, "error while start jobs")
	}

	// Run gRPC and block
	app.RunGRPCServer()

	return nil
}
