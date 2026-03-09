package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/INITOPS-TEAM/buried-marks-mail-microservice/internal/architect_service"
	"github.com/INITOPS-TEAM/buried-marks-mail-microservice/internal/daily_code_service"
	"github.com/INITOPS-TEAM/buried-marks-mail-microservice/internal/invite_service"
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

	http.HandleFunc("/api/send-invite", invite_service.HandleSendInvite)
	http.HandleFunc("/api/send-architect-email", architect_service.HandleSendArchitectEmail)

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("HTTP server failed %v", err)
		}
	}()

	log.Println("Mail service is running...")

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan

	log.Println("Server successfully shut down")
}
