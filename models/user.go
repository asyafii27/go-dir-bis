package models

import (
	"time"

	role "mobile-directory-bussines/models/role"
)

type User struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	RoleID          uint       `json:"role_id"`
	ReffID          *uint      `json:"reff_id"`
	Name            string     `json:"name"`
	Password        string     `json:"password"`
	Email           string     `json:"email"`
	MobileNo        *string    `json:"mobile_no"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	Role *role.Role `json:"role,omitempty" gorm:"foreignKey:RoleID;references:ID"`
}
