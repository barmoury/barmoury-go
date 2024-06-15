package controller

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/barmoury/barmoury-go/api/annotation"
	"github.com/barmoury/barmoury-go/api/config"
	"github.com/barmoury/barmoury-go/api/model"
	"github.com/barmoury/barmoury-go/copier"
	"github.com/barmoury/barmoury-go/eloquent"
	"github.com/barmoury/barmoury-go/meta"
	"github.com/barmoury/barmoury-go/util"
	"github.com/gin-gonic/gin"
)

type RouteMethod string

const (
	STAT             RouteMethod = "STAT"
	SHOW             RouteMethod = "SHOW"
	INDEX            RouteMethod = "INDEX"
	STORE            RouteMethod = "STORE"
	UPDATE           RouteMethod = "UPDATE"
	DESTROY          RouteMethod = "DESTROY"
	STORE_MULTIPLE   RouteMethod = "STORE_MULTIPLE"
	DESTROY_MULTIPLE RouteMethod = "DESTROY_MULTIPLE"
)

var (
	NO_RESOURCE_FORMAT_STRING = "no %s found with the specified id %s"
	ACCESS_DENIED             = "access denied, you do not have the required role to access this endpoint"
)

type Controller[T1 any, T2 any] struct {
	Self                 any
	Pageble              bool
	StoreAsynchronously  bool
	UpdateAsynchronously bool
	DeleteAsynchronously bool
	FineName             string
	Router               *gin.RouterGroup
	QueryArmoury         eloquent.QueryArmoury
	meta.IAnotation
}

func (c *Controller[T1, T2]) DefaultAttributesAnnotations() map[string]map[string]any {
	m := make(map[string]map[string]any)

	sm := make(map[string]any)
	sm["RequestMapping"] = annotation.RequestMapping{
		Value:  "/stat",
		Method: annotation.GET,
	}
	m["Stat"] = sm

	im := make(map[string]any)
	im["RequestMapping"] = annotation.RequestMapping{
		Method: annotation.GET,
	}
	m["Index"] = im

	//@Validate()
	stm := make(map[string]any)
	stm["RequestMapping"] = annotation.RequestMapping{
		Method: annotation.POST,
	}
	m["Store"] = stm

	stmm := make(map[string]any)
	stmm["RequestMapping"] = annotation.RequestMapping{
		Value:  "/multiple",
		Method: annotation.POST,
	}
	m["StoreMultiple"] = stmm

	gm := make(map[string]any)
	gm["RequestMapping"] = annotation.RequestMapping{
		Value:  "/:id",
		Method: annotation.GET,
	}
	m["Show"] = gm

	//@Validate({ groups: ["UPDATE"] })
	up := make(map[string]any)
	up["RequestMapping"] = annotation.RequestMapping{
		Value:  "/:id",
		Method: annotation.PATCH,
	}
	m["Update"] = up

	dm := make(map[string]any)
	dm["RequestMapping"] = annotation.RequestMapping{
		Value:  "/:id",
		Method: annotation.DELETE,
	}
	m["Destroy"] = dm

	dmm := make(map[string]any)
	dmm["RequestMapping"] = annotation.RequestMapping{
		Value:  "/multiple",
		Method: annotation.DELETE,
	}
	m["DestroyMultiple"] = dmm

	return m
}

func (c *Controller[T1, T2]) Setup(opts config.RouterOption) {
	var t1 T1
	c.Router = opts.RouterGroup
	c.QueryArmoury = eloquent.QueryArmoury{
		Db: opts.Db,
	}
	c.FineName = util.SplitByRegex(util.GetTypeName(t1), "\\[")[0]
}

func (c *Controller[T1, T2]) PreResponse(entity *T1) {
}

func (c *Controller[T1, T2]) PreResponses(entities []*T1) {
	for _, entity := range entities {
		util.InvokeSurefireMethod(c.Self, "PreResponse", reflect.ValueOf(entity))
	}
}

func (c *Controller[T1, T2]) ResolveSubEntities() bool {
	return true
}

func (c *Controller[T1, T2]) SkipRecursiveSubEntities() bool {
	return true
}

