package store

import (
	"fmt"
	"log/slog"

	"github.com/valkey-io/valkey-go"
)

func InitValkey(addr string) (valkey.Client, error) {
	if addr == "" {
		return nil, fmt.Errorf("VALKEY_ADDR variable is required")
	}

	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{addr},
	})

	if err != nil {
		slog.Error("valkey connection failed", "err", err)
		return nil, err
	}

	return client, nil
}
