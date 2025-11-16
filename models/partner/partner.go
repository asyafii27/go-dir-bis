package partner

import (
	"time"

	"mobile-directory-bussines/models/master/bestgroup"
	"mobile-directory-bussines/models/master/category"
	"mobile-directory-bussines/models/master/region"
	"mobile-directory-bussines/models/master/secondsubcategory"
	"mobile-directory-bussines/models/master/subcategory"
	"mobile-directory-bussines/models/partnerowner"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Partner struct {
	ID                          string     `gorm:"type:char(36);primaryKey" json:"id"`
	PartnerOwnerID              string     `gorm:"type:char(36)" json:"partner_owner_id"`
	SellerType                  string     `json:"seller_type"` // general_partner / prefered_partner
	SalesID                     string     `json:"sales_id"`
	BestGroupID                 string     `json:"best_group_id"`
	VerificatorID               string     `json:"verificator_id"`
	RangeMonth                  string     `json:"range_month"`
	QtyMonth                    int        `json:"qty_month"`
	Code                        string     `json:"code"`
	Name                        string     `json:"name"`
	Email                       string     `json:"email"`
	MobileNo                    string     `json:"mobile_no"`
	WorkingHourStart            string     `json:"working_hour_start"`
	WorkingHourEnd              string     `json:"working_hour_end"`
	Address                     string     `json:"address"`
	Rt                          string     `json:"rt"`
	Rw                          string     `json:"rw"`
	Latitude                    string     `json:"latitude"`
	Longitude                   string     `json:"longitude"`
	ProvinceCode                string     `json:"province_code"`
	CityCode                    string     `json:"city_code"`
	DistrictCode                string     `json:"district_code"`
	VillageCode                 string     `json:"village_code"`
	PaymentMethod               string     `json:"payment_method"`
	PaymentStatus               string     `json:"payment_status"`
	ManualPaymentStatus         string     `json:"manual_payment_status"`
	UrlWebsite                  string     `json:"url_website"`
	SubDomainRequest            string     `json:"sub_domain_request"`
	StatusTxt                   string     `json:"status_txt"`
	VerificationSellerStatusTxt string     `json:"verification_seller_status_txt"`
	VerificationSellerAt        *time.Time `json:"verification_seller_at"`
	StatusStreamTxt             string     `json:"status_stream_txt"`
	FbSosmed                    string     `json:"fb_sosmed"`
	IgSosmed                    string     `json:"ig_sosmed"`
	TiktokSosmed                string     `json:"tiktok_sosmed"`
	PhoneNo                     string     `json:"phone_no"`
	StoreDescription            string     `json:"store_description"`
	BussinesType                string     `json:"bussines_type"`
	CreatedAt                   time.Time  `json:"created_at"`
	UpdatedAt                   time.Time  `json:"updated_at"`

	PartnerOwner *partnerowner.PartnerOwner `gorm:"foreignKey:PartnerOwnerID" json:"partner_owner"`
	Province     *region.Province           `gorm:"foreignKey:ProvinceCode;references:Code" json:"province"`
	City         *region.City               `gorm:"foreignKey:CityCode;references:Code" json:"city"`
	District     *region.District           `gorm:"foreignKey:DistrictCode;references:Code" json:"district"`
	Village      *region.Village            `gorm:"foreignKey:VillageCode;references:Code" json:"village"`

	Categories          []*category.Category                   `gorm:"many2many:partners_categories;foreignKey:ID;joinForeignKey:PartnerID;References:ID;joinReferences:CategoryID" json:"categories"`
	SubCategories       []*subcategory.SubCategory             `gorm:"many2many:partners_sub_categories;foreignKey:ID;joinForeignKey:PartnerID;References:ID;joinReferences:SubCategoryID" json:"sub_categories"`
	SecondSubCategories []*secondsubcategory.SecondSubCategory `gorm:"many2many:partners_second_sub_categories;foreignKey:ID;joinForeignKey:PartnerID;References:ID;joinReferences:SecondSubCategoryID" json:"second_sub_categories"`

	BestGroup                           *bestgroup.BestGroup                  `gorm:"foreignKey:BestGroupID" json:"best_group"`
	HistoryPreferedPartnerVerifications []*HistoryPreferedPartnerVerification `gorm:"foreignKey:PartnerID;references:ID" json:"history_prefered_partner_verifications"`

	PartnerScheduleMonday    *PartnerSchedule `gorm:"foreignKey:PartnerID;references:ID;where:day_of_week = '1'" json:"partner_schedule_monday"`
	PartnerScheduleTuesday   *PartnerSchedule `gorm:"foreignKey:PartnerID;references:ID;where:day_of_week = '2'" json:"partner_schedule_tuesday"`
	PartnerScheduleWednesday *PartnerSchedule `gorm:"foreignKey:PartnerID;references:ID;where:day_of_week = '3'" json:"partner_schedule_wednesday"`
	PartnerScheduleThursday  *PartnerSchedule `gorm:"foreignKey:PartnerID;references:ID;where:day_of_week = '4'" json:"partner_schedule_thursday"`
	PartnerScheduleFriday    *PartnerSchedule `gorm:"foreignKey:PartnerID;references:ID;where:day_of_week = '5'" json:"partner_schedule_friday"`
	PartnerScheduleSaturday  *PartnerSchedule `gorm:"foreignKey:PartnerID;references:ID;where:day_of_week = '6'" json:"partner_schedule_saturday"`
	PartnerScheduleSunday    *PartnerSchedule `gorm:"foreignKey:PartnerID;references:ID;where:day_of_week = '7'" json:"partner_schedule_sunday"`
}

// Auto generate UUID before create
func (p *Partner) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return
}
