package wallet

import (
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

	_, er := svc.Pay(account.ID, 1000, "auto")

	if er == ErrAmountMustBePositive {
		t.Error(ErrAmountMustBePositive)
	}

	if er == nil {
		t.Error("error on pay")
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

	_, errPay := svc.Pay(account.ID, 1000, "auto")

	if errPay != ErrNotEnoughBalance {
		t.Error(ErrNotEnoughBalance)
	}
}

func TestService_Reject_notRejectPaymentNotFound(t *testing.T) {
	svc := &Service{}

	account, errAccount := svc.RegisterAccount("992000000001")

	if errAccount == ErrPhoneRegistered {
		t.Error(ErrPhoneRegistered)
	}

	_, errPay := svc.Pay(account.ID, 1000, "auto")

	if errPay != ErrNotEnoughBalance {
		t.Error(ErrNotEnoughBalance)
	}
}