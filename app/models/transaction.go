package models

import (
	"time"

	"github.com/gofrs/uuid"
)

const (
	TxTypeDeposit    = "deposit"
	TxTypeWithdrawal = "withdrawal"
	TxTypeBalance    = "balance_enquiry"
)

const (
	TxStatusPending  = "pending"
	TxStatusApproved = "approved"
)

type Transaction struct {
	ID        uuid.UUID `json:"transactionId"`
	Type      string    `json:"transactionType"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Amount    float64   `json:"amount"`
	UserID    uuid.UUID `json:"userId"`
	AccountID uuid.UUID `json:"accountId"`
}
