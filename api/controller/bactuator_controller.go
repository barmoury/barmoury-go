package controller

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/barmoury/barmoury-go/api/annotation"
	"github.com/barmoury/barmoury-go/api/config"
	"github.com/barmoury/barmoury-go/api/model"
	"github.com/barmoury/barmoury-go/eloquent"
	"github.com/barmoury/barmoury-go/util"
	"github.com/gin-gonic/gin"
)

var (
	SQL_QUERY_SUCCESSFUL    = "query successfully"
	SQL_QUERY_ERROR_MESSAGE = "you do not have the '%s' permission to perform this operation"
)

type BactuatorController struct {
	Self           any
	SpringLike     bool
	resourcesMap   map[string]any
	introspectMap  map[string]any
	controllersMap map[string]any
	QueryArmoury   eloquent.QueryArmoury
	config.BacuatorInterface
}

func (c *BactuatorController) DefaultAttributesAnnotations() map[string]map[string]any {
	m := make(map[string]map[string]any)

	m1 := make(map[string]any)
	m1["RequestMapping"] = annotation.RequestMapping{
		Value:  "/health",
		Method: annotation.GET,
	}
	m["HealthCheck"] = m1

	m2 := make(map[string]any)
	m2["RequestMapping"] = annotation.RequestMapping{
		Value:  "/introspect",
		Method: annotation.GET,
	}
	m["Introspect"] = m2

	m3 := make(map[string]any)
	m3["RequestMapping"] = annotation.RequestMapping{
		Method: annotation.POST,
		Value:  "/database/query/single",
	}
	m["ExecuteSingleQueries"] = m3

	m4 := make(map[string]any)
	m4["RequestMapping"] = annotation.RequestMapping{
		Method: annotation.POST,
		Value:  "/database/query/multiple",
	}
	m["ExecuteMultipleQueries"] = m4

	return m
}

func (c *BactuatorController) Setup(opts config.RouterOption) {
	c.QueryArmoury = eloquent.QueryArmoury{
		Db: opts.Db,
	}
}

func (c *BactuatorController) IsSnakeCase() bool {
	return false
}

func (c *BactuatorController) ProcessResponse(g *gin.Context, httpStatus int, apiResponseOrData any, message string) {
	if message == "" {
		g.JSON(httpStatus, apiResponseOrData)
		return
	}
	g.JSON(httpStatus, model.NewApiResponse(apiResponseOrData, message))
}

func (c *BactuatorController) processResponse(g *gin.Context, httpStatus int, apiResponseOrData any, message string) {
	util.InvokeSurefireMethod(c.Self, "ProcessResponse", reflect.ValueOf(g), reflect.ValueOf(httpStatus),
		reflect.ValueOf(apiResponseOrData), reflect.ValueOf(message))
}

func GetSchemaName(q string) string {
	rx := regexp.MustCompile(`\w+[.]\w+`)
	matches := rx.FindAllString(strings.Split(q, "FROM")[1], -1)
	if len(matches) > 0 {
		return strings.Split(matches[0], ".")[0]
	}
	return ""
}

// TODO validate and panic if the query is not for a specific schema
func (c *BactuatorController) executeQueryForResult(g *gin.Context, query string, includeColumnsName bool) any {
	//var t string
	qu := strings.ToUpper(query)
	if strings.Contains(qu, "SELECT") && !(util.InvokeSurefireMethod(c.Self, "PrincipalCan", reflect.ValueOf(g), reflect.ValueOf("SELECT"))[0].Interface().(bool)) {
		panic(util.StrFormat(SQL_QUERY_ERROR_MESSAGE, "SELECT"))
	} else if strings.Contains(qu, "UPDATE") && !(util.InvokeSurefireMethod(c.Self, "PrincipalCan", reflect.ValueOf(g), reflect.ValueOf("UPDATE"))[0].Interface().(bool)) {
		panic(util.StrFormat(SQL_QUERY_ERROR_MESSAGE, "UPDATE"))
	} else if strings.Contains(qu, "DELETE") && !(util.InvokeSurefireMethod(c.Self, "PrincipalCan", reflect.ValueOf(g), reflect.ValueOf("DELETE"))[0].Interface().(bool)) {
		panic(util.StrFormat(SQL_QUERY_ERROR_MESSAGE, "DELETE"))
	} else if strings.Contains(qu, "INSERT") && !(util.InvokeSurefireMethod(c.Self, "PrincipalCan", reflect.ValueOf(g), reflect.ValueOf("INSERT"))[0].Interface().(bool)) {
		panic(util.StrFormat(SQL_QUERY_ERROR_MESSAGE, "INSERT"))
	} else if strings.Contains(qu, "TRUNCATE") && !(util.InvokeSurefireMethod(c.Self, "PrincipalCan", reflect.ValueOf(g), reflect.ValueOf("TRUNCATE"))[0].Interface().(bool)) {
		panic(util.StrFormat(SQL_QUERY_ERROR_MESSAGE, "TRUNCATE"))
	} else if !(util.InvokeSurefireMethod(c.Self, "PrincipalCan", reflect.ValueOf(g), reflect.ValueOf("UNKNOWN"))[0].Interface().(bool)) {
		panic(util.StrFormat(SQL_QUERY_ERROR_MESSAGE, "UNKNOWN"))
	}
	var results []map[string]interface{}
	r := c.QueryArmoury.Db.Raw(query).Find(&results)
	if r.Error != nil {
		return r.Error.Error()
	}
	if !c.SpringLike || results == nil || !strings.Contains(qu, "SELECT") {
		return results
	}
	cr := [][]any{}
	for i, row := range results {
		crr := []any{}
		ccr := []any{}
		for k, v := range row {
			crr = append(crr, v)
			if i == 0 && includeColumnsName {
				ccr = append(ccr, k)
			}
		}
		if i == 0 && includeColumnsName {
			cr = append(cr, ccr)
		}
		cr = append(cr, crr)
	}
	return cr
}

