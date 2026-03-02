package daily_code_service

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"
)

type EmailsResponse struct {
	Emails []string `json:"emails"`
}

// Getting mails list from auth service
func fetchEmailsFromAuth() ([]string, error) {
	url := os.Getenv("SERVICE_URL")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	secretKey := os.Getenv("DJANGO_SECRET_KEY")

	req.Header.Set("X-Internal-Token", secretKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	var result EmailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Emails, nil
}
