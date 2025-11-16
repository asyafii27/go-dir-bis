package bestgroup

import (
	"time"
)

type BestGroup struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Code      *string    `gorm:"type:varchar(20)" json:"code"`
	Name      string     `gorm:"type:varchar(55);not null" json:"name"`
	Rank      *int       `json:"rank"`
	Price     *float64   `gorm:"type:decimal(15,2)" json:"price"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (BestGroup) TableName() string {
	return "best_groups"
}
