package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/pastorenue/kinance/internal/common"
)

type User struct {
	common.BaseModel
	Email          string     `json:"email" gorm:"uniqueIndex;not null"`
	Password       string     `json:"-" gorm:"not null"`
	FirstName      string     `json:"first_name" gorm:"not null"`
	LastName       string     `json:"last_name" gorm:"not null"`
	Phone          string     `json:"phone"`
	DateOfBirth    *time.Time `json:"date_of_birth"`
	ProfilePicture string     `json:"profile_picture"`
	IsActive       bool       `json:"is_active" gorm:"default:true"`
	FamilyID       *uuid.UUID `json:"family_id"`
	Family         *Family    `json:"family,omitempty"`
	Role           UserRole   `json:"role" gorm:"default:member"`
}

type Family struct {
	common.BaseModel
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	Members     []User `json:"members,omitempty"`
}

type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleParent UserRole = "parent"
	RoleMember UserRole = "member"
	RoleChild  UserRole = "child"
)

type CreateUserRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	Phone           string `json:"phone"`
}

type UpdateUserRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Phone       string `json:"phone"`
	DateOfBirth string `json:"date_of_birth"`
}
