package model

import (
	"math/rand"
	"time"
)

type Transaction struct {
	ID              uint    `gorm:"type:uuid;primary_key" json:"transaction_id,omitempty"`
	UserID          uint    `gorm:"type:uuid;not null" json:"user_id"`
	TopUpId         *uint   `gorm:"type:uuid" json:"top_up_id,omitempty"`
	TransferId      *uint   `gorm:"type:uuid" json:"transfer_id,omitempty"`
	PaymentId       *uint   `gorm:"type:uuid" json:"payment_id,omitempty"`
	TransactionType string  `gorm:"type:varchar(10);not null" json:"transaction_type"` // CREDIT or DEBIT
	Amount          float64 `gorm:"not null" json:"amount"`
	Remarks         string  `gorm:"type:varchar(255)" json:"remarks"`
	BalanceBefore   float64 `gorm:"not null" json:"balance_before"`
	BalanceAfter    float64 `gorm:"not null" json:"balance_after"`
	Status          string  `gorm:"type:varchar(10);not null" json:"status"`
	Created         string  `json:"created"`
}

func (t *Transaction) MockRandType() {

	rand.Seed(time.Now().UnixNano())
	if rand.Intn(2) == 0 {
		t.TransactionType = "DEBIT"
		return
	}
	t.TransactionType = "KREDIT"
}
