package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Balance DB MODEL
type Balance struct {
	gorm.Model
	UserId      uuid.UUID `json:"user_id"`
	UserBalance int       `json:"user_balance"`
	Info        string    `json:"info"`
}

type BalanceChangeRequest struct {
	UserId uuid.UUID `json:"user_id"`
	Income int       `json:"income"`
}

type BalanceChangeRespond struct {
	UserId      uuid.UUID `json:"user_id"`
	UserBalance int       `json:"user_balance"`
	Income      int       `json:"income"`
	Info        string    `json:"info"`
}

type BalanceChangeError struct {
	UserId uuid.UUID `json:"user_id"`
	Info   string    `json:"info"`
}

type BalanceRequest struct {
	UserId uuid.UUID `json:"user_id"`
}

type BalanceRespond struct {
	UserId      uuid.UUID `json:"user_id"`
	UserBalance int       `json:"user_balance"`
	Info        string    `json:"info"`
}

type BalanceError struct {
	UserId uuid.UUID `json:"user_id"`
	Info   string    `json:"info"`
}
