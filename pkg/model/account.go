package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Account DB model
type Account struct {
	gorm.Model
	UserId  uuid.UUID `json:"user_id"`
	ItemID  uuid.UUID `json:"item_id"`
	OrderID uuid.UUID `json:"order_id"`
	Price   int       `json:"price"`
	Info    string    `json:"info"`
}

type AccountRequest struct {
	AccountOrder Account
}

type AccountRespond struct {
	UserId      uuid.UUID `json:"user_id"`
	ItemID      uuid.UUID `json:"item_id"`
	OrderID     uuid.UUID `json:"order_id"`
	UserBalance int       `json:"user_balance"`
}

type AccountRespondError struct {
	UserId  uuid.UUID `json:"user_id"`
	ItemID  uuid.UUID `json:"item_id"`
	OrderID uuid.UUID `json:"order_id"`
	Info    string    `json:"info"`
}
