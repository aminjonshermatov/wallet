package wallet

import (
	"github.com/aminjonshermatov/wallet/pkg/types"
	"reflect"
	"testing"
)

func TestService_FindAccountByID_success(t *testing.T) {
	svc := &Service{}

	account, _ := svc.RegisterAccount("+992000000001")

	acc, e := svc.FindAccountByID(account.ID)

	if e != nil {
		t.Error(e)
	}

	if !reflect.DeepEqual(account, acc) {
		t.Error("Accounts doesn't match")
	}
}

func TestService_FindAccountByID_notFound(t *testing.T) {
	svc := &Service{}

	_, e := svc.FindAccountByID(123)

	if e != ErrAccountNotFound {
		t.Error(e)
	}
}

func TestService_FindPaymentByID_success(t *testing.T) {
	svc := &Service{}

	account, errorReg := svc.RegisterAccount("+992000000001")

	if errorReg != nil {
		t.Error("error on register account")
	}

	payment, er := svc.Pay(account.ID, 1000, "auto")

	if er != nil {
		t.Error("error on pay")
	}

	_, err := svc.FindPaymentByID(payment.ID)

	if err != nil {
		t.Error("payment not found")
	}
}

func TestService_FindPaymentByID_notFound(t *testing.T) {
	svc := &Service{}

	_, err := svc.FindPaymentByID("aaa")

	if err != ErrPaymentNotFound {
		t.Error("payment exist")
	}
}

func TestService_Reject_success(t *testing.T) {
	svc := &Service{}

	account, errAccount := svc.RegisterAccount("992000000001")

	if errAccount == ErrPhoneRegistered {
		t.Error(ErrPhoneRegistered)
	}

	payment, errPay := svc.Pay(account.ID, 1000, "auto")

	if errPay == ErrAccountNotFound {
		t.Error(ErrAccountNotFound)
	}

	if errPay == ErrNotEnoughBalance {
		t.Error(ErrNotEnoughBalance)
	}

	errReject := svc.Reject(payment.ID)

	if errReject == ErrPaymentNotFound {
		t.Error(ErrPaymentNotFound)
	}

	if errReject == ErrAccountNotFound {
		t.Error(ErrAccountNotFound)
	}

	if payment.Status != types.PaymentStatusFail {
		t.Error("payment status doesn't failed")
	}
}

func TestService_Reject_notRejectPaymentNotFound(t *testing.T) {
	svc := &Service{}

	account, errAccount := svc.RegisterAccount("992000000001")

	if errAccount == ErrPhoneRegistered {
		t.Error(ErrPhoneRegistered)
	}

	payment, errPay := svc.Pay(account.ID, 1000, "auto")

	if errPay == ErrAccountNotFound {
		t.Error(ErrAccountNotFound)
	}

	if errPay == ErrNotEnoughBalance {
		t.Error(ErrNotEnoughBalance)
	}

	errReject := svc.Reject("aaaa")

	if errReject != ErrPaymentNotFound {
		t.Error(ErrPaymentNotFound)
	}

	if errReject == ErrAccountNotFound {
		t.Error(ErrAccountNotFound)
	}

	if payment.Status != types.PaymentStatusInProgress {
		t.Error("payment status doesn't failed")
	}
}