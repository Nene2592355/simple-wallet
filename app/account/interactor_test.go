package account

import (
	"testing"

	"simple-wallet/app/data"
	"simple-wallet/app/errors"
	"simple-wallet/app/models"

	"github.com/gofrs/uuid"
)

// mockRepository is a test double for the Repository interface.
type mockRepository struct {
	account models.Account
	err     error
}

func (m *mockRepository) GetAccountByUserID(_ uuid.UUID) (models.Account, error) {
	return m.account, m.err
}

func (m *mockRepository) UpdateBalance(amount uint, _ uuid.UUID) (models.Account, error) {
	if m.err != nil {
		return models.Account{}, m.err
	}
	updated := m.account
	updated.AvailableBalance = amount
	return updated, nil
}

func (m *mockRepository) Create(_ uuid.UUID) (models.Account, error) {
	return m.account, m.err
}

func newTestInteractor(repo Repository) Interactor {
	usersChan := make(chan data.UserContract, 1)
	transChan := make(chan data.TransactionContract, 10)

	chanUsers := data.ChanNewUsers{
		Channel: usersChan,
		Reader:  usersChan,
		Writer:  usersChan,
	}
	chanTrans := data.ChanNewTransactions{
		Channel: transChan,
		Reader:  transChan,
		Writer:  transChan,
	}

	return NewInteractor(repo, chanUsers, chanTrans)
}

func TestDeposit_AmountBelowMinimum(t *testing.T) {
	userID, _ := uuid.NewV4()
	repo := &mockRepository{}
	intr := newTestInteractor(repo)

	_, err := intr.Deposit(userID, minimumDepositAmount-1)
	if err == nil {
		t.Fatal("expected error for amount below minimum, got nil")
	}
}

func TestDeposit_AccountNotFound(t *testing.T) {
	userID, _ := uuid.NewV4()
	repo := &mockRepository{
		err: errors.Error{Code: errors.ENOTFOUND},
	}
	intr := newTestInteractor(repo)

	_, err := intr.Deposit(userID, minimumDepositAmount)
	if err == nil {
		t.Fatal("expected error when account not found, got nil")
	}
}

func TestDeposit_AccountFrozen(t *testing.T) {
	userID, _ := uuid.NewV4()
	accID, _ := uuid.NewV4()
	repo := &mockRepository{
		account: models.Account{
			ID:               accID,
			UserID:           userID,
			Status:           models.StatusFrozen,
			AvailableBalance: 0,
		},
	}
	intr := newTestInteractor(repo)

	_, err := intr.Deposit(userID, minimumDepositAmount)
	if err == nil {
		t.Fatal("expected error for frozen account, got nil")
	}
}

func TestDeposit_Success(t *testing.T) {
	userID, _ := uuid.NewV4()
	accID, _ := uuid.NewV4()
	initialBalance := uint(0)
	repo := &mockRepository{
		account: models.Account{
			ID:               accID,
			UserID:           userID,
			Status:           models.StatusActive,
			AvailableBalance: initialBalance,
		},
	}
	intr := newTestInteractor(repo)

	depositAmount := uint(100)
	balance, err := intr.Deposit(userID, depositAmount)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// balance is stored in cents; convert to dollars for comparison
	expected := float64(initialBalance+depositAmount*100) / 100.0
	if balance != expected {
		t.Errorf("expected balance %v, got %v", expected, balance)
	}
}
