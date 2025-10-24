package region

import "time"

type City struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Code         string    `gorm:"type:char(2);not null;unique" json:"code"`
	ProvinceCode string    `gorm:"type:varchar" json:"province_code"`
	Meta         *string   `gorm:"type:text" json:"meta.omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (City) TableName() string {
	return "indonesia_cities"
}
