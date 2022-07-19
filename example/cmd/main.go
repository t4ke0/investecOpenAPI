package main

import (
	"log"

	client "github.com/t4ke0/investecOpenAPI"
)

// fill in the following vars
var (
	clientID string = ""
	secret   string = ""
)

func main() {

	clt := client.NewBankingClient(secret, clientID)

	if err := clt.GetAccessToken(); err != nil {
		log.Fatal(err)
	}

	accounts, err := clt.GetAccounts()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("debug accounts", accounts)

	// fill in the accountID var with your account id
	var accountID string = ""
	balance, err := clt.GetAccountBalance(accountID)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("debug", balance)

	// i'm not using here fromDate or toDate params
	transactions, err := clt.GetAccountTransactions(accountID, "", "")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("debug transactions %v", transactions)

}
