package investecOpenAPI

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/t4ke0/investecOpenAPI/api"
)

var APIurl string = "https://openapi.investec.com"

var IsDebug bool

type BankingClient struct {
	UserCreds   string
	AccessToken string

	httpClient *http.Client
}

func NewBankingClient(key, secret, clientID string) BankingClient {
	userCreds := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientID, secret)))
	return BankingClient{
		UserCreds:   userCreds,
		AccessToken: key,
		httpClient:  new(http.Client),
	}
}

type authMode int32

const (
	basic authMode = iota
	bearer
)

func (b BankingClient) requestAPI(url, method string, mode authMode, body []byte) (*http.Response, error) {
	var (
		req *http.Request
		err error
	)
	if mode == basic {
		req, err = http.NewRequest(method, url, strings.NewReader("grant_type=client_credentials&scope=accounts"))
	} else {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	}

	if err != nil {
		return nil, err
	}

	switch mode {
	case basic:
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", b.UserCreds))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("x-api-key", b.AccessToken)
	case bearer:
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", b.AccessToken))
		req.Header.Set("x-api-key", b.AccessToken)
		req.Header.Set("Accept", "application/json")
	}

	return b.httpClient.Do(req)
}

// GetAccessToken
func (b *BankingClient) GetAccessToken() error {
	if IsDebug {
		APIurl = "https://openapisandbox.investec.com"
	}
	url := fmt.Sprintf("%s/identity/v2/oauth2/token", APIurl)
	resp, err := b.requestAPI(url, http.MethodPost, basic, nil)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get access token oauth [%v] [%v]", resp.StatusCode, string(data))
	}

	var tokenResponse api.OAuth2Response

	if err := json.Unmarshal(data, &tokenResponse); err != nil {
		return err
	}

	b.AccessToken = tokenResponse.AccessToken
	return nil
}

// GetAccounts
func (b BankingClient) GetAccounts() (api.Accounts, error) {
	if IsDebug {
		APIurl = "https://openapisandbox.investec.com"
	}
	url := fmt.Sprintf("%s/za/pb/v1/accounts", APIurl)
	resp, err := b.requestAPI(url, http.MethodGet, bearer, nil)
	if err != nil {
		return api.Accounts{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return api.Accounts{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return api.Accounts{}, fmt.Errorf("failed to get accounts %v %v", resp.StatusCode, string(data))
	}

	var accounts api.Accounts

	return accounts, json.Unmarshal(data, &accounts)
}

// GetAccountBalance
func (b BankingClient) GetAccountBalance(accountId string) (api.Balance, error) {
	if IsDebug {
		APIurl = "https://openapisandbox.investec.com"
	}
	url := fmt.Sprintf("%s/za/pb/v1/accounts/%s/balance", APIurl, accountId)
	resp, err := b.requestAPI(url, http.MethodGet, bearer, nil)
	if err != nil {
		return api.Balance{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return api.Balance{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return api.Balance{}, fmt.Errorf("failed to get balance %v %v", resp.StatusCode, string(data))
	}

	var balance api.Balance

	return balance, json.Unmarshal(data, &balance)
}

// GetAccountTransactions gets account transactions you can filter the result
// using the fromDate and toDate parameters and the dates need to follow ISO
// 8601 time format
func (b BankingClient) GetAccountTransactions(accountID string, fromDate, toDate string) (api.Transactions, error) {
	if IsDebug {
		APIurl = "https://openapisandbox.investec.com"
	}
	var url string = fmt.Sprintf("%s/za/pb/v1/accounts/%s/transactions", APIurl, accountID)
	if fromDate != "" && toDate != "" {
		url += fmt.Sprintf("?fromDate=%s&toDate=%s", fromDate, toDate)
	}

	resp, err := b.requestAPI(url, http.MethodGet, bearer, nil)
	if err != nil {
		return api.Transactions{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return api.Transactions{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return api.Transactions{}, fmt.Errorf("failed to get account Transactions %v %v", resp.StatusCode, string(data))
	}

	var transactions api.Transactions
	return transactions, json.Unmarshal(data, &transactions)
}

// TransferMultiple
func (b BankingClient) TransferMultiple(accountID string, transfers []api.TransferTo) (api.MultipleTransfersResponse, error) {
	if IsDebug {
		APIurl = "https://openapisandbox.investec.com"
	}
	var url string = fmt.Sprintf("%s/za/pb/v1/accounts/%s/transfermultiple", APIurl, accountID)

	var requestBodyObj = struct {
		TransferList []api.TransferTo `json:"transferList"`
	}{transfers}

	data, err := json.Marshal(requestBodyObj)
	if err != nil {
		return api.MultipleTransfersResponse{}, err
	}
	resp, err := b.requestAPI(url, http.MethodPost, bearer, data)
	if err != nil {
		return api.MultipleTransfersResponse{}, err
	}
	defer resp.Body.Close()

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return api.MultipleTransfersResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return api.MultipleTransfersResponse{}, fmt.Errorf(
			"Error failed to TransferMultiple [%v] [%v]", resp.StatusCode, string(data),
		)
	}

	var response api.MultipleTransfersResponse

	if IsDebug {
		log.Printf("DEBUG RESPONSE TRANSFER MULTIPLE %v %v", resp.StatusCode, string(data))
		return response, nil
	}

	return response, json.Unmarshal(data, &response)
}

// PayMultiple
func (b BankingClient) PayMultiple(accountID string, payments []api.PaymentMultiple) (api.MultiplePaymentResponse, error) {
	if IsDebug {
		APIurl = "https://openapisandbox.investec.com"
	}
	var url string = fmt.Sprintf("%s/za/pb/v1/accounts/%s/paymultiple", APIurl, accountID)

	var requestBodyObj = struct {
		PaymentList []api.PaymentMultiple `json:"paymentList"`
	}{payments}

	data, err := json.Marshal(requestBodyObj)
	if err != nil {
		return api.MultiplePaymentResponse{}, err
	}

	resp, err := b.requestAPI(url, http.MethodPost, bearer, data)
	if err != nil {
		return api.MultiplePaymentResponse{}, err
	}
	defer resp.Body.Close()

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return api.MultiplePaymentResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return api.MultiplePaymentResponse{}, fmt.Errorf(
			"Error failed to PayMultiple [%v] [%v]", resp.StatusCode, string(data),
		)
	}

	var response api.MultiplePaymentResponse
	if IsDebug {
		log.Printf("DEBUG RESPONSE PAY MULTIPLE %v %v", resp.StatusCode, string(data))
		return response, nil
	}

	return response, json.Unmarshal(data, &response)
}

// GetBeneficiaries
func (b BankingClient) GetBeneficiaries() (api.Beneficiaries, error) {
	if IsDebug {
		APIurl = "https://openapisandbox.investec.com"
	}
	var url string = fmt.Sprintf("%s/za/pb/v1/accounts/beneficiaries", APIurl)

	resp, err := b.requestAPI(url, http.MethodGet, bearer, nil)
	if err != nil {
		return api.Beneficiaries{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return api.Beneficiaries{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return api.Beneficiaries{}, fmt.Errorf(
			"Error failed to GetBeneficiaries [%v] [%v]", resp.StatusCode, string(data),
		)
	}

	var beneficiaries api.Beneficiaries
	return beneficiaries, json.Unmarshal(data, &beneficiaries)
}
