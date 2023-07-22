package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`
	Creator     string `gorm:"not null" json:"creator"`
	Asignees    []User `gorm:"many2many:task_asignees;constraint:OnDelete:SET NULL" json:"asignees"`
}
