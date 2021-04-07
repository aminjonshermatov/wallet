package wallet

import (
	"reflect"
	"testing"
)

func TestService_FindAccountByID_success(t *testing.T) {
	svc := Service{}

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
	svc := Service{}

	_, e := svc.FindAccountByID(123)

	if e != ErrAccountNotFound {
		t.Error(e)
	}
}