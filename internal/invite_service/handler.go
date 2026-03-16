package invite_service

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/INITOPS-TEAM/buried-marks-mail-microservice/internal/daily_code_service"
	"github.com/valkey-io/valkey-go"
)

type InviteRequest struct {
	Email      string `json:"email"`
	InviteLink string `json:"invite_link"`
}

func HandleSendInvite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not Post", http.StatusMethodNotAllowed)
		return
	}

	var reqData InviteRequest
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		slog.Error("Failed to decode invite request", "err", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	go sendInviteEmail(reqData.Email, reqData.InviteLink)

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status": "processing"}`))
}

type WelcomeCodeRequest struct {
	Email string `json:"email"`
}

func HandleSendInviteDailyCode(valkeyClient valkey.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not Post", http.StatusMethodNotAllowed)
			return
		}

		var reqData WelcomeCodeRequest
		if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
			slog.Error("Failed to decode welcome code request", "err", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		ctx := context.Background()
		key := os.Getenv("DAILY_CODE_KEY")
		code, err := valkeyClient.Do(ctx, valkeyClient.B().Get().Key(key).Build()).ToString()
		if err != nil || code == "" {
			slog.Error("Failed to retrieve daily code from valkey for new user", "err", err)
			http.Error(w, "Daily code not available", http.StatusInternalServerError)
			return
		}

		go daily_code_service.SendEmail(reqData.Email, code)

		slog.Info("Triggered welcome daily code email", "email", reqData.Email)

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"status": "processing"}`))
	}
}
