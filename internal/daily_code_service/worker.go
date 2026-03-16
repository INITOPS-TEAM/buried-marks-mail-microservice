package daily_code_service

import (
	"context"
	"log/slog"
	"os"

	"github.com/go-co-op/gocron/v2"
	"github.com/valkey-io/valkey-go"
)

func StartWorker(valkeyClient valkey.Client) error {
	ctx := context.Background()
	key := os.Getenv("DATABASE_KEY")
	exists, err := valkeyClient.Do(ctx, valkeyClient.B().Exists().Key(key).Build()).AsInt64()

	if err != nil {
		return err
	} else if exists == 0 {
		slog.Info("daily code not found. Generating a new one")
		ProcessDailyCode(valkeyClient)
	} else {
		slog.Info("daily code exists")
	}

	s, err := gocron.NewScheduler()
	if err != nil {
		return err
	}

	_, err = s.NewJob(
		//Testing every 120 second code generation
		//gocron.DurationJob(120*time.Second),
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