func (c *BactuatorController) resolveCasing(q string) string {
	if util.InvokeSurefireMethod(c.Self, "IsSnakeCase")[0].Interface().(bool) {
		return util.ToSnakeCase(q)
	}
	return q
}

func (c *BactuatorController) resolveControllers(q string) {
	c.controllersMap = map[string]any{}
	fmt.Println("TO GET THE CONTrOLLER", q)
}

func (c *BactuatorController) resolveResources() {
	c.resourcesMap = map[string]any{}
	fmt.Println("TO GET THE resolveResources")
}

func (c *BactuatorController) resolveIntrospect() {
	c.introspectMap = map[string]any{}
	c.introspectMap["resources"] = c.resourcesMap
	c.introspectMap["controllers"] = c.controllersMap
	c.introspectMap["name"] = util.InvokeSurefireMethod(c.Self, "ServiceName")[0].Interface()
	c.introspectMap["description"] = util.InvokeSurefireMethod(c.Self, "ServiceDescription")[0].Interface()
	c.introspectMap[c.resolveCasing("logUrls")] = util.InvokeSurefireMethod(c.Self, "LogUrls")[0].Interface()
	c.introspectMap[c.resolveCasing("iconLocation")] = util.InvokeSurefireMethod(c.Self, "IconLocation")[0].Interface()
	c.introspectMap[c.resolveCasing("serviceApiName")] = util.InvokeSurefireMethod(c.Self, "ServiceApiName")[0].Interface()
	c.introspectMap[c.resolveCasing("databaseQueryRoute")] = util.InvokeSurefireMethod(c.Self, "DatabaseQueryRoute")[0].Interface()
	c.introspectMap[c.resolveCasing("databaseMultipleQueryRoute")] = util.InvokeSurefireMethod(c.Self, "DatabaseMultipleQueryRoute")[0].Interface()
}

// @RequestMapping{Value:  "/health", Method: annotation.GET}
func (c *BactuatorController) HealthCheck(g *gin.Context) {
	r := map[string]any{}
	ok := util.InvokeSurefireMethod(c.Self, "IsServiceOk")[0].Interface().(bool)
	r["status"] = util.If(ok, "ok", "not ok")
	c.processResponse(g, http.StatusOK, model.NewApiResponse(r, "health check successful"), "")
}

// @RequestMapping{Value:  "/introspect", Method: annotation.GET}
func (c *BactuatorController) Introspect(g *gin.Context) {
	if c.introspectMap == nil {
		c.resolveControllers(os.Getenv("SERVICE_BASE_URL"))
		c.resolveResources()
		c.resolveIntrospect()
	}
	c.introspectMap["users"] = util.InvokeSurefireMethod(c.Self, "UserStatistics")[0].Interface()
	c.introspectMap["earnings"] = util.InvokeSurefireMethod(c.Self, "EarningStatistics")[0].Interface()
	c.introspectMap[c.resolveCasing("downloadCounts")] = util.InvokeSurefireMethod(c.Self, "DownloadsCount")[0].Interface()
	c.processResponse(g, http.StatusOK, model.NewApiResponse(c.introspectMap, "introspect data fetched successfully"), "")
}

// @RequestMapping{Value:  "/database/query/single", Method: annotation.POST}
func (c *BactuatorController) ExecuteSingleQueries(g *gin.Context) {
	body := map[string]any{}
	if err := g.ShouldBindJSON(&body); err != nil {
		panic(errors.New("invalid request payload, " + err.Error()))
	}
	isc := util.InvokeSurefireMethod(c.Self, "IsSnakeCase")[0].Interface().(bool)
	icn := body[util.If(isc, "include_column_names", "includeColumnNames")]
	r := c.executeQueryForResult(g, body["query"].(string), icn.(bool))
	c.processResponse(g, http.StatusOK, model.NewApiResponse(r, SQL_QUERY_SUCCESSFUL), "")
}

// @RequestMapping{Value:  "/database/query/multiple", Method: annotation.POST}
func (c *BactuatorController) ExecuteMultipleQueries(g *gin.Context) {
	res := map[string]any{}
	body := map[string]any{}
	if err := g.ShouldBindJSON(&body); err != nil {
		panic(errors.New("invalid request payload, " + err.Error()))
	}
	isc := util.InvokeSurefireMethod(c.Self, "IsSnakeCase")[0].Interface().(bool)
	icn := body[util.If(isc, "include_column_names", "includeColumnNames")]
	for _, query := range body["queries"].([]interface{}) {
		res[query.(string)] = c.executeQueryForResult(g, query.(string), icn.(bool))
	}
	c.processResponse(g, http.StatusOK, model.NewApiResponse(res, SQL_QUERY_SUCCESSFUL), "")
}
