package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	UserName       string `gorm:"not null"`
	Password       string `gorm:"not null"`
	Email          string `gorm:"not null"`
	RateLimit      int
	RateLimitReset int
	Permisos       string `gorm:"not null"`
	Token          string
	Short          []Short
	Datos          []Datos
}
