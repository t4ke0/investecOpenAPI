package main

import (
	"log"

	client "github.com/t4ke0/investecOpenAPI"
)

// sandbox data
var (
	clientID string = "yAxzQRFX97vOcyQAwluEU6H6ePxMA5eY"
	secret   string = "4dY0PjEYqoBrZ99r"
	key      string = "eUF4elFSRlg5N3ZPY3lRQXdsdUVVNkg2ZVB4TUE1ZVk6YVc1MlpYTjBaV010ZW1FdGNHSXRZV05qYjNWdWRITXRjMkZ1WkdKdmVBPT0="
)

func main() {

	// set this to false if you want to test production API
	client.IsDebug = true

	clt := client.NewBankingClient(key, secret, clientID)

	if err := clt.GetAccessToken(); err != nil {
		log.Fatal(err)
	}

	accounts, err := clt.GetAccounts()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("debug accounts %v", accounts)

	// fill in the accountID var with your account id
	var accountID string = ""
	balance, err := clt.GetAccountBalance(accountID)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("debug %v", balance)

	// i'm not using here fromDate or toDate params
	transactions, err := clt.GetAccountTransactions(accountID, "", "")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("debug transactions %v", transactions)

}
