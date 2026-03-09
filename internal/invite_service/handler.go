package invite_service

import (
	"encoding/json"
	"log/slog"
	"net/http"
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