func (c *Controller[T1, T2]) PreQuery(g *gin.Context, authentication any) *gin.Context {
	return g
}

func (c *Controller[T1, T2]) PreCreate(g *gin.Context, authentication any, entity *T1, entityRequest *T2) {
}

func (c *Controller[T1, T2]) PostCreate(g *gin.Context, authentication any, entity *T1) {
}

func (c *Controller[T1, T2]) PreUpdate(g *gin.Context, authentication any, entity *T1, entityRequest *T2) {
}

func (c *Controller[T1, T2]) PostUpdate(g *gin.Context, authentication any, prevEntity *T1, entity *T1) {
}

func (c *Controller[T1, T2]) PreDelete(g *gin.Context, authentication any, entity *T1, id any) {
}

func (c *Controller[T1, T2]) PostDelete(g *gin.Context, authentication any, entity *T1) {
}

func (c *Controller[T1, T2]) OnAsynchronousError(t string, entity any, err error) {
}

func (c *Controller[T1, T2]) HandleSqlInjectionQuery(g *gin.Context, authentication any) {
	panic(errors.New("sql injection attack detected"))
}

func (c *Controller[T1, T2]) SanitizeAndGetRequestParameters(g *gin.Context, authentication any) *gin.Context {
	if g.Request.URL.Query().Has(eloquent.BARMOURY_RAW_SQL_PARAMETER_KEY) {
		util.InvokeSurefireMethod(c.Self, "HandleSqlInjectionQuery", reflect.ValueOf(g), reflect.ValueOf(authentication))
	}
	return g
}

func (c *Controller[T1, T2]) ProcessResponse(g *gin.Context, httpStatus int, apiResponseOrData any, message string) {
	if message == "" {
		g.JSON(httpStatus, apiResponseOrData)
		return
	}
	g.JSON(httpStatus, model.NewApiResponse(apiResponseOrData, message))
}

func (c *Controller[T1, T2]) GetResourceById(id any, authentication any) *T1 {
	var t T1
	return c.QueryArmoury.GetResourceById(t, id, util.StrFormat(NO_RESOURCE_FORMAT_STRING, c.FineName, id)).(*T1)
}

func (c *Controller[T1, T2]) PostGetResourceById(g *gin.Context, authentication any, entity *T1) {

}

func (c *Controller[T1, T2]) ValidateBeforeCommit(t1 *T1) string {
	return ""
}

func (c *Controller[T1, T2]) ShouldNotHonourMethod(routeMethod RouteMethod) bool {
	return false
}

func (c *Controller[T1, T2]) GetRouteMethodRoles(routeMethod RouteMethod) []string {
	return []string{}
}

func (c *Controller[T1, T2]) ValidateRouteAccess(g *gin.Context, routeMethod RouteMethod, errMessage string) {
	if util.InvokeSurefireMethod(c.Self, "ShouldNotHonourMethod", reflect.ValueOf(routeMethod))[0].Bool() {
		panic(errors.New(errMessage))
	}
	roles := util.InvokeSurefireMethod(c.Self, "GetRouteMethodRoles", reflect.ValueOf(routeMethod))[0].Interface().([]string)
	av, ok := g.Get("authoritiesValues")
	if len(roles) > 0 && (!ok || !util.SlicesIntercepts(roles, av.([]string))) {
		panic(errors.New(ACCESS_DENIED))
	}
}

func (c *Controller[T1, T2]) InjectUpdateFieldId(g *gin.Context, resourceRequest T2) T2 {
	if !(g.Request.Method == "POST" || g.Request.Method == "PUT" || g.Request.Method == "PATCH") {
		return resourceRequest
	}
	util.SetFieldValue(resourceRequest, "___BARMOURY_UPDATE_ENTITY_ID___", g.Param("userid"))
	return resourceRequest
}

func (c *Controller[T1, T2]) ResolveRequestPayload(authentication any, resourceRequest *T2) *T1 {
	var t1_ T1
	t1 := &t1_
	rm := util.GetDeclaredMethodValue(&t1, "Resolve")
	if rm.IsValid() {
		t1 = rm.Call([]reflect.Value{reflect.ValueOf(resourceRequest), reflect.ValueOf(c.QueryArmoury), reflect.ValueOf(authentication)})[0].Interface().(*T1)
	}
	return t1
}

