package investecOpenAPI

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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

func (b BankingClient) requestAPI(url, method string, mode authMode) (*http.Response, error) {
	var (
		req *http.Request
		err error
	)
	if mode == basic {
		req, err = http.NewRequest(method, url, strings.NewReader("grant_type=client_credentials&scope=accounts"))
	} else {
		req, err = http.NewRequest(method, url, nil)
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
	resp, err := b.requestAPI(url, http.MethodPost, basic)
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
	resp, err := b.requestAPI(url, http.MethodGet, bearer)
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
	resp, err := b.requestAPI(url, http.MethodGet, bearer)
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

	resp, err := b.requestAPI(url, http.MethodGet, bearer)
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
