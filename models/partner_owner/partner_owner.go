package partner_owner

import "time"

type PartnerOwner struct {
	ID           uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Code         string     `json:"code" gorm:"size:50;unique;not null"`
	Name         string     `json:"name" gorm:"size:100;not null"`
	Email        *string    `json:"email,omitempty" gorm:"size:255;unique"`
	MobileNo     string     `json:"mobile_no" gorm:"size:15;not null"`
	Address      *string    `json:"address,omitempty" gorm:"type:text"`
	RT           *string    `json:"rt,omitempty" gorm:"size:5"`
	RW           *string    `json:"rw,omitempty" gorm:"size:5"`
	Latitude     *float64   `json:"latitude,omitempty" gorm:"type:decimal(10,7)"`
	Longitude    *float64   `json:"longitude,omitempty" gorm:"type:decimal(10,7)"`
	ProvinceCode *string    `json:"province_code,omitempty" gorm:"size:10;index"`
	CityCode     *string    `json:"city_code,omitempty" gorm:"size:10;index"`
	DistrictCode *string    `json:"district_code,omitempty" gorm:"size:10;index"`
	VillageCode  *string    `json:"village_code,omitempty" gorm:"size:10;index"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

func (PartnerOwner) TableName() string {
	return "partner_owners"
}
