package models

import (
	"time"
)

// Создадим структуры для таблиц
// Структура для песни
type Song struct {
	ID        int        `gorm:"primary_key"`
	Duration  int        `json:"duration"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}
