package model

type TopUp struct {
	ID            uint    `gorm:"type:uuid;primary_key" json:"top_up_id"`
	UserID        uint    `gorm:"type:uuid;not null" json:"user_id,omitempty"`
	AmountTopUp   float64 `gorm:"not null" json:"amount_top_up"`
	BalanceBefore float64 `gorm:"not null" json:"balance_before"`
	BalanceAfter  float64 `gorm:"not null" json:"balance_after"`
	Created       string  `gorm:"autoCreateTime" json:"created"`
}
