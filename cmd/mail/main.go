package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/INITOPS-TEAM/buried-marks-mail-microservice/internal/daily_code_service"
	"github.com/INITOPS-TEAM/buried-marks-mail-microservice/internal/store"
)

func main() {
	valkeyAddr := os.Getenv("VALKEY_ADDR")

	client, err := store.InitValkey(valkeyAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Valkey: %v", err)
	}
	defer client.Close()

	if err := daily_code_service.StartWorker(client); err != nil {
		log.Fatalf("Failed to start worker: %v", err)
	}

	log.Println("Mail service is running...")

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan

	log.Println("Server successfully shut down")
}
