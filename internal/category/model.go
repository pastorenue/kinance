package category

import (
	"github.com/google/uuid"
	"github.com/pastorenue/kinance/internal/common"
)

type Category struct {
	common.BaseModel
	Name             string     `gorm:"not null;unique" json:"name"`
	UserID           uuid.UUID  `gorm:"not null;index" json:"user_id"`
	ParentCategoryID *uuid.UUID `gorm:"index" json:"parent_category_id,omitempty"`
	ParentCategory   *Category  `gorm:"foreignKey:ParentCategoryID" json:"parent_category"`
	SubCategories    []Category `gorm:"foreignKey:ParentCategoryID" json:"sub_categories"`
	ColorCode        string     `gorm:"type:char(7)" json:"color_code"` // e.g., #RRGGBB
}

type CreateCategoryRequest struct {
	Name             string     `json:"name" binding:"required"`
	ParentCategoryID *uuid.UUID `json:"parent_category_id,omitempty"`
	ColorCode        string     `json:"color_code" binding:"omitempty,hexcolor|len=7"`
}

type UpdateCategoryRequest struct {
	Name             *string    `json:"name"`
	ParentCategoryID *uuid.UUID `json:"parent_category_id,omitempty"`
	ColorCode        *string    `json:"color_code" binding:"omitempty,hexcolor|len=7"`
}
