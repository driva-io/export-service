package adapters

import (
	"bytes"
	"encoding/json"
	"export-service/internal/core/ports"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

type DrivaMailer struct {
	logger          *zap.Logger
	adminToken      string
	tokenExpiration time.Time
}

var _ ports.Mailer = (*DrivaMailer)(nil)

func NewDrivaMailer(logger *zap.Logger) *DrivaMailer {
	return &DrivaMailer{
		logger: logger,
	}
}

func (d *DrivaMailer) SendEmail(userEmail, userName, templateId, link string) error {
	payload := map[string]any{
		"from":        "dados@driva.io",
		"to":          userEmail,
		"template_id": templateId,
		"dynamicTemplateData": map[string]string{
			"link": link,
			"name": userName,
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://services.driva.io/automation/v1/emails/sendEmail", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	token, err := d.getAdminToken()
	if err != nil {
		d.logger.Error("Failed to get admin token", zap.Error(err))
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to send email, status: %d", resp.StatusCode)
	}

	return nil
}

func (d *DrivaMailer) getAdminToken() (string, error) {
	if d.adminToken != "" && d.tokenExpiration.After(time.Now()) {
		return d.adminToken, nil
	}
	d.logger.Info("Cache miss, fetching admin token")

	payload := map[string]any{"email": os.Getenv("ADMIN_EMAIL"), "password": os.Getenv("ADMIN_PASSWORD")}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://services.driva.io/auth/v1/authenticate", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("company", "Driva")
	req.URL.RawQuery = q.Encode()

	client := &http.Client{Timeout: 15 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("failed to authenticate, status: %d", resp.StatusCode)
	}

	var respData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", err
	}

	token, ok := respData["token"].(string)
	if !ok {
		return "", fmt.Errorf("token not found in response")
	}

	d.adminToken = token
	d.tokenExpiration = time.Now().Add(1 * time.Hour)

	return token, nil
}