func (c *Controller[T1, T2]) GetAuthentication(g *gin.Context) any {
	if a, ok := g.Get("user"); ok {
		return a
	}
	return 0
}

func (c *Controller[T1, T2]) processResponse(g *gin.Context, httpStatus int, apiResponseOrData any, message string) {
	util.InvokeSurefireMethod(c.Self, "ProcessResponse", reflect.ValueOf(g), reflect.ValueOf(httpStatus),
		reflect.ValueOf(apiResponseOrData), reflect.ValueOf(message))
}

// @RequestMapping{Value:  "/stat", Method: annotation.GET}
func (c *Controller[T1, T2]) Stat(g *gin.Context) {
	authentication := c.GetAuthentication(g)
	c.ValidateRouteAccess(g, STAT, "the GET '**/stat' route is not supported for this resource")
	st := util.InvokeSurefireMethod(c.Self, "SanitizeAndGetRequestParameters", reflect.ValueOf(g), reflect.ValueOf(authentication))[0].Interface().(*gin.Context)
	g = util.InvokeSurefireMethod(c.Self, "PreQuery", reflect.ValueOf(st), reflect.ValueOf(authentication))[0].Interface().(*gin.Context)
	c.processResponse(g, http.StatusOK, model.NewApiResponse([]T1{}, util.StrFormat("%s statistics fetched successfully", c.FineName)), "")
}

// @RequestMapping{Method: annotation.GET}
func (c *Controller[T1, T2]) Index(g *gin.Context) {
	var t1 *T1
	authentication := c.GetAuthentication(g)
	c.ValidateRouteAccess(g, INDEX, "the GET '**/' route is not supported for this resource")
	st := util.InvokeSurefireMethod(c.Self, "SanitizeAndGetRequestParameters", reflect.ValueOf(g), reflect.ValueOf(authentication))[0].Interface().(*gin.Context)
	g = util.InvokeSurefireMethod(c.Self, "PreQuery", reflect.ValueOf(st), reflect.ValueOf(authentication))[0].Interface().(*gin.Context)
	pageble := util.GetDeclaredFieldValueAs[bool](c, "Pageble")
	rses := util.InvokeSurefireMethod(c.Self, "ResolveSubEntities")[0].Interface().(bool)
	resources := c.QueryArmoury.PageQuery(g, t1, rses, pageble, nil)
	if pageble {
		util.InvokeSurefireMethod(c.Self, "PreResponses", reflect.ValueOf((resources.(map[string]any)["content"]).([]*T1)))
	} else {
		util.InvokeSurefireMethod(c.Self, "PreResponses", reflect.ValueOf(resources.([]*T1)))
	}
	c.processResponse(g, http.StatusOK, model.NewApiResponse(resources, util.StrFormat("%s list fetched successfully", c.FineName)), "")
}

