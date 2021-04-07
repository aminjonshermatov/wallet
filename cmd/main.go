package main

import (
	"fmt"
	"github.com/aminjonshermatov/wallet/pkg/wallet"
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
		switch err {
		case wallet.ErrAmountMustBePositive:
			fmt.Println("Amount must be greater then 0")
		case wallet.ErrAccountNotFound:
			fmt.Println("Account not found")
		}
		return
	}

	fmt.Println(account.Balance) // 10
}
