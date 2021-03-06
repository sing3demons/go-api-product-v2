package controllers

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type pagingResult struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	PrevPage  int `json:"prevPage"`
	NextPage  int `json:"nextPage"`
	Count     int `json:"count"`
	TotalPage int `json:"totalPage"`
}

type pagination struct {
	ctx     *gin.Context
	query   *gorm.DB
	records interface{}
}

func (p *pagination) pagingResource() *pagingResult {
	page, _ := strconv.Atoi(p.ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(p.ctx.DefaultQuery("limit", "12"))

	ch := make(chan int)
	go p.countRecords(ch)

	offset := (page - 1) * limit
	p.query.Offset(offset).Limit(limit).Find(p.records)

	count := <-ch
	totalPage := int(math.Ceil(float64(count) / float64(limit)))
	// 5. Find nextPage
	var nextPage int
	if nextPage == totalPage {
		nextPage = totalPage
	} else {
		nextPage = totalPage + 1
	}
	// 6. create pagingResult
	return &pagingResult{
		Page:      page,
		Limit:     limit,
		PrevPage:  page - 1,
		NextPage:  nextPage,
		Count:     count,
		TotalPage: totalPage,
	}

}

func (p *pagination) countRecords(ch chan int) {
	var count int64
	p.query.Model(p.records).Count(&count)

	ch <- int(count)
}
