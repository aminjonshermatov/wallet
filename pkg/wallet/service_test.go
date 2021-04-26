package wallet

import (
	"fmt"
	"github.com/aminjonshermatov/wallet/pkg/types"
	"github.com/google/uuid"
	"log"
	"reflect"
	"sort"
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

func TestService_FavoritePayment_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]

	_, err = s.FavoritePayment(payment.ID, "osh")
	if err != nil {
		t.Error(err)
	}
}

func TestService_FavoritePayment_fail(t *testing.T) {
	s := newTestService()

	_, err := s.FavoritePayment(uuid.New().String(), "osh")
	if err == nil {
		t.Error("FavoritePayment(): must return error, now nil")
	}
}

func TestService_PayFromFavorite_success(t *testing.T) {
	s := newTestService()

	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error("PayFromFavorite(): can't get payments")
		return
	}

	payment := payments[0]

	favorite, err := s.FavoritePayment(payment.ID, "osh")
	if err != nil {
		t.Error("PayFromFavorite(): can't add payment to favorite")
		return
	}

	_, err = s.PayFromFavorite(favorite.ID)
	if err != nil {
		t.Error("PayFromFavorite(): can't not pay from favorite")
		return
	}
}

func TestService_PayFromFavorite_fail(t *testing.T) {
	s := newTestService()

	_, err := s.PayFromFavorite(uuid.New().String())
	if err == nil {
		t.Error("PayFromFavorite(): must be error, now returned nil")
	}
}

func TestService_ExportToFile_success(t *testing.T) {
	s := newTestService()

	_, err := s.RegisterAccount("+992000000001")
	if err != nil {
		t.Error(err)
	}

	path := "export.txt"
	err = s.ExportToFile(path)
	if err != nil {
		t.Error(err)
	}
}

func TestService_ExportToFile_fail(t *testing.T) {
	s := newTestService()

	_, err := s.RegisterAccount("+992000000001")
	if err != nil {
		t.Error(err)
	}

	path := "data/export.txt"
	err = s.ExportToFile(path)
	if err == nil {
		t.Error(err)
	}
}

func TestService_ImportFromFile_success(t *testing.T) {
	s := newTestService()

	_, err := s.RegisterAccount("+992000000001")
	if err != nil {
		t.Error(err)
	}

	path := "export.txt"
	err = s.ExportToFile(path)
	if err != nil {
		t.Error(err)
	}

	s.accounts = s.accounts[:0]
	err = s.ImportFromFile(path)
	if err != nil {
		t.Error(err)
	}
}

func TestService_ImportFromFile_fail(t *testing.T) {
	s := newTestService()

	_, err := s.RegisterAccount("+992000000001")
	if err != nil {
		t.Error(err)
	}

	path := "export.txt"
	err = s.ExportToFile(path)
	if err != nil {
		t.Error(err)
	}

	s.accounts = s.accounts[:0]
	err = s.ImportFromFile("data/" + path)
	if err == nil {
		t.Error(err)
	}
}

func TestService_Export_regular(t *testing.T) {
	s := newTestService()

	account, err := s.RegisterAccount("+992000000001")
	if err != nil {
		t.Error(err)
	}

	err = s.Deposit(account.ID, types.Money(100_000))
	if err != nil {
		t.Error(err)
	}

	payment, err := s.Pay(account.ID, types.Money(1_000), "foo")
	if err != nil {
		t.Error(err)
	}

	_, err = s.FavoritePayment(payment.ID, "FOO")
	if err != nil {
		t.Error(err)
	}

	path := "data"

	err = s.Export(path)
	if err != nil {
		t.Error(err)
	}
}

func TestService_Export_emptySlices(t *testing.T) {
	s := newTestService()

	path := "data"

	err := s.Export(path)
	if err != nil {
		t.Error(err)
	}
}

func TestService_Import_regular(t *testing.T) {
	s := newTestService()

	account, err := s.RegisterAccount("+992000000001")
	if err != nil {
		t.Error(err)
	}

	err = s.Deposit(account.ID, types.Money(100_000))
	if err != nil {
		t.Error(err)
	}

	payment, err := s.Pay(account.ID, types.Money(1_000), "foo")
	if err != nil {
		t.Error(err)
	}

	_, err = s.FavoritePayment(payment.ID, "FOO")
	if err != nil {
		t.Error(err)
	}

	path := "data"

	err = s.Export(path)
	if err != nil {
		t.Error(err)
	}

	err = s.Import(path)
	if err != nil {
		t.Error(err)
	}
}

func TestService_Import_clearAfterExport(t *testing.T) {
	s := newTestService()

	account, err := s.RegisterAccount("+992000000001")
	if err != nil {
		t.Error(err)
	}

	err = s.Deposit(account.ID, types.Money(100_000))
	if err != nil {
		t.Error(err)
	}

	payment, err := s.Pay(account.ID, types.Money(1_000), "foo")
	if err != nil {
		t.Error(err)
	}

	_, err = s.FavoritePayment(payment.ID, "FOO")
	if err != nil {
		t.Error(err)
	}

	path := "data"

	err = s.Export(path)
	if err != nil {
		t.Error(err)
	}

	s.accounts = s.accounts[:0]
	s.payments = s.payments[:0]
	s.favorites = s.favorites[:0]

	err = s.Import(path)
	if err != nil {
		t.Error(err)
	}
}

func TestService_Import_emptySlices(t *testing.T) {
	s := newTestService()

	path := "data"

	err := s.Export(path)
	if err != nil {
		t.Error(err)
	}

	err = s.Import(path)
	if err != nil {
		t.Error(err)
	}
}

