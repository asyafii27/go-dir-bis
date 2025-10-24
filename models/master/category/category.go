package category

import "time"

type Category struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Code        string    `gorm:"type:varchar(50);not null;unique" json:"code"`
	Name        string    `gorm:"type:varchar(50)" json:"name"`
	StatusTxt   string    `gorm:"type:varchar(20);default:active" json:"active"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Category) TableName() string {
	return "categories"
}
