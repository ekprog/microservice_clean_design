package bootstrap

import (
	"microservice_clean_design/app/job"
	"microservice_clean_design/jobs"
)

func initJobs() error {
	job.NewJob("test1", jobs.TestJob, job.Time("1 second"))
	return nil
}
