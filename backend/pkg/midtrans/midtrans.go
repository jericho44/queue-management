package midtrans

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type MidtransClient struct {
	ServerKey string
	IsSandbox bool
}

func NewMidtransClient() *MidtransClient {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if serverKey == "" {
		serverKey = "SB-Mid-server-DummyKeyForDevelopment"
	}
	isSandbox := os.Getenv("MIDTRANS_IS_SANDBOX") != "false"

	return &MidtransClient{
		ServerKey: serverKey,
		IsSandbox: isSandbox,
	}
}

func (m *MidtransClient) getSnapURL() string {
	if m.IsSandbox {
		return "https://app.sandbox.midtrans.com/snap/v1/transactions"
	}
	return "https://app.midtrans.com/snap/v1/transactions"
}

type SnapTransactionRequest struct {
	TransactionDetails struct {
		OrderID     string `json:"order_id"`
		GrossAmount int64  `json:"gross_amount"`
	} `json:"transaction_details"`
	CustomerDetails struct {
		FirstName string `json:"first_name"`
		Email     string `json:"email"`
	} `json:"customer_details"`
	ItemDetails []struct {
		ID       string `json:"id"`
		Price    int64  `json:"price"`
		Quantity int    `json:"quantity"`
		Name     string `json:"name"`
	} `json:"item_details"`
	EnabledPayments []string `json:"enabled_payments,omitempty"`
}

type SnapTransactionResponse struct {
	Token         string   `json:"token"`
	RedirectURL   string   `json:"redirect_url"`
	ErrorMessages []string `json:"error_messages,omitempty"`
}

func (m *MidtransClient) CreateSnapToken(orderID string, amount int64, orgName string, email string, description string) (*SnapTransactionResponse, error) {
	reqBody := SnapTransactionRequest{}
	reqBody.TransactionDetails.OrderID = orderID
	reqBody.TransactionDetails.GrossAmount = amount

	reqBody.CustomerDetails.FirstName = orgName
	reqBody.CustomerDetails.Email = email

	// Daftar kanal pembayaran yang diaktifkan (Opsional: dapat disesuaikan)
	reqBody.EnabledPayments = []string{
		"credit_card",
		"gopay",
		"shopeepay",
		"qris",
		"bca_va",
		"bni_va",
		"bri_va",
		"permata_va",
		"other_va",
		"indomaret",
		"alfamart",
	}

	reqBody.ItemDetails = []struct {
		ID       string `json:"id"`
		Price    int64  `json:"price"`
		Quantity int    `json:"quantity"`
		Name     string `json:"name"`
	}{
		{
			ID:       orderID,
			Price:    amount,
			Quantity: 1,
			Name:     description,
		},
	}


	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", m.getSnapURL(), bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(m.ServerKey + ":"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+authHeader)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call Midtrans API: %w", err)
	}
	defer resp.Body.Close()

	var snapResp SnapTransactionResponse
	if err := json.NewDecoder(resp.Body).Decode(&snapResp); err != nil {
		return nil, fmt.Errorf("failed to parse Midtrans response: %w", err)
	}

	return &snapResp, nil
}

func (m *MidtransClient) VerifySignature(orderID, statusCode, grossAmount, signatureKey string) bool {
	raw := orderID + statusCode + grossAmount + m.ServerKey
	hash := sha512.Sum512([]byte(raw))
	expectedSig := hex.EncodeToString(hash[:])
	return expectedSig == signatureKey
}
