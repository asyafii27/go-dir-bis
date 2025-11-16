package partner

import (
	"time"
)

type HistoryPreferedPartnerVerification struct {
	ID            uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PartnerID     string     `gorm:"type:char(36);not null" json:"partner_id"`
	VerificatorID *uint64    `gorm:"type:bigint unsigned" json:"verificator_id"`
	StatusTxt     string     `gorm:"type:enum('pending','verified','rejected');not null;default:'pending'" json:"status_txt"`
	RejectReason  *string    `gorm:"type:varchar(255)" json:"reject_reason"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

func (HistoryPreferedPartnerVerification) TableName() string {
	return "history_prefered_partner_verifications"
}
