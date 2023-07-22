package models

import (
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	Name        string `gorm:"primaryKey;unique;uniqueIndex;not null" json:"name"`
	Description string `json:"description"`
	Table       string `gorm:"not null" json:"table"`
	Operation   string `gorm:"not null" json:"operation"`

	Groups []Group `gorm:"many2many:group_permissions;constraint:OnDelete:SET NULL" json:"groups"`
	Roles  []Role  `gorm:"many2many:role_permissions;constraint:OnDelete:SET NULL" json:"roles"`
}
