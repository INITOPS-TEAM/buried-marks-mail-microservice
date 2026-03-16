package daily_code_service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/valkey-io/valkey-go"
)

func GenerateToken() string {
	return uuid.New().String()
}

func ProcessDailyCode(valkeyClient valkey.Client) {
	slog.Info("starting daily code generation process")

	dailyCode := GenerateToken()
	key := os.Getenv("DAILY_CODE_KEY")
	ctx := context.Background()

	err := valkeyClient.Do(ctx, valkeyClient.B().Set().Key(key).Value(dailyCode).Ex(24*time.Hour).Build()).Error()
	if err != nil {
		slog.Error("failed to save daily code to valkey", "err", err)
		return
	}
	slog.Info("daily code successfully updated in valkey")

	var emails []string
	var fetchErr error
	maxRetries := 5

	for attempt := 1; attempt <= maxRetries; attempt++ {
		emails, fetchErr = fetchEmailsFromAuth()
		if fetchErr == nil {
			break
		}
		time.Sleep(3 * time.Second)
	}

	if fetchErr != nil {
		slog.Error("failed to fetch emails from auth service after all retries", "err", fetchErr)
		return
	}

	slog.Info("successfully fetched emails", "count", len(emails))

	if len(emails) == 0 {
		slog.Warn("no emails found to send the daily code to")
		return
	}

	for _, email := range emails {
		go SendEmail(email, dailyCode)
	}
}

type BrevoRequest struct {
	Sender      map[string]string   `json:"sender"`
	To          []map[string]string `json:"to"`
	Subject     string              `json:"subject"`
	HtmlContent string              `json:"htmlContent"`
}

func SendEmail(toEmail string, code string) {
	apiKey := os.Getenv("BREVO_API_KEY")
	senderEmail := os.Getenv("SENDER_EMAIL")

	if apiKey == "" || senderEmail == "" {
		slog.Error("missing required env variables for email", "email", toEmail)
		return
	}

	reqBody := BrevoRequest{
		Sender:      map[string]string{"email": senderEmail, "name": "Buried Marks"},
		To:          []map[string]string{{"email": toEmail}},
		Subject:     "Buried Marks: Your Daily Code",
		HtmlContent: fmt.Sprintf("<h3>Hello!</h3><p>Your new daily code is: <b>%s</b></p>", code),
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
		slog.Error("network error while sending email", "err", err, "email", toEmail)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 || resp.StatusCode == 200 {
		slog.Info("email successfully sent", "email", toEmail)
	} else {
		bodyBytes, _ := io.ReadAll(resp.Body)
		slog.Error("brevo server rejected request", "status", resp.StatusCode, "response", string(bodyBytes), "email", toEmail)
	}
}
