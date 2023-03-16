package bootstrap

import (
	"microservice/app/job"
	"microservice/jobs"
)

func initJobs() error {
	job.NewJob("test1", jobs.TestJob, job.Time("1 second"))
	return nil
}
