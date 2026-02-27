package transaction

import (
	"time"

	"simple-wallet/app/errors"
	"simple-wallet/app/models"
	"simple-wallet/app/storage"

	"github.com/gofrs/uuid"
)

type Repository interface {
	Add(models.Transaction) (models.Transaction, error)
	GetByID(id uuid.UUID) (models.Transaction, error)
	UpdateStatus(id uuid.UUID, status string) error
	GetTransactions(userId uuid.UUID, from time.Time, limit int) (*[]models.Transaction, error)
}

type repository struct {
	database *storage.Database
}

func NewRepository(db *storage.Database) Repository {
	return &repository{db}
}

func (r repository) Add(tx models.Transaction) (models.Transaction, error) {
	result := r.database.Create(&tx)
	if err := result.Error; err != nil {
		return models.Transaction{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return tx, nil
}

func (r repository) GetByID(id uuid.UUID) (models.Transaction, error) {
	var tx models.Transaction
	result := r.database.Where("id = ?", id).First(&tx)
	if err := result.Error; err != nil {
		return models.Transaction{}, errors.Error{Err: err, Code: errors.ENOTFOUND}
	}
	return tx, nil
}

func (r repository) UpdateStatus(id uuid.UUID, status string) error {
	result := r.database.Model(&models.Transaction{}).Where("id = ?", id).Update("status", status)
	if err := result.Error; err != nil {
		return errors.Error{Err: err, Code: errors.EINTERNAL}
	}
	return nil
}

func (r repository) GetTransactions(userId uuid.UUID, from time.Time, limit int) (*[]models.Transaction, error) {
	var transactions []models.Transaction

	result := r.database.Where(
		models.Transaction{UserID: userId},
	).Where(
		"timestamp <= ?", from,
	).Order("timestamp desc").Limit(limit).Find(&transactions)

	if err := result.Error; err != nil {
		return nil, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return &transactions, nil
}
