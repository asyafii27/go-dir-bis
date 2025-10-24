package subcategory

import (
	"mobile-directory-bussines/models/master/category"
	"time"
)

type SubCategory struct {
	ID          uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryId  uint64     `gorm:"not null;index" json:"category_id"`
	Code        string     `gorm:"type:varchar(50);not null;unique" json:"code"`
	Name        string     `gorm:"type:varchar(50);not null" json:"name"`
	StatusTxt   string     `gorm:"type:varchar(20);default:'active';not null" json:"status_txt"`
	Description *string    `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`

	Category category.Category `gorm:"foreignKey:CategoryID;references:ID" json:"category"`
}

func (SubCategory) TableName() string {
	return "sub_categories"
}
