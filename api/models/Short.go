package models

import (
	"time"

	"gorm.io/gorm"
)

type Short struct {
	gorm.Model

	Pagina        string `gorm:"not null"`
	Short         string `gorm:"not null"`
	Expiry        time.Time
	FechaCreacion time.Time `json:"fecha_creacion"`
	Abierto       int
	Permisos      string
	UserID        uint
}
