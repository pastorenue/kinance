package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/pastorenue/kinance/internal/common"
)

type Group struct {
	common.BaseModel
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	// Add other fields as needed
}

type Address struct {
	common.BaseModel
	Country     string `json:"country" gorm:"type:varchar(100)"`
	City        string `json:"city" gorm:"type:varchar(100)"`
	PostalCode  string `json:"postal_code" gorm:"type:varchar(20)"`
	Street      string `json:"street" gorm:"type:varchar(100)"`
	HouseNumber int    `json:"house_number" gorm:"type:int"`
}

type User struct {
	common.BaseModel
	Email          string     `json:"email" gorm:"uniqueIndex;not null"`
	Password       string     `json:"-"`
	FirstName      string     `json:"first_name" gorm:"not null"`
	LastName       string     `json:"last_name" gorm:"not null"`
	Phone          string     `json:"phone"`
	DateOfBirth    *time.Time `json:"date_of_birth"`
	ProfilePicture string     `json:"profile_picture"`
	IsActive       bool       `json:"is_active" gorm:"default:true"`
	FamilyID       *uuid.UUID `json:"family_id"`
	Family         *Family    `json:"family,omitempty"`
	Role           UserRole   `json:"role" gorm:"default:member"`
	Address        []Address  `json:"address" gorm:"type:jsonb"`
	Groups         []Group    `json:"groups" gorm:"many2many:user_groups;"`

	// For Super Admins Only
	IsSuperAdmin bool `json:"is_super_admin" gorm:"default:false"`
	IsStaff      bool `json:"is_staff" gorm:"default:false"`
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

type CreateAddressRequest struct {
	Country     string `json:"country" binding:"required"`
	City        string `json:"city" binding:"required"`
	PostalCode  string `json:"postal_code" binding:"required"`
	Street      string `json:"street" binding:"required"`
	HouseNumber int    `json:"house_number" binding:"required"`
}

type UpdateAddressRequest struct {
	ID          uuid.UUID `json:"id" binding:"required"`
	Country     string    `json:"country"`
	City        string    `json:"city"`
	PostalCode  string    `json:"postal_code"`
	Street      string    `json:"street"`
	HouseNumber int       `json:"house_number"`
}
