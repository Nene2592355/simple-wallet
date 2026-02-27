package transaction

import (
	"testing"
	"time"

	"simple-wallet/app/data"
	"simple-wallet/app/errors"
	"simple-wallet/app/models"

	"github.com/gofrs/uuid"
)

// mockRepository implements the Repository interface for testing
type mockRepository struct {
	transactions map[uuid.UUID]models.Transaction
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		transactions: make(map[uuid.UUID]models.Transaction),
	}
}

func (m *mockRepository) Add(tx models.Transaction) (models.Transaction, error) {
	m.transactions[tx.ID] = tx
	return tx, nil
}

func (m *mockRepository) GetByID(id uuid.UUID) (models.Transaction, error) {
	tx, ok := m.transactions[id]
	if !ok {
		return models.Transaction{}, errors.Error{Code: errors.ENOTFOUND}
	}
	return tx, nil
}

func (m *mockRepository) UpdateStatus(id uuid.UUID, status string) error {
	tx, ok := m.transactions[id]
	if !ok {
		return errors.Error{Code: errors.ENOTFOUND}
	}
	tx.Status = status
	m.transactions[id] = tx
	return nil
}

func (m *mockRepository) GetTransactions(userId uuid.UUID, from time.Time, limit int) (*[]models.Transaction, error) {
	var result []models.Transaction
	for _, tx := range m.transactions {
		if tx.UserID == userId {
			result = append(result, tx)
		}
	}
	return &result, nil
}

func TestApproveTransaction_Success(t *testing.T) {
	repo := newMockRepository()
	transChan := data.ChanNewTransactions{
		Channel: make(chan data.TransactionContract, 10),
		Reader:  make(chan data.TransactionContract, 10),
		Writer:  make(chan data.TransactionContract, 10),
	}

	intr := NewInteractor(repo, transChan)

	txID, _ := uuid.NewV4()
	userID, _ := uuid.NewV4()
	accID, _ := uuid.NewV4()

	tx := models.Transaction{
		ID:        txID,
		Type:      models.TxTypeDeposit,
		Status:    models.TxStatusPending,
		Timestamp: time.Now(),
		Amount:    100.0,
		UserID:    userID,
		AccountID: accID,
	}
	repo.transactions[txID] = tx

	err := intr.ApproveTransaction(txID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updated, _ := repo.GetByID(txID)
	if updated.Status != models.TxStatusApproved {
		t.Errorf("expected status %s, got %s", models.TxStatusApproved, updated.Status)
	}
}

func TestApproveTransaction_NotFound(t *testing.T) {
	repo := newMockRepository()
	transChan := data.ChanNewTransactions{
		Channel: make(chan data.TransactionContract, 10),
		Reader:  make(chan data.TransactionContract, 10),
		Writer:  make(chan data.TransactionContract, 10),
	}

	intr := NewInteractor(repo, transChan)

	txID, _ := uuid.NewV4()
	err := intr.ApproveTransaction(txID)
	if err == nil {
		t.Fatal("expected error for non-existent transaction")
	}

	e, ok := err.(errors.Error)
	if !ok {
		t.Fatal("expected errors.Error type")
	}
	if e.Code != errors.ENOTFOUND {
		t.Errorf("expected error code %s, got %s", errors.ENOTFOUND, e.Code)
	}
}

func TestApproveTransaction_AlreadyApproved(t *testing.T) {
	repo := newMockRepository()
	transChan := data.ChanNewTransactions{
		Channel: make(chan data.TransactionContract, 10),
		Reader:  make(chan data.TransactionContract, 10),
		Writer:  make(chan data.TransactionContract, 10),
	}

	intr := NewInteractor(repo, transChan)

	txID, _ := uuid.NewV4()
	userID, _ := uuid.NewV4()
	accID, _ := uuid.NewV4()

	tx := models.Transaction{
		ID:        txID,
		Type:      models.TxTypeDeposit,
		Status:    models.TxStatusApproved,
		Timestamp: time.Now(),
		Amount:    100.0,
		UserID:    userID,
		AccountID: accID,
	}
	repo.transactions[txID] = tx

	err := intr.ApproveTransaction(txID)
	if err == nil {
		t.Fatal("expected error for already approved transaction")
	}

	e, ok := err.(errors.Error)
	if !ok {
		t.Fatal("expected errors.Error type")
	}
	if e.Code != errors.EINVALID {
		t.Errorf("expected error code %s, got %s", errors.EINVALID, e.Code)
	}
}
