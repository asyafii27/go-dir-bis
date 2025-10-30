package partner

import (
	"time"

	partner_owner "mobile-directory-bussines/models/partner_owner"
)

type Partner struct {
	ID                  string     `json:"id" gorm:"type:char(36);primaryKey"`
	PartnerOwnerID      *uint64    `json:"partner_owner_id" gorm:"not null"`
	SalesID             *uint64    `json:"sales_id"`
	BestGroupID         *uint64    `json:"best_group_id"`
	ProductPacketID     *uint64    `json:"product_packet_id"`
	RangeMonth          int        `json:"range_month" gorm:"default:1"`
	QtyMonth            int        `json:"qty_month" gorm:"default:1"`
	Code                string     `json:"code" gorm:"type:varchar(50);unique;not null"`
	Name                string     `json:"name" gorm:"type:varchar(100);not null;index"`
	Email               *string    `json:"email" gorm:"type:varchar(255);unique"`
	MobileNo            string     `json:"mobile_no" gorm:"type:varchar(15);not null"`
	WorkingHourStart    *time.Time `json:"working_hour_start"`
	WorkingHourEnd      *time.Time `json:"working_hour_end"`
	StatusTxt           string     `json:"status_txt" gorm:"type:varchar(20);default:'active'"`
	Description         *string    `json:"description" gorm:"type:text"`
	Address             *string    `json:"address" gorm:"type:text"`
	Rt                  *string    `json:"rt" gorm:"type:varchar(5)"`
	Rw                  *string    `json:"rw" gorm:"type:varchar(5)"`
	Latitude            float64    `json:"latitude" gorm:"not null"`
	Longitude           float64    `json:"longitude" gorm:"not null"`
	ProvinceCode        *string    `json:"province_code" gorm:"type:varchar(10)"`
	CityCode            *string    `json:"city_code" gorm:"type:varchar(10)"`
	DistrictCode        *string    `json:"district_code" gorm:"type:varchar(10)"`
	VillageCode         *string    `json:"village_code" gorm:"type:varchar(10)"`
	PaymentMethod       string     `json:"payment_method" gorm:"type:enum('manual','payment_gateway','not_set');default:'not_set'"`
	PaymentStatus       *string    `json:"payment_status" gorm:"type:varchar(255)"`
	ManualPaymentStatus string     `json:"manual_payment_status" gorm:"type:varchar(100);default:'pending'"`
	UrlWebsite          *string    `json:"url_website" gorm:"type:varchar(255)"`
	FbSosmed            *string    `json:"fb_sosmed" gorm:"type:varchar(155)"`
	IgSosmed            *string    `json:"ig_sosmed" gorm:"type:varchar(155)"`
	TiktokSosmed        *string    `json:"tiktok_sosmed" gorm:"type:varchar(155)"`
	PhoneNo             *string    `json:"phone_no" gorm:"type:varchar(50)"`
	SubDomainRequest    *string    `json:"sub_domain_request" gorm:"type:varchar(255)"`
	CreatedAt           *time.Time `json:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at"`

	PartnerOwner *partner_owner.PartnerOwner `json:"partner_owner,omitempty" gorm:"foreignKey:PartnerOwnerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (Partner) TableName() string {
	return "partners"
}
