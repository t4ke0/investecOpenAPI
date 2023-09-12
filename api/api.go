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

type TransferTo struct {
	BeneficiaryAccountId string `json:"beneficiaryAccountId"`
	Amount               string `json:"amount"`
	MyReference          string `json:"myReference"`
	TheirReference       string `json:"theirReference"`
}

type MultipleTransfersResponse struct {
	Data struct {
		TransferResponses []struct {
			PaymentReferenceNumber string `json:"PaymentReferenceNumber"`
			PaymentDate            string `json:"PaymentDate"`
			Status                 string `json:"Status"`
			BeneficiaryName        string `json:"BeneficiaryName"`
			BeneficiaryAccountId   string `json:"BeneficiaryAccountId"`
		} `json:"TransferResponses"`
		ErrorMessage string `json:"ErrorMessage"`
	} `json:"data"`
}

type PaymentMultiple struct {
	BeneficiaryId  string `json:"beneficiaryId"`
	Amount         string `json:"amount"`
	MyReference    string `json:"myReference"`
	TheirReference string `json:"theirReference"`
}

type MultiplePaymentResponse struct {
	MultipleTransfersResponse
}

type Beneficiaries struct {
	Data []struct {
		BeneficiaryId          string `json:"BeneficiaryId"`
		AccountNumber          string `json:"AccountNumber"`
		Code                   string `json:"Code"`
		Bank                   string `json:"Bank"`
		BeneficiaryName        string `json:"BeneficiaryName"`
		LastPaymentAmount      string `json:"LastPaymentAmount"`
		LastPaymentDate        string `json:"LastPaymentDate"`
		CellNo                 string `json:"CellNo"`
		EmailAddress           string `json:"EmailAddress"`
		Name                   string `json:"Name"`
		ReferenceAccountNumber string `json:"ReferenceAccountNumber"`
		ReferenceName          string `json:"ReferenceName"`
		CategoryId             string `json:"CategoryId"`
		ProfileId              string `json:profileId`
	} `json:"data"`
}
