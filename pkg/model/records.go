package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Records DB model
type Records struct {
	gorm.Model
	UserId      uuid.UUID `json:"user_id"`
	ItemID      uuid.UUID `json:"item_id"`
	OrderID     uuid.UUID `json:"order_id"`
	IncomePrice int       `json:"income_price"`
}

type RecordsRequest struct {
	gorm.Model
	UserId      uuid.UUID `json:"user_id"`
	ItemID      uuid.UUID `json:"item_id"`
	OrderID     uuid.UUID `json:"order_id"`
	IncomePrice int       `json:"income_price"`
}

type RecordsRespond struct {
	Status bool   `json:"status"`
	Info   string `json:"info"`
}
