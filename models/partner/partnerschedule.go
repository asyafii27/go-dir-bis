package partner

import (
	"time"
)

type PartnerSchedule struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	PartnerID string    `gorm:"type:char(36);not null" json:"partner_id"`
	Type      string    `gorm:"type:enum('open','close');not null;default:'open'" json:"type"`
	DayOfWeek string    `gorm:"type:enum('1','2','3','4','5','6','7');not null" json:"day_of_week"`
	OpenTime  *string   `gorm:"type:string" json:"open_time"`
	CloseTime *string   `gorm:"type:string" json:"close_time"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (PartnerSchedule) TableName() string {
	return "partner_schedules"
}
