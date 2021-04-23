package main

import (
	"github.com/aminjonshermatov/wallet/pkg/wallet"
	"log"
)

func main() {
	svc := &wallet.Service{}

	//for i := 0; i < 5; i++ {
	//	account, err := svc.RegisterAccount("+99200000000" + types.Phone(i))
	//	if err != nil {
	//		log.Print(err)
	//		break
	//	}
	//
	//	err = svc.Deposit(account.ID, types.Money(1_000 * (i + 1)))
	//	if err != nil {
	//		log.Print(err)
	//		break
	//	}
	//
	//	payment, err := svc.Pay(account.ID, types.Money(500 * (i + 1)), "foo")
	//	if err != nil {
	//		log.Print(err)
	//		break
	//	}
	//
	//	_, err = svc.FavoritePayment(payment.ID, "FOO")
	//	if err != nil {
	//		log.Print(err)
	//		break
	//	}
	//}
	//
	//err := svc.Export("data")
	//if err != nil {
	//	log.Print(err)
	//	return
	//}
	svc.Log("favorites")
	err := svc.Import("data")
	if err != nil {
		log.Print(err)
		return
	}
	svc.Log("favorites")
}
