package partnerowner

import (
	"mobile-directory-bussines/config"
	"mobile-directory-bussines/helpers"
	"mobile-directory-bussines/models/partnerowner"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPartnerOwners(c *gin.Context) {
	var partnerowners []partnerowner.PartnerOwner
	db := config.Database

	db = db.Preload("Province").Preload("City").Preload("District").Preload("Village")

	db = ApplyPartnerOwnerFilters(c, db)

	page := c.Query("page")

	if page != "" {
		meta, err := helpers.LaravelPaginate(c, db, &partnerowners)
		if err != nil {
			helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, meta)
		return
	}

	if err := db.Find(&partnerowners).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, partnerowners)
}

func ApplyPartnerOwnerFilters(c *gin.Context, db *gorm.DB) *gorm.DB {
	if globalSearch := c.Query("global_search"); globalSearch != "" {
		db = db.Where("name LIKE ? OR code LIKE ?", "%"+globalSearch+"%", "%"+globalSearch+"%")
	}

	if address := c.Query("address"); address != "" {
		db = db.Where("address = ?", address)
	}

	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortDir := strings.ToLower(c.DefaultQuery("sort_dir", "desc"))

	allowedSorts := map[string]bool{
		"id":         true,
		"name":       true,
		"created_at": true,
		"updated_at": true,
	}

	if !allowedSorts[sortBy] {
		sortBy = "created_at"
	}

	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}

	return db.Order(sortBy + " " + sortDir)
}

func PaginateData(c *gin.Context, db *gorm.DB, result interface{}) (helpers.PaginationMeta, error) {
	meta, err := helpers.Paginate(c, db, result)
	if err != nil {
		return meta, err
	}

	return meta, nil
}

func GetPartnerOwnerByID(c *gin.Context) {
	id := c.Param("id")
	var partnerowner partnerowner.PartnerOwner

	db := config.Database

	db = db.Preload("Province").Preload("City").Preload("District").Preload("Village")

	if err := db.First(&partnerowner, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ErrorResponse(c, http.StatusNotFound, "Data pemilik tidak ditemukan")
		} else {
			helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}

		return
	}

	helpers.SuccessResponse(c, http.StatusOK, "Success", partnerowner)
}

// type UpdatePartnerOwnerRequest struct {
// 	Name            string  `json:"name" binding:"required,max=100"`
// 	UserName        string  `json:"user_name" binding:"required,max=50"`
// 	Email           string  `json:"email" binding:"required,email"`
// 	Password        *string `json:"password" binding:"omitempty,min=8"`
// 	PasswordConfirm *string `json:"password_confirmation" binding:"omitempty,min=8,eqfield=Password"`
// 	MobileNo        *string `json:"mobile_no" binding:"omitempty,max=20"`
// 	Address         *string `json:"address" binding:"omitempty,max=255"`
// 	RT              *string `json:"rt" binding:"omitempty,len=3"`
// 	RW              *string `json:"rw" binding:"omitempty,len=3"`
// 	ProvinceCode    *string `json:"province_code"`
// 	CityCode        *string `json:"city_code"`
// 	DistrictCode    *string `json:"district_code"`
// 	VillageCode     *string `json:"village_code"`
// }

// // validateEmailUniqueness checks if email is unique in partner_owners table
// func validateEmailUniqueness(tx *gorm.DB, email string, excludeID string) error {
// 	var emailCount int64
// 	if err := tx.Model(&partnerowner.PartnerOwner{}).
// 		Where("email = ? AND id != ?", email, excludeID).
// 		Count(&emailCount).Error; err != nil {
// 		return err
// 	}
// 	if emailCount > 0 {
// 		return fmt.Errorf("email sudah digunakan")
// 	}
// 	return nil
// }

// // validateUserEmailUniqueness checks if email is unique in users table
// func validateUserEmailUniqueness(tx *gorm.DB, email string, excludeID uint64) error {
// 	var userEmailCount int64
// 	if err := tx.Model(&models.User{}).
// 		Where("email = ? AND id != ?", email, excludeID).
// 		Count(&userEmailCount).Error; err != nil {
// 		return err
// 	}
// 	if userEmailCount > 0 {
// 		return fmt.Errorf("email sudah digunakan oleh user lain")
// 	}
// 	return nil
// }

