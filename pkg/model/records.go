package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Records struct {
	gorm.Model
	UserID  uuid.UUID
	ItemID  uuid.UUID
	OrderID uuid.UUID
	Income  int
}
