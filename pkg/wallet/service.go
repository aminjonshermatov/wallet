package wallet

import (
	"bufio"
	"errors"
	"github.com/aminjonshermatov/wallet/pkg/types"
	"github.com/google/uuid"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater then 0")
var ErrAccountNotFound = errors.New("account not found")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrNotEnoughBalance = errors.New("account balance least then amount")
var ErrFavoriteNotFound = errors.New("favorite payment not found")

type Service struct {
	nextAccountID	int64
	accounts		[]*types.Account
	payments		[]*types.Payment
	favorites		[]*types.Favorite
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}

	s.nextAccountID++
	account := &types.Account{
		ID: 		s.nextAccountID,
		Phone: 		phone,
		Balance: 	0,
	}

	s.accounts = append(s.accounts, account)
	return account, nil
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}

	account, err := s.FindAccountByID(accountID)
	if err != nil {
		return err
	}

	account.Balance += amount
	return nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	account, err := s.FindAccountByID(accountID)
	if err != nil {
		return nil, err
	}

	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount

	paymentID := uuid.New().String()

	payment := &types.Payment{
		ID:			paymentID,
		AccountID: 	accountID,
		Amount: 	amount,
		Category: 	category,
		Status: 	types.PaymentStatusInProgress,
	}

	s.payments = append(s.payments, payment)

	return payment, nil
}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	var account *types.Account

	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	return account, nil
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}

	return nil, ErrPaymentNotFound
}

