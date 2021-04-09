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

	phone := types.Phone("+992000000001")
	account, err := svc.RegisterAccount(phone)
	if err != nil {
		t.Errorf("FindPaymentByID(): can't register account, error = %v", err)
		return
	}

	err = svc.Deposit(account.ID, 10_000_00)
	if err != nil {
		t.Errorf("FindPaymentByID(): can't deposit account error = %v", err)
		return
	}

	payment, err := svc.Pay(account.ID, 1000_00, "auto")
	if err != nil {
		t.Errorf("FindPaymentByID(): can't create payment, error = %v", err)
		return
	}

	got, err := svc.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindPaymentByID(): error = %v", err)
		return
	}

	if !reflect.DeepEqual(payment, got) {
		t.Errorf("FindPaymentByID(): wrong payment returned, %v", err)
		return
	}
}

func TestService_FindPaymentByID_notFound(t *testing.T) {
	svc := &Service{}

	phone := types.Phone("+992000000001")
	account, err := svc.RegisterAccount(phone)
	if err != nil {
		t.Errorf("FindPaymentByID(): can't register account, error = %v", err)
		return
	}

	err = svc.Deposit(account.ID, 10_000_00)
	if err != nil {
		t.Errorf("FindPaymentByID(): can't deposit account error = %v", err)
		return
	}

	_, err = svc.Pay(account.ID, 1000_00, "auto")
	if err != nil {
		t.Errorf("FindPaymentByID(): can't create payment, error = %v", err)
		return
	}

	_, err = svc.FindPaymentByID("payment.ID")
	if err == nil {
		t.Errorf("FindPaymentByID(): must be error")
		return
	}
}

func TestService_Reject_success(t *testing.T) {
	svc := &Service{}

	phone := types.Phone("+992000000001")
	account, err := svc.RegisterAccount(phone)
	if err != nil {
		t.Errorf("Reject(): can't register account, error = %v", err)
		return
	}

	err = svc.Deposit(account.ID, 10_000_00)
	if err != nil {
		t.Errorf("Reject(): can't deposit account error = %v", err)
		return
	}

	payment, err := svc.Pay(account.ID, 1000_00, "auto")
	if err != nil {
		t.Errorf("Reject(): can't create payment, error = %v", err)
		return
	}

	err = svc.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
		return
	}
}

func TestService_Reject_notRejectPaymentNotFound(t *testing.T) {
	svc := &Service{}

	phone := types.Phone("+992000000001")
	account, err := svc.RegisterAccount(phone)
	if err != nil {
		t.Errorf("Reject(): can't register account, error = %v", err)
		return
	}

	err = svc.Deposit(account.ID, 10_000_00)
	if err != nil {
		t.Errorf("Reject(): can't deposit account error = %v", err)
		return
	}

	_, err = svc.Pay(account.ID, 1000_00, "auto")
	if err != nil {
		t.Errorf("Reject(): can't create payment, error = %v", err)
		return
	}

	err = svc.Reject("payment.ID")
	if err == nil {
		t.Errorf("Reject(): must be err")
		return
	}
}