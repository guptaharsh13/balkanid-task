package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string `gorm:"primaryKey;unique;uniqueIndex;not null" json:"name"`
	Description string `json:"description"`

	Users       []User       `gorm:"foreignKey:Role;references:Name;constraint:OnDelete:SET NULL" json:"users"`
	Permissions []Permission `gorm:"many2many:role_permissions;constraint:OnDelete:SET NULL" json:"permissions"`
}
