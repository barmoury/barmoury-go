package eloquent

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/barmoury/barmoury-go/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	BARMOURY_RAW_SQL_PARAMETER_KEY = "___BARMOURY__RAW__SQL___"
)

type QueryArmoury struct {
	Db *gorm.DB
}

func (q QueryArmoury) PageQuery(g *gin.Context, clazz any, resolveSubEntities bool, pageable bool, logger any) any {
	var count int64
	cl := util.GetFieldPtrType(clazz)
	ct := util.GetFieldPtrValue(clazz).Interface()
	t1_ := reflect.MakeSlice(reflect.SliceOf(cl), 0, 0).Interface()
	db, page, offset, limit, sorted, paged := q.buildPageFilter(q.Db, g)
	if db = db.Find(&t1_); db.Error != nil { // TODO use the filters
		panic(db.Error)
	}
	rc := db.RowsAffected
	if err := q.Db.Model(ct).Count(&count).Error; err != nil { // TODO use the filters
		panic(err)
	}
	if pageable {
		return q.PaginateResult(t1_, rc, page, offset, limit, count, sorted, paged)
	}
	return t1_
}

func (q QueryArmoury) buildPageFilter(db *gorm.DB, c *gin.Context) (*gorm.DB, int, int64, int64, bool, bool) {
	page := 1
	limit := 10
	sorted := false
	paged := false
	query := c.Request.URL.Query()
	if size, err := strconv.Atoi(query.Get("size")); err == nil && size > 0 {
		limit = size
	}
	if page_, err := strconv.Atoi(query.Get("page")); err == nil && page_ > 0 {
		page = page_
		paged = true
	}
	offset := ((page - 1) * limit)
	db = db.Offset(offset).Limit(limit)
	if sorts_ := c.QueryArray("sort"); len(sorts_) > 0 {
		sorted = true
		for _, sorts := range sorts_ {
			sort_ := strings.Split(sorts, ",")
			sort := "`" + sort_[0] + "`"
			if len(sort_) > 0 {
				sort += " " + sort_[1]
			}
			db = db.Order(sort)
		}
	}
	return db, page, int64(offset), int64(limit), sorted, paged
}

func (q QueryArmoury) PaginateResult(rows any, rowsCount int64, page int, offset int64, limit int64, count int64, sorted bool, paged bool) map[string]any {
	sort := map[string]any{
		"empty":    rowsCount == 0,
		"sorted":   sorted,
		"unsorted": !sorted,
	}
	return map[string]any{
		"content": rows,
		"pageable": map[string]any{
			"sort":        sort,
			"offset":      offset,
			"page_number": page,
			"page_size":   limit,
			"paged":       paged,
			"unpaged":     !paged,
		},
		"last":               offset >= (count - limit),
		"total_pages":        (count / limit) + 1,
		"total_elements":     count,
		"first":              offset == 0,
		"size":               limit,
		"number":             page,
		"sort":               sort,
		"number_of_elements": rowsCount,
		"empty":              rowsCount == 0,
	}
}

func (q QueryArmoury) GetResourceById(t any, id any, message string) any {
	tt := util.GetFieldNonPtrType(t)
	v := reflect.New(tt).Interface()
	if q.Db.First(&v, id).Error != nil {
		panic(errors.New(message))
	}
	return v
}
