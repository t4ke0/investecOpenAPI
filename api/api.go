package api

// OAuth2Response
type OAuth2Response struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// Accounts
type Accounts struct {
	Data struct {
		AccountsArr []struct {
			AccountID     string `json:"accountId"`
			AccountNumber string `json:"accountNumber"`
			AccountName   string `json:"accountName"`
			ReferenceName string `json:"referenceName"`
			ProductName   string `json:"productName"`
		} `json:"accounts"`
	} `json:"data"`
}

// Balance
type Balance struct {
	Data struct {
		AccountID        string  `json:"accountId"`
		CurrentBalance   float64 `json:"currentBalance"`
		AvailableBalance float64 `json:"availableBalance"`
		Currency         string  `json:"currency"`
	} `json:"data"`
}

// Transactions
type Transactions struct {
	Data struct {
		AccountTransactions []struct {
			AccountID       string  `json:"accountId"`
			Type            string  `json:"type"`
			TransactionType string  `json:"transactionType"`
			Status          string  `json:"status"`
			Description     string  `json:"description"`
			CardNumber      string  `json:"cardNumber"`
			PostedOrder     int     `json:"postedOrder"`
			PostingDate     string  `json:"postingDate"`
			ValueDate       string  `json:"valueDate"`
			ActionDate      string  `json:"actionDate"`
			TransactionDate string  `json:"transactionDate"`
			Amount          float64 `json:"amount"`
			RunningBalance  float64 `json:"runningBalance"`
		} `json:"transactions"`
	} `json:"data"`
}