// @Validated()
// @RequestMapping{Method: annotation.POST}
func (c *Controller[T1, T2]) Store(g *gin.Context) {
	authentication := c.GetAuthentication(g)
	c.ValidateRouteAccess(g, STORE, "the POST '**/' route is not supported for this resource")
	var request *T2
	if err := g.ShouldBindJSON(&request); err != nil {
		panic(errors.New("invalid request payload, " + err.Error()))
	}
	async := util.GetDeclaredFieldValueAs[bool](c, "StoreAsynchronously")
	resource := util.InvokeSurefireMethod(c.Self, "ResolveRequestPayload", reflect.ValueOf(authentication), reflect.ValueOf(request))[0].Interface().(*T1)
	util.InvokeSurefireMethod(c.Self, "PreCreate", reflect.ValueOf(g), reflect.ValueOf(authentication), reflect.ValueOf(resource), reflect.ValueOf(request))
	msg := util.InvokeSurefireMethod(c.Self, "ValidateBeforeCommit", reflect.ValueOf(resource))[0].Interface().(string)
	if msg != "" {
		panic(errors.New(msg))
	}
	if async {
		go func() {
			if err := c.QueryArmoury.Db.Create(resource).Error; err != nil {
				util.InvokeSurefireMethod(c.Self, "OnAsynchronousError", reflect.ValueOf("Store"), reflect.ValueOf(resource), reflect.ValueOf(err))
				return
			}
			util.InvokeSurefireMethod(c.Self, "PostCreate", reflect.ValueOf(g), reflect.ValueOf(authentication), reflect.ValueOf(resource))
		}()
		c.processResponse(g, http.StatusAccepted, model.NewApiResponse[any](nil, util.StrFormat("%s is being created", c.FineName)), "")
		return
	}
	if err := c.QueryArmoury.Db.Create(resource).Error; err != nil {
		panic(err)
	}
	util.InvokeSurefireMethod(c.Self, "PostCreate", reflect.ValueOf(g), reflect.ValueOf(authentication), reflect.ValueOf(resource))
	util.InvokeSurefireMethod(c.Self, "PreResponse", reflect.ValueOf(resource))
	c.processResponse(g, http.StatusCreated, model.NewApiResponse(resource, util.StrFormat("%s created successfully", c.FineName)), "")
}

// @RequestMapping{Value:  "/multiple", Method: annotation.POST}
func (c *Controller[T1, T2]) StoreMultiple(g *gin.Context) {
	authentication := c.GetAuthentication(g)
	c.ValidateRouteAccess(g, STORE_MULTIPLE, "the POST '**/multiple' route is not supported for this resource")
	var requests []*T2
	var resources []*T1
	if err := g.ShouldBindJSON(&requests); err != nil {
		panic(errors.New("invalid request payload, " + err.Error()))
	}
	async := util.GetDeclaredFieldValueAs[bool](c, "StoreAsynchronously")
	for _, request := range requests {
		resource := util.InvokeSurefireMethod(c.Self, "ResolveRequestPayload", reflect.ValueOf(authentication), reflect.ValueOf(request))[0].Interface().(*T1)
		util.InvokeSurefireMethod(c.Self, "PreCreate", reflect.ValueOf(g), reflect.ValueOf(authentication), reflect.ValueOf(resource), reflect.ValueOf(request))
		msg := util.InvokeSurefireMethod(c.Self, "ValidateBeforeCommit", reflect.ValueOf(resource))[0].Interface().(string)
		if msg != "" {
			panic(errors.New(msg))
		}
		if async {
			go func() {
				if err := c.QueryArmoury.Db.Create(resource).Error; err != nil {
					util.InvokeSurefireMethod(c.Self, "OnAsynchronousError", reflect.ValueOf("Store"), reflect.ValueOf(resource), reflect.ValueOf(err))
					return
				}
				util.InvokeSurefireMethod(c.Self, "PostCreate", reflect.ValueOf(g), reflect.ValueOf(authentication), reflect.ValueOf(resource))
			}()
			continue
		}
		if err := c.QueryArmoury.Db.Create(resource).Error; err != nil {
			panic(err)
		}
		util.InvokeSurefireMethod(c.Self, "PostCreate", reflect.ValueOf(g), reflect.ValueOf(authentication), reflect.ValueOf(resource))
		util.InvokeSurefireMethod(c.Self, "PreResponse", reflect.ValueOf(resource))
		resources = append(resources, resource)
	}
	if async {
		c.processResponse(g, http.StatusAccepted, model.NewApiResponse[any](nil, util.StrFormat("%ss are being created", c.FineName)), "")
		return
	}
	c.processResponse(g, http.StatusAccepted, model.NewApiResponse[any](resources, util.StrFormat("%ss created successfully", c.FineName)), "")
}

