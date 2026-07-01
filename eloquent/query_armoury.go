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

func (q QueryArmoury) PageQuery(g *gin.Context, clazz any, resolveSubEntities bool, pageable bool, logger any) any { // TODO accepts ignore pagination, page and size
	var count int64
	cl := util.GetFieldPtrType(clazz)
	ct := util.GetFieldPtrValue(clazz).Interface()
	t1_ := reflect.MakeSlice(reflect.SliceOf(cl), 0, 0).Interface()
	request_fields := q.resolveQueryFields(g, clazz, false)
	db, page, offset, limit, sorted, paged := q.buildPageFilter(q.Db, g)
	db = q.getRawSqlQueriesFromRequestParameter(db, g)
	db = q.buildWhereFilter(db, request_fields, nil)
	if db = db.Find(&t1_).Select(util.StrFormat("%s.*", util.GetDeclaredSurefireMethod(t1_, "TableName"))); db.Error != nil { // TODO use the filters, pass which to select
		panic(db.Error)
	}
	rc := db.RowsAffected
	countDb := q.Db
	countDb = q.buildWhereFilter(countDb, request_fields, nil)
	countDb = q.getRawSqlQueriesFromRequestParameter(countDb, g)
	if err := countDb.Model(ct).Count(&count).Error; err != nil { // TODO use the filters
		panic(err)
	}
	if pageable {
		return q.PaginateResult(t1_, rc, page, offset, limit, count, sorted, paged)
	}
	return t1_
}

func (q QueryArmoury) getRawSqlQueriesFromRequestParameter(db *gorm.DB, g *gin.Context) *gorm.DB {
	queryParts := [][]string{}
	rawQueries := g.QueryArray(BARMOURY_RAW_SQL_PARAMETER_KEY)
	if len(rawQueries) > 0 {
		for _, rawQuery := range rawQueries {
			//re := regexp.MustCompile(`(where|WHERE)`)
			//queryParts = append(queryParts, re.Split(rawQuery, -1))
			queryParts = append(queryParts, strings.Split(rawQuery, "WHERE"))
		}
	}
	for _, queryPart := range queryParts {
		if queryPart[0] != "" {
			db = db.Joins(queryPart[0])
		}
		if len(queryPart) > 1 {
			db = db.Where(queryPart[1])
		}
	}
	return db
}

func (q QueryArmoury) buildWhereFilter(db *gorm.DB, request_fields map[string][]any, join_tables any) *gorm.DB {
	for matching_field_name, values := range request_fields {
		column_name := values[0]
		request_param_filter := values[2]
		db = q.getRelationQueryPart(db, column_name.(string), false, matching_field_name, request_param_filter.(map[string]any)["operator"].(string), values[3].([]string))
	}
	return db
}

func (q QueryArmoury) getRelationQueryPart(db *gorm.DB, column_name string, is_entity_field bool, matching_field_name string, operator string, values []string) *gorm.DB {
	/*matching_field_name_parts := strings.Split(matching_field_name, ".")
	object_field := matching_field_name_parts[0]
	if len(matching_field_name_parts) > 1 {
		object_field = matching_field_name_parts[1]
	}*/
	if operator == "EQ" {
		db = db.Where(util.StrFormat("%s = ?", column_name), values[0])
	} else if operator == "GT" {
		db = db.Where(util.StrFormat("%s > ?", column_name), values[0])
	} else if operator == "LT" {
		db = db.Where(util.StrFormat("%s < ?", column_name), values[0])
	} else if operator == "NE" {
		db = db.Where(util.StrFormat("%s != ?", column_name), values[0])
	} else if operator == "IN" {
		db = db.Where(util.StrFormat("%s IN ?", column_name), values)
	} else if operator == "GT_EQ" {
		db = db.Where(util.StrFormat("%s >= ?", column_name), values[0])
	} else if operator == "LT_EQ" {
		db = db.Where(util.StrFormat("%s <= ?", column_name), values[0])
	} else if operator == "LIKE" || operator == "CONTAINS" {
		db = db.Where(util.StrFormat("%s LIKE ?", column_name), util.StrFormat("%%%v%%", values[0]))
	} else if operator == "ILIKE" {
		db = db.Where(util.StrFormat("%s ILIKE ?", column_name), util.StrFormat("%%%v%%", values[0]))
	} else if operator == "NOT_LIKE" || operator == "NOT_CONTAINS" {
		db = db.Where(util.StrFormat("%s NOT LIKE ?", column_name), util.StrFormat("%%%v%%", values[0]))
	} else if operator == "NOT_ILIKE" {
		db = db.Where(util.StrFormat("%s NOT_ILIKE ?", column_name), util.StrFormat("%%%v%%", values[0]))
	} else if operator == "ENDS_WITH" {
		db = db.Where(util.StrFormat("%s LIKE ?", column_name), util.StrFormat("%%%v", values[0]))
	} else if operator == "STARTS_WITH" {
		db = db.Where(util.StrFormat("%s LIKE ?", column_name), util.StrFormat("%v%%", values[0]))
	} else if operator == "NOT_IN" {
		db = db.Where(util.StrFormat("%s NOT IN ?", column_name), values)
	} else if operator == "BETWEEN" {
		db = db.Where(util.StrFormat("%s BETWEEN ?", column_name), values)
	} else if operator == "NOT_BETWEEN" {
		db = db.Where(util.StrFormat("%s NOT BETWEEN ?", column_name), values)
	}
	return db
}

