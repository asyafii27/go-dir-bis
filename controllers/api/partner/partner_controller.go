package partner

import (
	"database/sql"
	"mobile-directory-bussines/config"
	"mobile-directory-bussines/helpers"
	"mobile-directory-bussines/models/partner"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPartners(c *gin.Context) {
	var partners []partner.Partner

	db := config.Database

	db = db.Preload("Categories").Preload("SubCategories").Preload("SecondSubCategories")

	db = ApplyPartnerFilters(c, db)

	page := c.Query("page")

	if page != "" {
		meta, err := helpers.LaravelPaginate(c, db, &partners)
		if err != nil {
			helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, meta)
		return
	}

	if err := db.Find(&partners).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, partners)
}

func ApplyPartnerFilters(c *gin.Context, db *gorm.DB) *gorm.DB {

	if code := c.Query("code"); code != "" {
		db = db.Where("code LIKE ?", "%"+code+"%")
	}

	if name := c.Query("name"); name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	if partnerOwnerID := c.Query("partner_owner_id"); partnerOwnerID != "" {
		db = db.Where("partner_owner_id = ?", partnerOwnerID)
	}

	if mobileNo := c.Query("mobile_no"); mobileNo != "" {
		db = db.Where("mobile_no LIKE ?", "%"+mobileNo+"%")
	}

	if address := c.Query("address"); address != "" {
		db = db.Where("address LIKE ?", "%"+address+"%")
	}

	if provinceCode := c.Query("province_code"); provinceCode != "" {
		db = db.Where("province_code = ?", provinceCode)
	}

	if cityCode := c.Query("city_code"); cityCode != "" {
		db = db.Where("city_code = ?", cityCode)
	}

	if districtCode := c.Query("district_code"); districtCode != "" {
		db = db.Where("district_code = ?", districtCode)
	}

	if VillageCode := c.Query("village_code"); VillageCode != "" {
		db = db.Where("village_code = ?", VillageCode)
	}

	if statusTxt := c.Query("status_txt"); statusTxt != "" {
		db = db.Where("status_txt = ?", statusTxt)
	}

	if verificationTxt := c.Query("verification_seller_status_txt"); verificationTxt != "" {
		statuses := strings.Split(verificationTxt, ",")
		db = db.Where("verification_seller_status_txt IN ?", statuses)
	}

	if globalSearch := c.Query("global_search"); globalSearch != "" {
		like := "%" + globalSearch + "%"

		db = db.Joins("LEFT JOIN partner_owners AS po ON po.id = partners.partner_owner_id").
			Where(`
            (
                partners.code LIKE @like
                OR partners.name LIKE @like
                OR partners.email LIKE @like
                OR partners.mobile_no LIKE @like
                OR po.name LIKE @like
            )
        `, sql.Named("like", like))
	}

	if bestGroupId := c.Query("best_group_id"); bestGroupId != "" {
		explodeBestGroupId := strings.Split(bestGroupId, ",")
		db = db.Where("best_group_id = ?", explodeBestGroupId)
	}

	if paymentMethod := c.Query("payment_method"); paymentMethod != "" {
		db = db.Where("payment_method = ?", paymentMethod)
	}

	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortDir := strings.ToLower(c.DefaultQuery("sort_dir", "desc"))

	allowedSorts := map[string]bool{
		"id":   true,
		"name": true,
	}

	if !allowedSorts[sortBy] {
		sortBy = "created_at"
	}

	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}

	return db.Order(sortBy + " " + sortDir)
}