// @RequestMapping{Value:  "/:id", Method: annotation.GET}
func (c *Controller[T1, T2]) Show(g *gin.Context) {
	c.ValidateRouteAccess(g, SHOW, "the GET '**/:id' route is not supported for this resource")
	id := g.Param("id")
	authentication := c.GetAuthentication(g)
	resource := util.InvokeSurefireMethod(c.Self, "GetResourceById", reflect.ValueOf(id), reflect.ValueOf(authentication))[0].Interface().(*T1)
	util.InvokeSurefireMethod(c.Self, "PostGetResourceById", reflect.ValueOf(g), reflect.ValueOf(authentication), reflect.ValueOf(resource))
	util.InvokeSurefireMethod(c.Self, "PreResponse", reflect.ValueOf(resource))
	c.processResponse(g, http.StatusOK, model.NewApiResponse(resource, util.StrFormat("%s fetch successfully", c.FineName)), "")
}

// @Validated(Groups: []string{"UPDATE"})
// @RequestMapping{Value:  "/:id", Method: annotation.PATCH}
func (c *Controller[T1, T2]) Update(g *gin.Context) {
	c.ValidateRouteAccess(g, UPDATE, "the PATCH '**/:id' route is not supported for this resource")
	id := g.Param("id")
	var request *T2
	if err := g.BindJSON(&request); err != nil {
		panic(errors.New("invalid request payload, " + err.Error()))
	}
	var t1 T1
	authentication := c.GetAuthentication(g)
	prevResource := reflect.New(reflect.TypeOf(t1))
	async := util.GetDeclaredFieldValueAs[bool](c, "UpdateAsynchronously")
	resource := util.InvokeSurefireMethod(c.Self, "GetResourceById", reflect.ValueOf(id), reflect.ValueOf(authentication))[0].Interface().(*T1)
	copier.Copy(&prevResource, resource)
	util.InvokeSurefireMethod(c.Self, "PostGetResourceById", reflect.ValueOf(g), reflect.ValueOf(authentication), reflect.ValueOf(resource))
	rm := util.GetDeclaredMethodValue(&resource, "Resolve")
	if rm.IsValid() {
		resource = rm.Call([]reflect.Value{reflect.ValueOf(request), reflect.ValueOf(c.QueryArmoury), reflect.ValueOf(authentication)})[0].Interface().(*T1)
	}
	util.InvokeSurefireMethod(c.Self, "PreUpdate", reflect.ValueOf(g), reflect.ValueOf(authentication), reflect.ValueOf(resource), reflect.ValueOf(request))
	msg := util.InvokeSurefireMethod(c.Self, "ValidateBeforeCommit", reflect.ValueOf(resource))[0].Interface().(string)
	if msg != "" {
		panic(errors.New(msg))
	}
	if async {
		go func() {
			if err := c.QueryArmoury.Db.Save(resource).Error; err != nil {
				util.InvokeSurefireMethod(c.Self, "OnAsynchronousError", reflect.ValueOf("Update"), reflect.ValueOf(resource), reflect.ValueOf(err))
				return
			}
			util.InvokeSurefireMethod(c.Self, "PostUpdate", reflect.ValueOf(g), reflect.ValueOf(authentication), reflect.ValueOf(resource), reflect.ValueOf(resource))
		}()
		c.processResponse(g, http.StatusAccepted, model.NewApiResponse[any](nil, util.StrFormat("%s is being created", c.FineName)), "")
		return
	}
	if err := c.QueryArmoury.Db.Save(resource).Error; err != nil {
		panic(err)
	}
	util.InvokeSurefireMethod(c.Self, "PostUpdate", reflect.ValueOf(g), reflect.ValueOf(authentication), reflect.ValueOf(resource), reflect.ValueOf(resource))
	util.InvokeSurefireMethod(c.Self, "PreResponse", reflect.ValueOf(resource))
	c.processResponse(g, http.StatusCreated, model.NewApiResponse(resource, util.StrFormat("%s updated successfully", c.FineName)), "")
}