func (q QueryArmoury) resolveQueryFields(c *gin.Context, clazz any, resolve_stat_query_annotations bool) map[string][]any {
	request_fields := map[string][]any{}
	util.TranverseFields(clazz, func(f reflect.StructField) {
		rpfs_tag := f.Tag.Get("request_param_filters")
		if rpfs_tag == "" {
			return
		}
		main_field_name := f.Name
		rpfs := strings.Split(rpfs_tag, "|")
		for _, rpf := range rpfs {
			extra_field_names := []string{}
			params := strings.Split(rpf, ",")
			field_name := main_field_name
			column_name := main_field_name // get and process from the gorm tag
			request_param_filter := map[string]any{
				"multi_filter_separator": "__",
				"operator":               "NONE",
			}
			for _, param := range params {
				param_parts := strings.Split(param, "=")
				value := ""
				key := param_parts[0]
				if len(param_parts) > 1 {
					value = param_parts[1]
				}
				if key == "column" {
					column_name = value
				} else if key == "column_is_snake_case" {
					column_name = util.ToSnakeCase(column_name)
				} else if key == "boolean_to_int" {
					request_param_filter["boolean_to_int"] = true
				} else if key == "field_is_snake_case" {
					request_param_filter["field_is_snake_case"] = true
					field_name = util.ToSnakeCase(field_name)
				} else if key == "multi_filter_separator" {
					request_param_filter["multi_filter_separator"] = value
				} else if key == "operator" {
					request_param_filter["operator"] = value
					if len(rpfs) > 1 {
						extension := request_param_filter["operator"].(string)
						extra_field_names = append(extra_field_names, util.StrFormat("%v%v%s", field_name, request_param_filter["multi_filter_separator"], strings.ToLower(extension)))
						field_name = util.StrFormat("%v%v%s", field_name, request_param_filter["multi_filter_separator"], extension)
					}
				}
				// handle aliases
			}
			extra_field_names = append(extra_field_names, field_name)
			if !resolve_stat_query_annotations {

			}
			q.resolveQueryForSingleField(c, request_fields, request_param_filter, resolve_stat_query_annotations, extra_field_names, column_name, f)
		}
	})
	return request_fields
}

func (q QueryArmoury) resolveQueryForSingleField(c *gin.Context, request_fields map[string][]any, request_param_filter map[string]any,
	resolve_stat_query_annotations bool, query_params []string, column_name string, field reflect.StructField) {
	query := c.Request.URL.Query()
	for _, query_param := range query_params {
		is_present := false
		values := []string{}
		is_entity := !resolve_stat_query_annotations && request_param_filter["operator"] == "ENTITY"
		object_filter := (!resolve_stat_query_annotations &&
			(request_param_filter["operator"] == "OBJECT_EQ" ||
				request_param_filter["operator"] == "OBJECT_NE" ||
				request_param_filter["operator"] == "OBJECT_LIKE" ||
				request_param_filter["operator"] == "OBJECT_STR_EQ" ||
				request_param_filter["operator"] == "OBJECT_STR_NE" ||
				request_param_filter["operator"] == "OBJECT_NOT_LIKE" ||
				request_param_filter["operator"] == "OBJECT_CONTAINS" ||
				request_param_filter["operator"] == "OBJECT_ENDS_WITH" ||
				request_param_filter["operator"] == "OBJECT_STARTS_WITH" ||
				request_param_filter["operator"] == "OBJECT_NOT_CONTAINS" ||
				request_param_filter["operator"] == "OBJECT_STR_ENDS_WITH" ||
				request_param_filter["operator"] == "OBJECT_STR_STARTS_WITH"))
		if !resolve_stat_query_annotations {
			for key, e_values := range query {
				if key == query_param || ((object_filter || is_entity) && strings.HasPrefix(key, query_param+".")) {
					any_value_present := false
					for _, value := range e_values {
						if value == "" {
							continue
						}
						if request_param_filter["boolean_to_int"] == true {
							value = util.If(value == "true", "1", "0")
						}
						values = append(values, value)
						is_present = true
						if !any_value_present {
							any_value_present = true
						}
						if !any_value_present {
							continue
						}
					}
					query_param = util.If(object_filter && request_param_filter["column_object_fields_is_snake_case"] != "", util.ToSnakeCase(key), key)
					break
				}
			}
			_, ok := request_param_filter["always_query"]
			if !resolve_stat_query_annotations && !is_present && !ok {
				continue
			}
			if _, ok = request_fields[query_param]; ok {
				continue
			}
			if resolve_stat_query_annotations {
				query_param = field.Name
			}
			// handle join columns
			// end handle join column
			rf_value := []any{
				column_name,
				is_present,
				request_param_filter,
				values,
			}
			if resolve_stat_query_annotations {

			}
			/*if (!resolveStatQueryAnnotations && entity != null) {
			    requestFields.put(queryParam, fieldClass);
			    joinTables.put(entity.name(), joinColumn);
			}*/
			request_fields[query_param] = rf_value
		}
	}
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
