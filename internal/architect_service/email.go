package architect_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func sendArchitectEmail(toEmail string, subject string, customText string) {
	apiKey := os.Getenv("BREVO_API_KEY")
	senderEmail := os.Getenv("SENDER_EMAIL")

	if apiKey == "" || senderEmail == "" {
		slog.Error("missing required env variables for email", "email", toEmail)
		return
	}

	reqBody := BrevoRequest{
		Sender:      map[string]string{"email": senderEmail, "name": "Buried Marks"},
		To:          []map[string]string{{"email": toEmail}},
		Subject:     subject,
		HtmlContent: fmt.Sprintf("<p>%s</p>", customText),
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "https://api.brevo.com/v3/smtp/email", bytes.NewBuffer(jsonData))

	req.Header.Set("accept", "application/json")
	req.Header.Set("api-key", apiKey)
	req.Header.Set("content-type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)

	if err != nil {
		slog.Error("network error while sending architect email", "err", err, "email", toEmail)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 || resp.StatusCode == 200 || resp.StatusCode == 202 {
		slog.Info("architect email successfully sent", "email", toEmail)
	} else {
		bodyBytes, _ := io.ReadAll(resp.Body)
		slog.Error("brevo server rejected architect request", "status", resp.StatusCode, "response", string(bodyBytes), "email", toEmail)
	}
}
