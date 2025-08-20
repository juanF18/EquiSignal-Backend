package models

import (
	"time"

	"github.com/google/uuid"
)

type Stock struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Ticker     string    `gorm:"column:ticker"`
	Company    string    `gorm:"column:company"`
	Brokerage  string    `gorm:"column:brokerage"`
	Action     string    `gorm:"column:action"`
	RatingFrom string    `gorm:"column:rating_from"`
	RatingTo   string    `gorm:"column:rating_to"`
	TargetFrom string    `gorm:"column:target_from"`
	TargetTo   string    `gorm:"column:target_to"`
	Time       time.Time `gorm:"column:time"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
