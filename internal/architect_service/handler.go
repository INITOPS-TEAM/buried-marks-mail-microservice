package architect_service

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ArchitectEmailRequest struct {
	Emails     []string `json:"emails"`
	Subject    string   `json:"subject"`
	CustomText string   `json:"custom_text"`
}

func HandleSendArchitectEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqData ArchitectEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		slog.Error("Failed to decode architect email request", "err", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	go func(req ArchitectEmailRequest) {
		slog.Info("Starting architect mailing", "target_count", len(req.Emails))
		for _, email := range req.Emails {
			go sendArchitectEmail(email, req.Subject, req.CustomText)
		}
	}(reqData)

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status": "architect_mailing_started"}`))
}
