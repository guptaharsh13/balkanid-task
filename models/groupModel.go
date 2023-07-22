package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Name        string `gorm:"primaryKey;unique;uniqueIndex;not null" json:"name"`
	Description string `json:"description"`

	Users       []User       `gorm:"many2many:user_groups;constraint:OnDelete:SET NULL" json:"users"`
	Permissions []Permission `gorm:"many2many:group_permissions;constraint:OnDelete:SET NULL" json:"permissions"`
}
