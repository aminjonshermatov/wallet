package main

import (
	"github.com/aminjonshermatov/wallet/pkg/types"
	"github.com/aminjonshermatov/wallet/pkg/wallet"
	"log"
)

func main() {
	svc := &wallet.Service{}

	for i := 0; i < 5; i++ {
		account, err := svc.RegisterAccount("+99200000000" + types.Phone(rune(i)))
		if err != nil {
			log.Print(err)
			break
		}

		err = svc.Deposit(account.ID, types.Money(1_000 * (i + 1)))
		if err != nil {
			log.Print(err)
			break
		}

		payment, err := svc.Pay(account.ID, types.Money(500 * (i + 1)), "foo")
		if err != nil {
			log.Print(err)
			break
		}

		_, err = svc.FavoritePayment(payment.ID, "FOO")
		if err != nil {
			log.Print(err)
			break
		}
	}

	err := svc.Export("data/c258dee9-e7be-4a19-909e-a7d883c166a7")
	if err != nil {
		log.Print(err)
		return
	}
	err = svc.Import("data/c258dee9-e7be-4a19-909e-a7d883c166a7")
	if err != nil {
		log.Print(err)
		return
	}
	svc.Log("accounts")
	svc.Log("payments")
	svc.Log("favorites")
}
