package cron

import (
	"context"
	"time"

	"github.com/dmmitrenko/weather-app/internal/application"
	"github.com/dmmitrenko/weather-app/internal/domain"
	"github.com/robfig/cron/v3"
)

func StartJobs(proc *application.SubscriptionProcessor) *cron.Cron {
	scheduler := cron.New(
		cron.WithLocation(time.FixedZone("Europe/Kyiv", 3*3600)),
		cron.WithChain(cron.Recover(cron.DefaultLogger)),
	)

	scheduler.AddFunc("@hourly", func() {
		proc.Process(context.Background(), domain.Hourly)
	})

	scheduler.AddFunc("0 0 * * *", func() {
		proc.Process(context.Background(), domain.Daily)
	})
	scheduler.Start()
	return scheduler
}
