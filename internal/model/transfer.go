package model

type Transfer struct {
	ID            uint    `gorm:"type:uuid;primary_key" json:"transfer_id"`
	UserID        uint    `gorm:"type:uuid;not null" json:"user_id,omitempty"`
	Amount        float64 `gorm:"not null" json:"amount"`
	Remarks       string  `gorm:"type:varchar(255)" json:"remarks"`
	BalanceBefore float64 `gorm:"not null" json:"balance_before"`
	BalanceAfter  float64 `gorm:"not null" json:"balance_after"`
	Created       string  `json:"created"`
}