// @RequestMapping{Value:  "/:id", Method: annotation.DELETE}
func (c *Controller[T1, T2]) Destroy(g *gin.Context) {
	c.ValidateRouteAccess(g, DESTROY, "the DELETE '**/:id' route is not supported for this resource")
	id := g.Param("id")
	authentication := c.GetAuthentication(g)
	async := util.GetDeclaredFieldValueAs[bool](c, "DeleteAsynchronously")
	resource := util.InvokeSurefireMethod(c.Self, "GetResourceById", reflect.ValueOf(id), reflect.ValueOf(authentication))[0].Interface().(*T1)
	util.InvokeSurefireMethod(c.Self, "PostGetResourceById", reflect.ValueOf(g), reflect.ValueOf(authentication), reflect.ValueOf(resource))
	util.InvokeSurefireMethod(c.Self, "PreDelete", reflect.ValueOf(g), reflect.ValueOf(authentication), reflect.ValueOf(resource), reflect.ValueOf(id))
	if async {
		go func() {
			if err := c.QueryArmoury.Db.Delete(resource).Error; err != nil {
				util.InvokeSurefireMethod(c.Self, "OnAsynchronousError", reflect.ValueOf("Delete"), reflect.ValueOf(resource), reflect.ValueOf(err))
				return
			}
			util.InvokeSurefireMethod(c.Self, "PostDelete", reflect.ValueOf(g), reflect.ValueOf(authentication), reflect.ValueOf(resource))
		}()
		c.processResponse(g, http.StatusAccepted, model.NewApiResponse[any](nil, util.StrFormat("%s is being deleted", c.FineName)), "")
		return
	}
	if err := c.QueryArmoury.Db.Delete(resource).Error; err != nil {
		panic(err)
	}
	util.InvokeSurefireMethod(c.Self, "PostDelete", reflect.ValueOf(g), reflect.ValueOf(authentication), reflect.ValueOf(resource))
	util.InvokeSurefireMethod(c.Self, "PreResponse", reflect.ValueOf(resource))
	g.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	g.Writer.WriteHeader(http.StatusNoContent)
}

// @RequestMapping{Value:  "/multiple", Method: annotation.DELETE}
func (c *Controller[T1, T2]) DestroyMultiple(g *gin.Context) {
	c.ValidateRouteAccess(g, DESTROY_MULTIPLE, "the DELETE '**/multiple' route is not supported for this resource")
	var ids []any
	var resources []*T1
	if err := g.ShouldBindJSON(&ids); err != nil {
		panic(errors.New("invalid request payload for multiple deletion, " + err.Error()))
	}
	authentication := c.GetAuthentication(g)
	async := util.GetDeclaredFieldValueAs[bool](c, "DeleteAsynchronously")
	for _, id := range ids {
		resources = append(resources,
			util.InvokeSurefireMethod(c.Self, "GetResourceById", reflect.ValueOf(id), reflect.ValueOf(authentication))[0].Interface().(*T1))
	}
	for i, resource := range resources {
		util.InvokeSurefireMethod(c.Self, "PostGetResourceById", reflect.ValueOf(g), reflect.ValueOf(authentication), reflect.ValueOf(resource))
		util.InvokeSurefireMethod(c.Self, "PreDelete", reflect.ValueOf(g), reflect.ValueOf(authentication), reflect.ValueOf(resource), reflect.ValueOf(ids[i]))
		if async {
			go func() {
				if err := c.QueryArmoury.Db.Delete(resource).Error; err != nil {
					util.InvokeSurefireMethod(c.Self, "OnAsynchronousError", reflect.ValueOf("Delete"), reflect.ValueOf(resource), reflect.ValueOf(err))
					return
				}
				util.InvokeSurefireMethod(c.Self, "PostDelete", reflect.ValueOf(g), reflect.ValueOf(authentication), reflect.ValueOf(resource))
			}()
			c.processResponse(g, http.StatusAccepted, model.NewApiResponse[any](nil, util.StrFormat("%s is being deleted", c.FineName)), "")
			return
		}
		if err := c.QueryArmoury.Db.Delete(resource).Error; err != nil {
			panic(err)
		}
		util.InvokeSurefireMethod(c.Self, "PostDelete", reflect.ValueOf(g), reflect.ValueOf(authentication), reflect.ValueOf(resource))
		util.InvokeSurefireMethod(c.Self, "PreResponse", reflect.ValueOf(resource))
	}
	g.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	if async {
		c.processResponse(g, http.StatusAccepted, model.NewApiResponse[any](nil, util.StrFormat("%ss are being deleted", c.FineName)), "")
		return
	}
	g.Writer.WriteHeader(http.StatusNoContent)
}
