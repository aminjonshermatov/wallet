package main

import (
	"fmt"
	"github.com/aminjonshermatov/wallet/pkg/wallet"
	"github.com/google/uuid"
)

func main() {
	svc := &wallet.Service{}

	account, err := svc.RegisterAccount("+992000000001")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = svc.Deposit(account.ID, 10)
	if err != nil {
		return
	}

	payment, err := svc.PayFromFavorite(uuid.New().String())
	id := payment.ID
	fmt.Println(id)
}
