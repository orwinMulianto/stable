package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type BrevoSender struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type BrevoRecipient struct {
	Email string `json:"email"`
}

type BrevoRequest struct {
	Sender      BrevoSender      `json:"sender"`
	To          []BrevoRecipient `json:"to"`
	Subject     string           `json:"subject"`
	HTMLContent string           `json:"htmlContent"`
}

type BrevoResponse struct {
	MessageID string `json:"messageId"`
	Code      string `json:"code"`
	Message   string `json:"message"`
}

func SendEmail(to string, subject string, body string) error {
	log.Println("SendEmail called")
	log.Println("TO:", to)
	log.Println("SUBJECT:", subject)
	apiKey := os.Getenv("BREVO_API_KEY")
	senderEmail := os.Getenv("BREVO_SENDER_EMAIL")
	senderName := os.Getenv("BREVO_SENDER_NAME")

	if senderEmail == "" {
		senderEmail = "noreply@stable.com"
	}
	if senderName == "" {
		senderName = "Stable"
	}

	log.Printf("[EMAIL] Attempting to send email to: %s via Brevo", to)
	log.Printf("[EMAIL] BREVO_API_KEY present: %v, length: %d", apiKey != "", len(apiKey))

	if apiKey == "" {
		log.Printf("[EMAIL] ERROR: BREVO_API_KEY is not set")
		return fmt.Errorf("BREVO_API_KEY is not set")
	}

	htmlBody := fmt.Sprintf("<p>%s</p>", body)

	reqBody := BrevoRequest{
		Sender: BrevoSender{
			Email: senderEmail,
			Name:  senderName,
		},
		To: []BrevoRecipient{
			{Email: to},
		},
		Subject:     subject,
		HTMLContent: htmlBody,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("[EMAIL] ERROR marshaling request: %v", err)
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Printf("[EMAIL] Request body: %s", string(jsonBody))

	// Create HTTP request to Brevo API
	req, err := http.NewRequest("POST", "https://api.brevo.com/v3/smtp/email", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("[EMAIL] ERROR creating request: %v", err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[EMAIL] ERROR sending request: %v", err)
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[EMAIL] ERROR reading response: %v", err)
		return fmt.Errorf("failed to read response: %w", err)
	}

	log.Printf("[EMAIL] Response status: %d, body: %s", resp.StatusCode, string(respBody))

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("[EMAIL] SUCCESS: Email sent to %s via Brevo", to)
		return nil
	}

	var brevoResp BrevoResponse
	if err := json.Unmarshal(respBody, &brevoResp); err != nil {
		log.Printf("[EMAIL] ERROR parsing response: %v", err)
		return fmt.Errorf("Brevo error: status %d", resp.StatusCode)
	}

	log.Printf("[EMAIL] ERROR from Brevo: %s - %s", brevoResp.Code, brevoResp.Message)
	return fmt.Errorf("Brevo error: %s - %s", brevoResp.Code, brevoResp.Message)
}

func GenerateVerificationCode() string {
	n := rand.Intn(900000)
	code := int64(n) + 100000
	return strconv.FormatInt(code, 10)
}