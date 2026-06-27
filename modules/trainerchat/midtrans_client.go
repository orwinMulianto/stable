package trainerchat

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type midtransClient struct {
	serverKey    string
	isProduction bool
	httpClient   *http.Client
}

type snapTransactionRequest struct {
	TransactionDetails snapTransactionDetails `json:"transaction_details"`
	ItemDetails        []snapItemDetail       `json:"item_details,omitempty"`
	CustomerDetails    snapCustomerDetails    `json:"customer_details,omitempty"`
	Callbacks          snapCallbacks          `json:"callbacks,omitempty"`
}

type snapTransactionDetails struct {
	OrderID     string `json:"order_id"`
	GrossAmount int    `json:"gross_amount"`
}

type snapItemDetail struct {
	ID       string `json:"id"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
	Name     string `json:"name"`
}

type snapCustomerDetails struct {
	FirstName string `json:"first_name,omitempty"`
	Email     string `json:"email,omitempty"`
}

type snapCallbacks struct {
	Finish string `json:"finish,omitempty"`
}

type snapTransactionResponse struct {
	Token         string   `json:"token"`
	RedirectURL   string   `json:"redirect_url"`
	ErrorMessages []string `json:"error_messages"`
}

func newMidtransClient() *midtransClient {
	return &midtransClient{
		serverKey:    os.Getenv("MIDTRANS_SERVER_KEY"),
		isProduction: strings.EqualFold(os.Getenv("MIDTRANS_IS_PRODUCTION"), "true"),
		httpClient:   &http.Client{Timeout: 12 * time.Second},
	}
}

func (client *midtransClient) snapEndpoint() string {
	if client.isProduction {
		return "https://app.midtrans.com/snap/v1/transactions"
	}

	return "https://app.sandbox.midtrans.com/snap/v1/transactions"
}

func (client *midtransClient) CreateSnapTransaction(
	sessionID uint,
	orderID string,
	trainer TrainerCatalogItem,
	request CheckoutRequest,
) (snapTransactionResponse, error) {
	if strings.TrimSpace(client.serverKey) == "" {
		return snapTransactionResponse{}, errors.New("MIDTRANS_SERVER_KEY is required")
	}

	payload := snapTransactionRequest{
		TransactionDetails: snapTransactionDetails{
			OrderID:     orderID,
			GrossAmount: trainer.Price,
		},
		ItemDetails: []snapItemDetail{
			{
				ID:       fmt.Sprintf("trainer-%d", trainer.ID),
				Price:    trainer.Price,
				Quantity: 1,
				Name:     fmt.Sprintf("Chat 10 menit - %s", trainer.Name),
			},
		},
		CustomerDetails: snapCustomerDetails{
			FirstName: request.CustomerName,
			Email:     request.CustomerEmail,
		},
		Callbacks: snapCallbacks{
			Finish: os.Getenv("MIDTRANS_FINISH_URL"),
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return snapTransactionResponse{}, err
	}

	httpRequest, err := http.NewRequest(http.MethodPost, client.snapEndpoint(), bytes.NewReader(body))
	if err != nil {
		return snapTransactionResponse{}, err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(client.serverKey + ":"))
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Basic "+auth)

	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return snapTransactionResponse{}, err
	}
	defer httpResponse.Body.Close()

	var response snapTransactionResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return snapTransactionResponse{}, err
	}

	if httpResponse.StatusCode < 200 || httpResponse.StatusCode >= 300 {
		if len(response.ErrorMessages) > 0 {
			return snapTransactionResponse{}, errors.New(strings.Join(response.ErrorMessages, ", "))
		}
		return snapTransactionResponse{}, fmt.Errorf("midtrans error status %d", httpResponse.StatusCode)
	}

	if response.Token == "" {
		return snapTransactionResponse{}, errors.New("midtrans snap token is empty")
	}
	fmt.Printf("DEBUG serverKey: %q\n", client.serverKey[:min(8, len(client.serverKey))])
    fmt.Printf("DEBUG endpoint: %s\n", client.snapEndpoint())
    fmt.Printf("DEBUG payload: %s\n", body) // setelah json.Marshal

	_ = sessionID
	return response, nil
	
}

func (client *midtransClient) VerifyNotificationSignature(notification MidtransNotification) bool {
	if strings.TrimSpace(client.serverKey) == "" {
		return false
	}

	raw := notification.OrderID + notification.StatusCode + notification.GrossAmount + client.serverKey
	sum := sha512.Sum512([]byte(raw))
	expected := hex.EncodeToString(sum[:])

	return strings.EqualFold(expected, notification.SignatureKey)
}
