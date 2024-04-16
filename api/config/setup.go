package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/barmoury/barmoury-go/api/annotation"
	"github.com/barmoury/barmoury-go/meta"
	"github.com/barmoury/barmoury-go/util"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var registeredRoutes bool
var BarmouryGormDb *gorm.DB
var deferedhandlers map[string][]gin.HandlerFunc

type BacuatorInterface interface {
	Resources() []any
	IsServiceOk() bool
	Controllers() []any
	ServiceName() string
	IconLocation() string
	ServiceApiName() string
	DownloadsCount() uint64
	ServiceDescription() string
	DatabaseQueryRoute() string
	LogUrls() map[string][]string
	UserStatistics() map[string]uint64
	DatabaseMultipleQueryRoute() string
	EarningStatistics() map[string]uint64
	PrincipalCan(g *gin.Context, dbMethod string) bool
	meta.IAnotation
}

type RouterOption struct {
	Prefix      string
	Db          *gorm.DB
	RouterGroup *gin.RouterGroup
	Bacuator    BacuatorInterface
}

func RegisterControllers(engine *gin.Engine, opts RouterOption, controllers []meta.IAnotation) RouterOption {
	registeredRoutes = true
	UseDeferedHandlers(engine)
	if opts.RouterGroup == nil {
		opts.RouterGroup = engine.Group(util.GetOrDefault(opts.Prefix, ""))
	}
	if opts.Db == nil {
		opts.Db = createDefaultGormDbConnection()
		BarmouryGormDb = opts.Db
	}
	if opts.Bacuator != nil {
		controllers = append(controllers, opts.Bacuator)
	}
	for _, controller := range controllers {
		registerController(controller, opts)
	}
	return opts
}

func registerController(controller meta.IAnotation, opts RouterOption) {
	var e annotation.RequestMapping
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
		registerRoutes(controller, e, opts)
	}()
	ds, _ := controller.Annotations()
	e = meta.GetAnnotation[annotation.RequestMapping](ds, "RequestMapping")
}

func registerRoutes(controller interface{}, requestMapping annotation.RequestMapping, opts RouterOption) {
	var routesSet []string
	methodsAnnotations := meta.GetAnnotationsFromAttributesAnnotationsMethods(controller)
	controllerRoute := util.GetOrDefault(requestMapping.Value, "")
	if util.GetDeclaredSurefireFieldFromPtr(controller, "Controller").CanAddr() {
		if util.GetDeclaredSurefireFieldFromPtr(controller, "Controller").IsNil() {
			// TODO auto set a value for the Controller
			panic(errors.New("the controller did not initialize it Controller struct"))
		}
		util.SetFieldValue(util.GetDeclaredSurefireFieldFromPtr(controller, "Controller").Interface(), "Self", controller)
	} else {
		if util.GetDeclaredSurefireFieldFromPtr(controller, "BactuatorController").IsNil() {
			// TODO auto set a value for the BactuatorController
			panic(errors.New("the controller did not initialize it BactuatorController struct"))
		}
		util.SetFieldValue(util.GetDeclaredSurefireFieldFromPtr(controller, "BactuatorController").Interface(), "Self", controller)
	}
	util.TranverseDeclaredMethods(controller, func(m reflect.Method, v reflect.Value) {
		if m.Name == "Setup" {
			v.Call([]reflect.Value{reflect.ValueOf(opts)})
			return
		}
		rm, ok := meta.GetAttributesAnnotations[annotation.RequestMapping](methodsAnnotations, m.Name, "RequestMapping")
		if !ok {
			return
		}

		method := util.GetOrDefault(rm.Method, annotation.GET)
		route := controllerRoute + util.GetOrDefault(rm.Value, "")
		routerPath := fmt.Sprintf("%s__%s%s", method, opts.Prefix, route)
		if util.ValueInSlice(routesSet, routerPath) {
			return
		}
		routesSet = append(routesSet, routerPath)
		switch method {
		case annotation.PUT:
			opts.RouterGroup.PUT(route, v.Interface().(func(*gin.Context)))
		case annotation.HEAD:
			opts.RouterGroup.HEAD(route, v.Interface().(func(*gin.Context)))
		case annotation.POST:
			opts.RouterGroup.POST(route, v.Interface().(func(*gin.Context)))
		case annotation.PATCH:
			opts.RouterGroup.PATCH(route, v.Interface().(func(*gin.Context)))
		case annotation.TRACE:
			opts.RouterGroup.OPTIONS(route, v.Interface().(func(*gin.Context)))
		case annotation.DELETE:
			opts.RouterGroup.DELETE(route, v.Interface().(func(*gin.Context)))
		case annotation.OPTIONS:
			opts.RouterGroup.OPTIONS(route, v.Interface().(func(*gin.Context)))
		default:
			opts.RouterGroup.GET(route, v.Interface().(func(*gin.Context)))
		}
	})
}

func createDefaultGormDbConnection() *gorm.DB {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_SCHEMA"),
		os.Getenv("DATABASE_QUERY_STRING"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("barmoury - Error connecting to database : error=%s", err))
	}

	return db
}

func UseDeferedHandlers(engine *gin.Engine) {
	if len(deferedhandlers) == 0 {
		return
	}
	jms := deferedhandlers["JWT_MANAGERS"]
	eas := deferedhandlers["ERROR_ADVISERS"]
	rvs := deferedhandlers["ROUTE_VALIDATORS"]
	for _, rv := range rvs {
		engine.Use(rv)
	}
	for _, ea := range eas {
		engine.Use(ea)
	}
	for _, jm := range jms {
		engine.Use(jm)
	}
	delete(deferedhandlers, "JWT_MANAGERS")
	delete(deferedhandlers, "ROUTE_VALIDATORS")
	delete(deferedhandlers, "ERROR_ADVISERS")
	deferedhandlers = nil
}

func deferedHandler(group string, handler gin.HandlerFunc) {
	var dh []gin.HandlerFunc
	if deferedhandlers == nil {
		deferedhandlers = map[string][]gin.HandlerFunc{}
	}
	if dh_, ok := deferedhandlers[group]; ok {
		dh = dh_
	}
	dh = append(dh, handler)
	deferedhandlers[group] = dh
}
