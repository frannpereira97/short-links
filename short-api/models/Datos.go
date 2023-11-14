package models

import "gorm.io/gorm"

type Datos struct {
	gorm.Model

	Nombre          string `gorm:"not null"`
	Apellido        string `gorm:"not null"`
	Sexo            string
	FechaNacimiento string
	Nacionalidad    string
	Provincia       string
	Ciudad          string
	Domicilio       string
	UserID          uint
}
