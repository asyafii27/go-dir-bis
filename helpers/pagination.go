package helpers

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func Paginate(c *gin.Context, db *gorm.DB, out interface{}) (PaginationMeta, error) {
	var meta PaginationMeta

	pageStr := c.Query("page")
	if pageStr == "" {
		if err := db.Find(out).Error; err != nil {
			return meta, err
		}
		return meta, nil
	}

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	var total int64
	if err := db.Model(out).Count(&total).Error; err != nil {
		return meta, err
	}

	if err := db.Limit(limit).Offset(offset).Find(out).Error; err != nil {
		return meta, err
	}

	meta = PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	}

	return meta, nil
}

type LarvaelPaginationMeta struct {
	CurrentPage  int         `json:"current_page"`
	Data         interface{} `json:"data"`
	FirstPageURL string      `json:"first_page_url"`
	From         int         `json:"from"`
	LastPage     int         `json:"last_page"`
	LastPageURL  string      `json:"last_page_url"`
	Links        []PageLink  `json:"links"`
	NextPageURL  string      `json:"next_page_url"`
	Path         string      `json:"path"`
	PerPage      int         `json:"per_page"`
	PrevPageURL  string      `json:"prev_page_url"`
	To           int         `json:"to"`
	Total        int64       `json:"total"`
}
type PageLink struct {
	URL    *string `json:"url"`
	Label  string  `json:"label"`
	Page   *int    `json:"page"`
	Active bool    `json:"active"`
}

func LaravelPaginate(c *gin.Context, db *gorm.DB, result interface{}) (LarvaelPaginationMeta, error) {
	var total int64

	if err := db.Model(result).Count(&total).Error; err != nil {
		return LarvaelPaginationMeta{}, err
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	if err := db.Limit(limit).Offset(offset).Find(result).Error; err != nil {
		return LarvaelPaginationMeta{}, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	path := scheme + "://" + c.Request.Host + c.FullPath()

	links := []PageLink{
		{URL: nil, Label: "&laquo; Sebelumnya", Page: nil, Active: false},
	}
	for i := 1; i <= lastPage; i++ {
		url := path + "?page=" + strconv.Itoa(i)
		pageNum := i
		links = append(links, PageLink{
			URL:    &url,
			Label:  strconv.Itoa(i),
			Page:   &pageNum,
			Active: i == page,
		})
	}
	if page < lastPage {
		next := page + 1
		url := path + "?page=" + strconv.Itoa(next)
		links = append(links, PageLink{
			URL:    &url,
			Label:  "Berikutnya &raquo;",
			Page:   &next,
			Active: false,
		})
	}

	var prevPageURL *string
	if page > 1 {
		prev := path + "?page=" + strconv.Itoa(page-1)
		prevPageURL = &prev
	}

	var nextPageURL *string
	if page < lastPage {
		next := path + "?page=" + strconv.Itoa(page+1)
		nextPageURL = &next
	}

	meta := LarvaelPaginationMeta{
		CurrentPage:  page,
		Data:         result,
		FirstPageURL: path + "?page=1",
		From:         offset + 1,
		LastPage:     lastPage,
		LastPageURL:  path + "?page=" + strconv.Itoa(lastPage),
		Links:        links,
		NextPageURL:  derefStr(nextPageURL),
		Path:         path,
		PerPage:      limit,
		PrevPageURL:  derefStr(prevPageURL),
		To:           offset + limit,
		Total:        total,
	}

	return meta, nil
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
