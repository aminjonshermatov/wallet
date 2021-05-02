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

	err = svc.Deposit(account.ID, types.Money(9_000_000))
	if err != nil {
		log.Print(err)
		return
	}

	for i := 0; i < 9_000_000; i++ {
		_, err := svc.Pay(account.ID, types.Money(1), "foo")
		if err != nil {
			log.Print(err)
			break
		}
	}

	//svc.Log("payments")

	ch := svc.SumPaymentsWithProgress()
	result := types.Money(0)
	for val := range ch {
		result += val.Result
		log.Printf("part: %d, result: %v", val.Part, val.Result)
	}

	log.Printf("done, sum: %v", result)

	//payments, err := svc.ExportAccountHistory(account.ID)
	//if err != nil {
	//	log.Print(err)
	//	return
	//}
	//
	//log.Print(len(payments))
	//err = svc.HistoryToFiles(payments, "data", 3)
}
