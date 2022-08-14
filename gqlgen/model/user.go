package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