// // updatePartnerOwnerData updates the partner owner model with new data
// func updatePartnerOwnerData(partnerOwner *partnerowner.PartnerOwner, request *UpdatePartnerOwnerRequest) {
// 	partnerOwner.Name = strings.Title(request.Name)
// 	partnerOwner.Email = &request.Email
// 	if request.MobileNo != nil {
// 		partnerOwner.MobileNo = *request.MobileNo
// 	}
// 	partnerOwner.Address = request.Address
// 	partnerOwner.RT = request.RT
// 	partnerOwner.RW = request.RW
// 	partnerOwner.ProvinceCode = request.ProvinceCode
// 	partnerOwner.CityCode = request.CityCode
// 	partnerOwner.DistrictCode = request.DistrictCode
// 	partnerOwner.VillageCode = request.VillageCode
// }

// // updateUserData updates the associated user data
// func updateUserData(tx *gorm.DB, user *models.User, partnerOwner *partnerowner.PartnerOwner, request *UpdatePartnerOwnerRequest) error {
// 	user.Name = partnerOwner.Name
// 	user.UserName = request.UserName
// 	user.Email = request.Email
// 	user.MobileNo = partnerOwner.MobileNo

// 	if request.Password != nil {
// 		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*request.Password), bcrypt.DefaultCost)
// 		if err != nil {
// 			return err
// 		}
// 		user.Password = string(hashedPassword)
// 	}

// 	return tx.Save(user).Error
// }

// func createNewUser(tx *gorm.DB, partnerOwner *partnerowner.PartnerOwner, request *UpdatePartnerOwnerRequest) error {
// 	if request.Password == nil {
// 		return fmt.Errorf("password is required for new user")
// 	}

// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*request.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return err
// 	}

// 	newUser := models.User{
// 		ReffID:   partnerOwner.ID,
// 		RoleID:   4,
// 		Name:     partnerOwner.Name,
// 		UserName: request.UserName,
// 		Email:    request.Email,
// 		MobileNo: partnerOwner.MobileNo,
// 		Password: string(hashedPassword),
// 	}

// 	return tx.Create(&newUser).Error
// }

// func UpdatePartnerOwner(c *gin.Context) {
// 	id := c.Param("id")
// 	var request UpdatePartnerOwnerRequest

// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		helpers.ErrorResponse(c, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	db := config.Database
// 	tx := db.Begin()
// 	if tx.Error != nil {
// 		helpers.ErrorResponse(c, http.StatusInternalServerError, tx.Error.Error())
// 		return
// 	}
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 		}
// 	}()

// 	var partnerOwner partnerowner.PartnerOwner
// 	if err := tx.First(&partnerOwner, id).Error; err != nil {
// 		tx.Rollback()
// 		if err == gorm.ErrRecordNotFound {
// 			helpers.ErrorResponse(c, http.StatusNotFound, "Data tidak ditemukan")
// 		} else {
// 			helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
// 		}
// 		return
// 	}

// 	if err := validateEmailUniqueness(tx, request.Email, id); err != nil {
// 		tx.Rollback()
// 		helpers.ErrorResponse(c, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	updatePartnerOwnerData(&partnerOwner, &request)
// 	if err := tx.Save(&partnerOwner).Error; err != nil {
// 		tx.Rollback()
// 		helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	var user models.User
// 	err := tx.Where("reff_id = ?", partnerOwner.ID).First(&user).Error

// 	if err == nil { // User exists
// 		if err := validateUserEmailUniqueness(tx, request.Email, user.ID); err != nil {
// 			tx.Rollback()
// 			helpers.ErrorResponse(c, http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		if err := updateUserData(tx, &user, &partnerOwner, &request); err != nil {
// 			tx.Rollback()
// 			helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 	} else if err == gorm.ErrRecordNotFound {
// 		if err := createNewUser(tx, &partnerOwner, &request); err != nil {
// 			tx.Rollback()
// 			helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 	} else {
// 		tx.Rollback()
// 		helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	if err := tx.Commit().Error; err != nil {
// 		helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	helpers.SuccessResponse(c, http.StatusOK, "Success", partnerOwner)
// }
