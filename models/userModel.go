package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"primaryKey;unique;uniqueIndex;not null" json:"username"`
	Email    string `gorm:"unique;uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"password"`
	IsActive bool   `gorm:"default:false" json:"is_active"`
	IsAdmin  bool   `gorm:"default:false" json:"is_admin"`

	Role   string  `gorm:"default:NULL" json:"role"`
	Groups []Group `gorm:"many2many:user_groups;constraint:OnDelete:SET NULL" json:"groups"`

	CreatedTasks  []Task `gorm:"foreignKey:Creator;references:Username;constraint:OnDelete:SET NULL" json:"created_tasks"`
	AssignedTasks []Task `gorm:"many2many:task_asignees;constraint:OnDelete:SET NULL" json:"assigned_tasks"`
}
