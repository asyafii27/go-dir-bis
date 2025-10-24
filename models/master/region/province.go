package region

import "time"

type Province struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Code      string    `gorm:"type:char(2);not null;unique" json:"code"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Meta      *string   `gorm:"type:text" json:"meta,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Province) TableName() string {
	return "indonesia_provinces"
}
