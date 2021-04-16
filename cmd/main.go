package main

import (
	"github.com/aminjonshermatov/wallet/pkg/wallet"
	"log"
)

func main() {
	svc := &wallet.Service{}

	/*account1, err := svc.RegisterAccount("+992000000001")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = svc.Deposit(account1.ID, 1_000)
	if err != nil {
		fmt.Println(err)
		return
	}

	account2, err := svc.RegisterAccount("+992000000002")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = svc.Deposit(account2.ID, 2_000)
	if err != nil {
		fmt.Println(err)
		return
	}

	account3, err := svc.RegisterAccount("+992000000003")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = svc.Deposit(account3.ID, 3_000)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = svc.ExportToFile("accounts.txt")
	if err != nil {
		log.Print(err)
	}*/
	err := svc.ImportFromFile("accounts.txt")
	if err != nil {
		log.Print(err)
	}
}