func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return err
	}

	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return err
	}

	account.Balance += payment.Amount
	payment.Amount = 0
	payment.Status = types.PaymentStatusFail
	return nil
}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	newPayment, err := s.Pay(payment.AccountID, payment.Amount, payment.Category)
	if err != nil {
		return nil, err
	}

	return newPayment, nil
}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	favorite := &types.Favorite{
		ID:			uuid.New().String(),
		AccountID: 	payment.AccountID,
		Name: 		name,
		Amount: 	payment.Amount,
		Category: 	payment.Category,
	}

	s.favorites = append(s.favorites, favorite)
	return favorite, nil
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	var targetFavorite *types.Favorite

	for _, favorite := range s.favorites {
		if favorite.ID == favoriteID {
			targetFavorite = favorite
			break
		}
	}

	if targetFavorite == nil {
		return nil, ErrFavoriteNotFound
	}

	payment, err := s.Pay(targetFavorite.AccountID, targetFavorite.Amount, targetFavorite.Category)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *Service) ExportToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		cerr := file.Close()
		if cerr != nil {
			if err == nil {
				err = cerr
			}
		}
	} ()

	content := make([]byte, 0)
	for _, account := range s.accounts {
		content = append(content, []byte(strconv.FormatInt(account.ID, 10) + ";")...)
		content = append(content, []byte(account.Phone + ";")...)
		content = append(content, []byte(strconv.FormatInt(int64(account.Balance), 10) + "|")...)
	}

	_, err = file.Write(content)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ImportFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		cerr :=file.Close()
		if cerr != nil {
			if err == nil {
				err = cerr
			}
		}
	}()

	content := make([]byte, 0)
	buf := make([]byte, 4)
	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			content = append(content, buf[:read]...)
			break
		}

		content = append(content, buf[:read]...)
	}

	for _, row := range strings.Split(string(content), "|") {
		col := strings.Split(row, ";")
		if len(col) == 3 {
			_, err = s.RegisterAccount(types.Phone(col[1]))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func create(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}

func (s *Service) Export(dir string) error {
	err := ExportAccounts(s, dir)
	if err != nil {
		return err
	}

	err = ExportPayments(s, dir)
	if err != nil {
		return err
	}

	err = ExportFavorites(s, dir)
	if err != nil {
		return err
	}

	return nil
}

func ExportAccounts(s *Service, dir string) (err error) {
	if len(s.accounts) == 0 {
		return nil
	}

	file, err := create(dir + "/" + "accounts.dump")
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			if err == nil {
				err = cerr
			}
		}
	}()
	data := make([]byte, 0)

	for _, account := range s.accounts {

		data = append(data, []byte(strconv.FormatInt(account.ID, 10) + ";")...)
		data = append(data, []byte(account.Phone + ";")...)
		data = append(data, []byte(strconv.FormatInt(int64(account.Balance), 10) + "\n")...)
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}
func ExportPayments(s *Service, dir string) (err error) {
	if len(s.payments) == 0 {
		return nil
	}

	file, err := create(dir + "/" + "payments.dump")
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			if err == nil {
				err = cerr
			}
		}
	}()
	data := make([]byte, 0)

	for _, payment := range s.payments {

		data = append(data, []byte(payment.ID + ";")...)
		data = append(data, []byte(strconv.FormatInt(payment.AccountID, 10) + ";")...)
		data = append(data, []byte(strconv.FormatInt(int64(payment.Amount), 10) + ";")...)
		data = append(data, []byte(payment.Category + ";")...)
		data = append(data, []byte(payment.Status + "\n")...)
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}
func ExportFavorites(s *Service, dir string) (err error) {
	if len(s.favorites) == 0 {
		return nil
	}

	file, err := create(dir + "/" + "favorites.dump")
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			if err == nil {
				err = cerr
			}
		}
	}()
	data := make([]byte, 0)

	for _, favorite := range s.favorites {

		data = append(data, []byte(favorite.ID + ";")...)
		data = append(data, []byte(strconv.FormatInt(favorite.AccountID, 10) + ";")...)
		data = append(data, []byte(favorite.Name + ";")...)
		data = append(data, []byte(strconv.FormatInt(int64(favorite.Amount), 10) + ";")...)
		data = append(data, []byte(favorite.Category + "\n")...)
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Import(dir string) error {
	err := ImportAccounts(s, dir)
	if err != nil {
		return err
	}

	err = ImportPayments(s, dir)
	if err != nil {
		return err
	}

	err = ImportFavorites(s, dir)
	if err != nil {
		return err
	}
	return nil
}

func ImportAccounts(s *Service, dir string) (err error) {
	_, err = os.Stat(dir + "/" + "accounts.dump")
	if !os.IsNotExist(err) {
		src, err := os.Open(dir + "/" + "accounts.dump")
		if err != nil {
			return err
		}
		defer func() {
			if cerr := src.Close(); cerr != nil {
				if err == nil {
					err = cerr
				}
			}
		}()

		reader := bufio.NewReader(src)

		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			line = strings.Replace(line, "\n", "", 1)
			col := strings.Split(line, ";")
			newAccount := &types.Account{
				Phone: types.Phone(col[1]),
			}
			num, err := strconv.Atoi(col[0])
			if  err != nil {
				return err
			}
			newAccount.ID = int64(num)

			balance, err := strconv.Atoi(col[2])
			if err != nil {
				return err
			}
			newAccount.Balance = types.Money(balance)

			isFind := false
			for _, account := range s.accounts {
				if account.ID == newAccount.ID {
					isFind = true
					break
				}
			}

			if !isFind {
				s.accounts = append(s.accounts, newAccount)
			}
		}
		return nil
	}
	return nil
}
func ImportPayments(s *Service, dir string) (err error) {
	_, err = os.Stat(dir + "/" + "payments.dump")
	if !os.IsNotExist(err) {
		src, err := os.Open(dir + "/" + "payments.dump")
		if err != nil {
			return err
		}
		defer func() {
			if cerr := src.Close(); cerr != nil {
				if err == nil {
					err = cerr
				}
			}
		}()

		reader := bufio.NewReader(src)

		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			line = strings.Replace(line, "\n", "", 1)
			col := strings.Split(line, ";")
			newPayment := &types.Payment{
				ID: col[0],
			}
			num, err := strconv.Atoi(col[1])
			if  err != nil {
				return err
			}
			newPayment.AccountID = int64(num)
			amount, err := strconv.Atoi(col[2])
			if  err != nil {
				return err
			}
			newPayment.Amount = types.Money(int64(amount))

			newPayment.Category = types.PaymentCategory(col[3])
			newPayment.Status = types.PaymentStatus(col[4])

			isFind := false
			for _, payment := range s.payments {
				if payment.ID == newPayment.ID {
					isFind = true
					break
				}
			}

			if !isFind {
				s.payments = append(s.payments, newPayment)
			}
		}
		return nil
	}

	return nil
}
func ImportFavorites(s *Service, dir string) (err error) {
	_, err = os.Stat(dir + "/" + "favorites.dump")
	if !os.IsNotExist(err) {
		src, err := os.Open(dir + "/" + "favorites.dump")
		if err != nil {
			return err
		}
		defer func() {
			if cerr := src.Close(); cerr != nil {
				if err == nil {
					err = cerr
				}
			}
		}()

		reader := bufio.NewReader(src)

		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			line = strings.Replace(line, "\n", "", 1)
			col := strings.Split(line, ";")
			newFavorite := &types.Favorite{
				ID: col[0],
			}
			num, err := strconv.Atoi(col[1])
			if  err != nil {
				return err
			}
			newFavorite.AccountID = int64(num)

			newFavorite.Name = col[2]

			amount, err := strconv.Atoi(col[3])
			if  err != nil {
				return err
			}
			newFavorite.Amount = types.Money(int64(amount))

			newFavorite.Category = types.PaymentCategory(col[4])

			isFind := false
			for _, favorite := range s.favorites {
				if favorite.ID == newFavorite.ID {
					isFind = true
					break
				}
			}

			if !isFind {
				s.favorites = append(s.favorites, newFavorite)
			}
		}
		return nil
	}

	return nil
}

