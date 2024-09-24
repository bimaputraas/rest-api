package model

type User struct {
	ID          uint   `gorm:"type:uuid;primary_key" json:"user_id"`
	FirstName   string `gorm:"type:varchar(100);not null" json:"first_name" validate:"required"`
	LastName    string `gorm:"type:varchar(100);not null" json:"last_name" validate:"required"`
	PhoneNumber string `gorm:"type:varchar(15);unique;not null" json:"phone_number" validate:"required,e164"`
	Address     string `gorm:"type:text" json:"address" validate:"required"`
	Pin         string `gorm:"type:varchar(6);not null" json:"pin,omitempty" validate:"required"`
	Created     string `json:"created"`
}

type Balance struct {
	ID             uint    `gorm:"type:uuid;primary_key" json:"id"`
	UserID         uint    `gorm:"column:user_id;index" json:"user_id"`
	CurrentBalance float64 `gorm:"column:current_balance" json:"current_balance"`
	Updated        string  `json:"updated"`
}
