# investecOpenAPI


## installation

make sure you are in your root of your project folder
and you have a go.mod file there if not you just need to
create it with

```bash
go mod init <project name>
```

```bash
go get -v github.com/t4ke0/investecOpenAPI
```


## Usage

before you start you might need to have the `clientID` and `secret` from the API provider.
then you are good to proceed.

```golang
package main

import (
    ***

	client "github.com/t4ke0/investecOpenAPI"

    ***
)


func main() {
    // testing using sandbox data
	var clientID string = "yAxzQRFX97vOcyQAwluEU6H6ePxMA5eY"
	var secret   string = "4dY0PjEYqoBrZ99r"
	var key      string = "eUF4elFSRlg5N3ZPY3lRQXdsdUVVNkg2ZVB4TUE1ZVk6YVc1MlpYTjBaV010ZW1FdGNHSXRZV05qYjNWdWRITXRjMkZ1WkdKdmVBPT0="
    //

    client.IsDebug = true
    clt := client.NewBankingClient(key, secret, clientID)

    // OAuth to the API
    if err := clt.GetAccessToken(); err != nil {
        log.Fatal(err)
    }

    // Get Accounts.
    accounts, err := clt.GetAccounts()
    if err != nil {
        log.Fatal(err)
    }

    // see the specification of the account type in `api/api.go`
    fmt.Println("Accounts", accounts) 

    // Get Account Balance
    var accountID string = accounts.Data.AccountArr[0].AccountID
    balance, err := clt.GetAccountBalance(accountID)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("balance", balance)

    // Get account transactions
    // you can use fromDate and toData params to filter the transactions.
    // see investecOpenAPI.go file the dates need to have ISO 8601 format.
    transactions, err := clt.GetAccountTransactions(accountID, "", "")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("transactions", transactions)
}
```
## Testing

to Run an example please build `example/cmd/main.go`

```bash
go build example/cmd/main.go
```

```golang
package main

import (
	"encoding/json"
	"fmt"
	"log"

	client "github.com/t4ke0/investecOpenAPI"
	"github.com/t4ke0/investecOpenAPI/api"
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
	var accountID string = accounts.Data.AccountsArr[0].AccountID
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

	beneficiaries, err := clt.GetBeneficiaries()
	if err != nil {
		log.Fatal(err)
	}
	data, _ := json.MarshalIndent(beneficiaries, "", " ")

	fmt.Println(string(data))

	transfermultipleResponse, err := clt.TransferMultiple(accountID, []api.TransferTo{
		api.TransferTo{
			BeneficiaryAccountId: "MTAxOTAwMjQ2MTI2NjA=",
			Amount:               "10.00",
			MyReference:          "Test",
			TheirReference:       "STD Ben Ref",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	data, _ = json.MarshalIndent(transfermultipleResponse, "", " ")
	fmt.Println(string(data))

	payMultipleResponse, err := clt.PayMultiple(accountID, []api.PaymentMultiple{
		api.PaymentMultiple{
			BeneficiaryAccountId: "MTAxOTAwMjQ2MTI2NjA=",
			Amount:               "10.00",
			MyReference:          "Test",
			TheirReference:       "STD Ben Ref",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	data, _ = json.MarshalIndent(payMultipleResponse, "", " ")
	fmt.Println(string(data))

}
```
