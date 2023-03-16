package job

import (
	"github.com/jasonlvhit/gocron"
	"github.com/spf13/viper"
	"go.uber.org/dig"
	"microservice/app/core"
	"strconv"
	"strings"
)

var log core.Logger
var di *dig.Container

func Init(logger core.Logger, di_ *dig.Container) error {
	log = logger
	di = di_
	return nil
}

func NewJob(name string, fn interface{}, job *gocron.Job) {
	// DO
	err := job.Do(func() {
		err := di.Invoke(fn)
		if err != nil {
			log.ErrorWrap(err, "error in CRON %s", name)
		}
	})
	if err != nil {
		log.Fatal("cannot DO cron %s", name)
	}
	log.Info("New job was successfully registered - %s", name)
}

func Start() error {

	enabled := viper.GetBool("jobs.enabled")
	if !enabled {
		return nil
	}

	go gocron.Start()
	return nil
}

func Time(time string) *gocron.Job {

	split := strings.Split(time, " ")
	if len(split) != 2 {
		return nil
	}
	timeVal, err := strconv.ParseUint(split[0], 10, 64)
	if err != nil {
		return nil
	}
	timeId := split[1]

	job := gocron.Every(timeVal)
	switch timeId {
	case "seconds":
		job = job.Seconds()
	case "second":
		job = job.Second()
	case "minute":
		job = job.Minute()
	case "minutes":
		job = job.Minutes()
	case "hour":
		job = job.Hour()
	case "hours":
		job = job.Hours()
	}

	return job
}
