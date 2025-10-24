package secondsubcategory

import (
	"mobile-directory-bussines/models/master/category"
	"mobile-directory-bussines/models/master/subcategory"
	"time"
)

type SecondSubCategory struct {
	ID            uint64     `gorm:"primaryKey:autoIncrement" json:"id"`
	CategoryId    uint64     `gorm:"not null;index" json:"category_id"`
	SubCategoryId uint64     `gorm:"not null;index" json:"sub_category_id"`
	Code          string     `gorm:"type:varchar(50);not null;unique" json:"code"`
	Name          string     `gorm:"type:varchar(50);not null;unique" json:"name"`
	StatusTxt     string     `gorm:"type:varchar(20);default:'active';not null" json:"status_txt"`
	Description   *string    `gorm:"type:text" json:"description,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`

	Category    category.Category       `gorm:"foreignKey:CategoryId:references:ID" json:"category"`
	SubCategory subcategory.SubCategory `gorm:"foreignKey:SubCategoryId;references:ID" json:"subcategory"`
}

func (SecondSubCategory) TableName() string {
	return "second_sub_categories"
}
