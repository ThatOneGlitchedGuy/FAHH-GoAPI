package service

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

var (
	Scheduler *gocron.Scheduler
)

func InitScheduler() {
	Scheduler = gocron.NewScheduler(time.UTC)
	Scheduler.StartAsync()
}

func StartJobs() {
	_, err := Scheduler.Every(1).Minute().Do(func() {
		log.Println("This is a scheduled job running every minute.")
	})
	if err != nil {
		log.Printf("Failed to schedule job: %v", err)
	}
}
