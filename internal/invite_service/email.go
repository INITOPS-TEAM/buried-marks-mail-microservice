package invite_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type BrevoRequest struct {
	Sender      map[string]string   `json:"sender"`
	To          []map[string]string `json:"to"`
	Subject     string              `json:"subject"`
	HtmlContent string              `json:"htmlContent"`
}

func sendInviteEmail(toEmail string, inviteLink string) {
	apiKey := os.Getenv("BREVO_API_KEY")
	senderEmail := os.Getenv("SENDER_EMAIL")

	if apiKey == "" || senderEmail == "" {
		slog.Error("Missing required env variables for email", "email", toEmail)
		return
	}

	htmlBody := fmt.Sprintf(`
		<h3>Congratulations!</h3>
		<p>You have been invited to join our system.</p>
		<p>Click on the link below to create an account::</p>
		<a href="%s" style="display:inline-block; padding:10px 20px; background:#007bff; color:#fff; text-decoration:none; border-radius:5px;">
			Sign up
		</a>
		<p>Or copy this link into your browser: <br> %s</p>
	`, inviteLink, inviteLink)

	reqBody := BrevoRequest{
		Sender:      map[string]string{"email": senderEmail, "name": "Secret Society"},
		To:          []map[string]string{{"email": toEmail}},
		Subject:     "Invitation to register",
		HtmlContent: htmlBody,
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "https://api.brevo.com/v3/smtp/email", bytes.NewBuffer(jsonData))

	req.Header.Set("accept", "application/json")
	req.Header.Set("api-key", apiKey)
	req.Header.Set("content-type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Network error while sending invite email", "err", err, "email", toEmail)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 || resp.StatusCode == 200 {
		slog.Info("Invite email successfully sent", "email", toEmail)
	} else {
		slog.Error("brevo server rejected request", "status", resp.StatusCode, "email", toEmail)
	}
}