func TestService_ExportAccountHistory_success(t *testing.T)  {
	s := newTestService()

	account, err := s.RegisterAccount("+992000000001")
	if err != nil {
		t.Error(err)
	}

	err = s.Deposit(account.ID, types.Money(100_000))
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 7; i++ {
		_, err := s.Pay(account.ID, types.Money(1 + i), "foo")
		if err != nil {
			t.Error(err)
			break
		}
	}

	payments, err := s.ExportAccountHistory(account.ID)
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < len(s.payments); i++ {
		if *s.payments[i] != payments[i] {
			t.Errorf("payments is not matches, got %v, want %v", payments, s.payments)
		}
	}
}

func TestService_ExportAccountHistory_fail(t *testing.T)  {
	s := newTestService()

	account, err := s.RegisterAccount("+992000000001")
	if err != nil {
		t.Error(err)
	}

	err = s.Deposit(account.ID, types.Money(100_000))
	if err != nil {
		t.Error(err)
	}

	payments, err := s.ExportAccountHistory(321)
	if err == nil {
		t.Error(err)
	}

	for i := 0; i < len(s.payments); i++ {
		if *s.payments[i] == payments[i] {
			t.Errorf("payments is not matches, got %v, want %v", payments, s.payments)
		}
	}
}

func TestService_HistoryToFiles_noData(t *testing.T) {
	s := newTestService()

	payments := make([]types.Payment, 0)

	for _, payment := range s.payments {
		payments = append(payments, *payment)
	}

	err := s.HistoryToFiles(payments, "data", 3)
	if err != nil {
		t.Error(err)
	}
}

func TestService_HistoryToFiles_negativeRecords(t *testing.T) {
	s := newTestService()

	payments := make([]types.Payment, 0)

	for _, payment := range s.payments {
		payments = append(payments, *payment)
	}

	err := s.HistoryToFiles(payments, "data", -3)
	if err == nil {
		t.Error(err)
	}
}

func TestService_HistoryToFiles_OneFile(t *testing.T) {
	s := newTestService()

	account, err := s.RegisterAccount("+992000000001")
	if err != nil {
		log.Print(err)
		return
	}

	err = s.Deposit(account.ID, types.Money(100_000))
	if err != nil {
		log.Print(err)
		return
	}

	for i := 0; i < 7; i++ {
		_, err := s.Pay(account.ID, types.Money(1 + i), "foo")
		if err != nil {
			log.Print(err)
			break
		}
	}

	payments := make([]types.Payment, 0)

	for _, payment := range s.payments {
		payments = append(payments, *payment)
	}

	err = s.HistoryToFiles(payments, "data", 8)
	if err != nil {
		t.Error(err)
	}
}

func TestService_HistoryToFiles_multiFile(t *testing.T) {
	s := newTestService()

	account, err := s.RegisterAccount("+992000000001")
	if err != nil {
		log.Print(err)
		return
	}

	err = s.Deposit(account.ID, types.Money(100_000))
	if err != nil {
		log.Print(err)
		return
	}

	for i := 0; i < 7; i++ {
		_, err := s.Pay(account.ID, types.Money(1 + i), "foo")
		if err != nil {
			t.Error(err)
			break
		}
	}

	payments := make([]types.Payment, 0)

	for _, payment := range s.payments {
		payments = append(payments, *payment)
	}

	err = s.HistoryToFiles(payments, "data", 3)
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkService_SumPayments(b *testing.B) {
	s := newTestService()
	account, err := s.RegisterAccount("+992000000001")
	if err != nil {
		b.Fatal(err)
		return
	}

	err = s.Deposit(account.ID, types.Money(100_000))
	if err != nil {
		b.Fatal(err)
		return
	}

	for i := 0; i < 7; i++ {
		_, err := s.Pay(account.ID, types.Money(1 + i), "foo")
		if err != nil {
			b.Fatal(err)
		}
	}

	want := types.Money(28)

	for i := 0; i < b.N; i++ {
		result := s.SumPayments(i)
		if result != want {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
	}
}

func BenchmarkService_FilterPayments_success(b *testing.B) {
	s := newTestService()

	account1, err := s.RegisterAccount("+992000000000")
	if err != nil {
		b.Fatal(err)
	}
	err = s.Deposit(account1.ID, 100_000)
	if err != nil {
		b.Fatal(err)
	}
	wantPayments := make([]types.Payment, 0)
	payment1Acc1, err := s.Pay(account1.ID, 1_000, "goo")
	if err != nil {
		b.Fatal(err)
	}
	payment2Acc1, err := s.Pay(account1.ID, 2_000, "goo")
	if err != nil {
		b.Fatal(err)
	}

	wantPayments = append(wantPayments, *payment1Acc1)
	wantPayments = append(wantPayments, *payment2Acc1)

	account2, err := s.RegisterAccount("+992000000001")
	if err != nil {
		b.Fatal(err)
	}
	err = s.Deposit(account2.ID, 100_000)
	if err != nil {
		b.Fatal(err)
	}

	_, err = s.Pay(account2.ID, 1_000, "goo")
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		payments, err := s.FilterPayments(account1.ID, i)
		if err != nil {
			b.Fatal(err)
		}
		sort.Slice(payments, func(i, j int) bool {
			return payments[i].ID < payments[j].ID
		})

		sort.Slice(wantPayments, func(i, j int) bool {
			return wantPayments[i].ID < wantPayments[j].ID
		})
		if !reflect.DeepEqual(wantPayments, payments) {
			b.Fatalf("payments is not equal, got %v, want %v", payments, wantPayments)
		}
	}

}

func BenchmarkService_FilterPayments_invalidID(b *testing.B) {
	s := newTestService()

	for i := 0; i < b.N; i++ {
		_, err := s.FilterPayments(213456789, i)
		if err == nil {
			b.Fatal("result must be error, because invalid accountID")
		}
	}

}