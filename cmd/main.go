package main

import (
	"github.com/aminjonshermatov/wallet/pkg/types"
	"github.com/aminjonshermatov/wallet/pkg/wallet"
	"log"
)

func main() {
	svc := &wallet.Service{}

	account, err := svc.RegisterAccount("+992000000001")
	if err != nil {
		log.Print(err)
		return
	}

	err = svc.Deposit(account.ID, types.Money(100_000))
	if err != nil {
		log.Print(err)
		return
	}

	for i := 0; i < 9; i++ {
		_, err := svc.Pay(account.ID, types.Money(500 + i), "foo")
		if err != nil {
			log.Print(err)
			break
		}
	}

	payments, err := svc.ExportAccountHistory(account.ID)
	if err != nil {
		log.Print(err)
		return
	}

	log.Print(len(payments))
	err = svc.HistoryToFiles(payments, "data", 3)
}
