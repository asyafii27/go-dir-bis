package best_group

import "time"

type BestGroup struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Code      string    `gorm:"column:code;size:20" json:"code"`
	Name      string    `gorm:"column:name;size:55;not null" json:"name"`
	Rank      *int      `gorm:"column:rank" json:"rank"`
	Price     *float64  `gorm:"column:price;type:decimal(15,2)" json:"price"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (BestGroup) TableName() string {
	return "best_groups"
}