func (s *Service) ExportAccountHistory(accountID int64) ([]types.Payment, error) {
	res := make([]types.Payment, 0)

	for _, payment := range s.payments {
		if payment.AccountID == accountID {
			res = append(res, *payment)
		}
	}

	if len(res) == 0 {
		return nil, ErrAccountNotFound
	}
	return res, nil
}

func ExportToFileFrom(dir string, payments []types.Payment, start int, end int, idx string) error {
	file, err := create(dir + "/" + "payments" + idx + ".dump")
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			if err == nil {
				cerr = err
			}
		}
	}()

	data := make([]byte, 0)

	for i := start; i <= end; i++ {
		data = append(data, []byte(payments[i].ID + ";")...)
		data = append(data, []byte(strconv.FormatInt(payments[i].AccountID, 10) + ";")...)
		data = append(data, []byte(strconv.FormatInt(int64(payments[i].Amount), 10) + ";")...)
		data = append(data, []byte(payments[i].Category + ";")...)
		data = append(data, []byte(payments[i].Status + "\n")...)
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) HistoryToFiles(payments []types.Payment, dir string, records int) error {
	if records <= 0 {
		return errors.New("records must be non zero")
	}
	if len(payments) > 0 && len(payments) <= records {
		return ExportToFileFrom(dir, payments, 0, len(payments) - 1, "")
	} else {
		for i := 1; i <= int(math.Ceil(float64(len(payments)) / float64(records))); i++ {
			end := i * records - 1
			if end >= len(payments) {
				end = len(payments) - 1
			}
			err := ExportToFileFrom(dir, payments, records * (i - 1), end, strconv.FormatInt(int64(i), 10))
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func (s *Service) SumPayments(goroutines int) types.Money {
	wg := sync.WaitGroup{}
	if goroutines == 0 {
		goroutines = 1
	}
	wg.Add(goroutines)

	mu := sync.Mutex{}
	sum := types.Money(0)

	lenPayments := len(s.payments)
	for i := 0; i < goroutines; i++ {
		remainder := 0
		if lenPayments % goroutines != 0 {
			remainder = 1
		}

		start := i * lenPayments / goroutines + remainder
		if i == 0 {
			start = 0
		}
		end := (i + 1) * lenPayments / goroutines + remainder
		if end > lenPayments {
			end = lenPayments
		}

		go func(partPayments []*types.Payment) {
			sumPart := types.Money(0)
			defer wg.Done()
			for _, payment := range partPayments {
				sumPart += payment.Amount
			}
			mu.Lock()
			defer mu.Unlock()
			sum += sumPart
		}(s.payments[start:end])
	}

	wg.Wait()
	return sum
}

func (s *Service) FilterPayments(accountID int64, goroutines int) ([]types.Payment, error) {
	_, err := s.FindAccountByID(accountID)
	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	if goroutines == 0 {
		goroutines = 1
	}
	wg.Add(goroutines)

	mu := sync.Mutex{}
	filteredPayments := make([]types.Payment, 0)

	lenPayments := len(s.payments)

	for i := 0; i < goroutines; i++ {
		remainder := 0
		if lenPayments % goroutines != 0 {
			remainder = 1
		}

		start := i * lenPayments / goroutines + remainder
		if i == 0 {
			start = 0
		}
		end := (i + 1) * lenPayments / goroutines + remainder
		if end > lenPayments {
			end = lenPayments
		}

		go func(partPayments []*types.Payment) {
			filtered := make([]types.Payment, 0)
			defer wg.Done()
			for _, payment := range partPayments {
				if payment.AccountID == accountID {
					filtered = append(filtered, *payment)
				}
			}
			mu.Lock()
			defer mu.Unlock()
			filteredPayments = append(filteredPayments, filtered...)
		}(s.payments[start:end])
	}

	wg.Wait()
	return filteredPayments, nil
}

func (s *Service) FilterPaymentsByFn(filter func(payment types.Payment) bool, goroutines int) ([]types.Payment, error) {
	wg := sync.WaitGroup{}
	if goroutines == 0 {
		goroutines = 1
	}
	wg.Add(goroutines)

	mu := sync.Mutex{}
	filteredPayments := make([]types.Payment, 0)

	lenPayments := len(s.payments)

	for i := 0; i < goroutines; i++ {
		remainder := 0
		if lenPayments % goroutines != 0 {
			remainder = 1
		}

		start := i * lenPayments / goroutines + remainder
		if i == 0 {
			start = 0
		}
		end := (i + 1) * lenPayments / goroutines + remainder
		if end > lenPayments {
			end = lenPayments
		}

		go func(partPayments []*types.Payment) {
			filtered := make([]types.Payment, 0)
			defer wg.Done()

			for _, payment := range partPayments {
				if filter(*payment) {
					filtered = append(filtered, *payment)
				}
			}

			mu.Lock()
			defer mu.Unlock()
			filteredPayments = append(filteredPayments, filtered...)
		}(s.payments[start:end])
	}

	wg.Wait()
	return filteredPayments, nil
}
