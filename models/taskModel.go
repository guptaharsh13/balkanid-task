package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Name        string `gorm:"not null" json:"name"`
	Description string `gorm:"not null" json:"description"`
	Creator     string `gorm:"unique;not null;uniqueIndex" json:"creator"`
	Asignees    []User `gorm:"many2many:task_asignees"`
}
