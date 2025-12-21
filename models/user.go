package models

import (
	"time"

	role "mobile-directory-bussines/models/role"
)

type User struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	RoleID          uint       `json:"role_id"`
	ReffID          *uint      `json:"reff_id"`
	Name            string     `json:"name" gorm:"type:varchar(255)"`
	Password        string     `json:"password" gorm:"type:varchar(255)"`
	Email           string     `json:"email" gorm:"type:varchar(255);unique"`
	MobileNo        *string    `json:"mobile_no" gorm:"type:varchar(20)"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	Role *role.Role `json:"role,omitempty" gorm:"foreignKey:RoleID;references:ID"`
}
