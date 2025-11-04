package partnerowner

import (
	region "mobile-directory-bussines/models/master/region"
	"time"
)

type PartnerOwner struct {
	ID           uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Code         string     `json:"code" gorm:"size:50;unique;not null"`
	Name         string     `json:"name" gorm:"size:100;not null"`
	Email        *string    `json:"email" gorm:"size:255;unique"`
	MobileNo     string     `json:"mobile_no" gorm:"size:15;not null"`
	Address      *string    `json:"address" gorm:"type:text"`
	RT           *string    `json:"rt" gorm:"size:5"`
	RW           *string    `json:"rw" gorm:"size:5"`
	Latitude     *float64   `json:"latitude" gorm:"type:decimal(10,7)"`
	Longitude    *float64   `json:"longitude" gorm:"type:decimal(10,7)"`
	ProvinceCode *string    `json:"province_code" gorm:"size:10;index"`
	CityCode     *string    `json:"city_code" gorm:"size:10;index"`
	DistrictCode *string    `json:"district_code" gorm:"size:10;index"`
	VillageCode  *string    `json:"village_code" gorm:"size:10;index"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`

	// relations
	Province *region.Province `json:"province" gorm:"foreignKey:ProvinceCode;references:Code"`
	City     *region.City     `json:"city" gorm:"foreignKey:CityCode;references:Code"`
	District *region.District `json:"district" gorm:"foreignKey:DistrictCode;references:Code"`
	Village  *region.Village  `json:"village" gorm:"foreignKey:VillageCode;references:Code"`
}

func (PartnerOwner) TableName() string {
	return "partner_owners"
}
