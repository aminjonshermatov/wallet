package wallet

import (
	"fmt"
	"github.com/aminjonshermatov/wallet/pkg/types"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

type testService struct {
	*Service
}

type testAccount struct {
	phone		types.Phone
	balance		types.Money
	payments	[]struct{
		amount		types.Money
		category	types.PaymentCategory
	}
}

var defaultTestAccount = testAccount{
	phone:		"+992000000001",
	balance: 	10_000_00,
	payments: 	[]struct{
		amount		types.Money
		category	types.PaymentCategory
	}{
		{amount: 1_000_00, category: "auto"},
	},
}

func newTestService() *testService {
	return &testService{Service: &Service{}}
}

func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can't regist account,  error = %v", err)
	}

	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't deposity account, error = %v", err)
	}

	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't make payment, error = %v", err)
		}
	}

	return account, payments, nil
}

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
	s := newTestService()

	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	got, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindPaymentByID(): error = %v", err)
		return
	}

	if !reflect.DeepEqual(payment, got) {
		t.Errorf("FindPaymentByID(): wrong paymen returned = %v", err)
		return
	}
}

func TestService_FindPaymentByID_fail(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.FindPaymentByID(uuid.New().String())
	if err == nil {
		t.Error("FindPaymentByID(): must return error, returned nil")
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentByID(): must retunrn ErrPaymentNotFound, returned = %v", err)
		return
	}
}

func TestService_Reject_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
		return
	}

	savedPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Reject(): can't find payment by id, error = %v", err)
		return
	}
	if savedPayment.Status != types.PaymentStatusFail {
		t.Errorf("Reject(): status can't changed, error = %v", savedPayment)
		return
	}

	savedAccount, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("Reject(): can't find account by id, error = %v", err)
		return
	}
	if savedAccount.Balance != defaultTestAccount.balance {
		t.Errorf("Reject(), balance didn't cahnged, account = %v", savedAccount)
	}
}

func TestService_Reject_fail(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	err = s.Reject(uuid.New().String())
	if err == nil {
		t.Error("Reject(): must be error, returned nil")
		return
	}
}

func TestService_Repeat_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	newPayment, nil := s.Repeat(payment.ID)
	if err != nil {
		t.Error(err)
		return
	}

	if payment.ID == newPayment.ID {
		t.Error("repeated payment id not different")
		return
	}

	if payment.AccountID != newPayment.AccountID ||
		payment.Status != newPayment.Status ||
		payment.Category != newPayment.Category ||
		payment.Amount != newPayment.Amount {
		t.Error("some field is not equal the original")
	}
}
