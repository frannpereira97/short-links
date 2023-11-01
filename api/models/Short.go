package models

import "gorm.io/gorm"

type Short struct {
	gorm.Model

	Pagina        string `gorm:"not null"`
	Short         string `gorm:"not null"`
	Expiry        string
	FechaCreacion string
	Abierto       int
	UserID        uint
}
