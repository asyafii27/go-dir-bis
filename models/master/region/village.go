package region

import "time"

type Village struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Code         string    `gorm:"type:char(20);not null;unique" json:"code"`
	DistrictCode string    `gorm:"type:varchar" json:"district_code"`
	Meta         *string   `gorm:"type:text" json:"meta.omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (Village) TableName() string {
	return "indonesia_villages"
}
