package daily_code_service

import (
	"log/slog"

	"github.com/go-co-op/gocron/v2"
	"github.com/valkey-io/valkey-go"
)

func StartWorker(valkeyClient valkey.Client) error {
	s, err := gocron.NewScheduler()
	if err != nil {
		return err
	}

	_, err = s.NewJob(
		//Testing every 60 second code generation
		//gocron.DurationJob(60*time.Second),
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(gocron.NewAtTime(0, 0, 0)),
		),
		gocron.NewTask(
			func() {
				ProcessDailyCode(valkeyClient)
			},
		),
	)

	if err != nil {
		return err
	}

	s.Start()
	slog.Info("worker scheduler successfully started")

	return nil
}
