package config

import (
	"errors"
	"reflect"
	"strings"

	"github.com/barmoury/barmoury-go/api/annotation"
	"github.com/barmoury/barmoury-go/api/model"
	"github.com/barmoury/barmoury-go/meta"
	"github.com/barmoury/barmoury-go/util"
	"github.com/gin-gonic/gin"
)

var ActveGinContext *gin.Context

type ErrorAdviserOption struct {
	Log func(...any)
}

type ErrorAdviser struct {
}

type errorAdviserMapEntry struct {
	Instance   any
	StatusCode int
	MethodName string
}

var errorAdviserMap map[string]errorAdviserMapEntry

func RegisterErrorAdvisers(engine *gin.Engine, opts ErrorAdviserOption, advisers []any) {
	if errorAdviserMap == nil {
		errorAdviserMap = map[string]errorAdviserMapEntry{}
	}
	for _, adviser := range advisers {
		annos := meta.GetAnnotationsFromAttributesAnnotationsMethods(adviser)
		util.TranverseDeclaredMethods(adviser, func(m reflect.Method, v reflect.Value) {
			ea, ok := meta.GetAttributesAnnotations[annotation.ErrorAdvise](annos, m.Name, "ErrorAdvise")
			if !ok {
				return
			}
			for _, err := range ea.Errors {
				errorAdviserMap[err] = errorAdviserMapEntry{
					MethodName: m.Name,
					Instance:   adviser,
					StatusCode: ea.StatusCode,
				}
			}
		})
	}
	handler := func(c *gin.Context) {
		ActveGinContext = c
		defer func() {
			errs := c.Errors.Errors()
			if r := recover(); r != nil {
				errs = append(errs, util.StrFormat("%s", r))
			}
			if len(errs) > 0 {
				var ok bool
				var cb errorAdviserMapEntry
				firstError := errs[0]
				if cb, ok = errorAdviserMap[firstError]; !ok {
					for err, v := range errorAdviserMap {
						if strings.Contains(firstError, err) {
							cb = v
							ok = true
							break
						}
					}
				}
				if !ok {
					if cb, ok = errorAdviserMap["___UnknownError___"]; !ok {
						panic(firstError)
					}
				}
				res := util.InvokeSurefireMethod(cb.Instance, cb.MethodName, reflect.ValueOf(errors.New(firstError)),
					reflect.ValueOf(opts))[0]
				c.JSON(cb.StatusCode, res.Interface().(model.ApiResponse[any]))
				c.Abort()
				return
			}
		}()
		c.Next()
	}
	engine.NoRoute(func(c *gin.Context) {
		c.Error(errors.New("404 page not found"))
	})
	engine.NoMethod(func(c *gin.Context) {
		c.Error(errors.New("the method is not supported for this route"))
	})
	if registeredRoutes {
		engine.Use(handler)
		return
	}
	deferedHandler("ERROR_ADVISERS", handler)
}

func (e ErrorAdviser) AttributesAnnotations() map[string]map[string]any {
	m := make(map[string]map[string]any)

	sm := make(map[string]any)
	sm["ErrorAdvise"] = annotation.ErrorAdvise{
		StatusCode: 404,
		Errors: []string{
			"404 page not found",
			"found with the specified id",
		},
	}
	m["RouteNotFound"] = sm

	sa := make(map[string]any)
	sa["ErrorAdvise"] = annotation.ErrorAdvise{
		StatusCode: 501,
		Errors: []string{
			"the method is not supported for this route",
		},
	}
	m["MethodNotImplemented"] = sa

	s3 := make(map[string]any)
	s3["ErrorAdvise"] = annotation.ErrorAdvise{
		StatusCode: 403,
		Errors: []string{
			"sql injection attack detected",
			"user details validation failed",
			"you do not have the required role",
		},
	}
	m["ForbiddenErrors"] = s3

	s4 := make(map[string]any)
	s4["ErrorAdvise"] = annotation.ErrorAdvise{
		StatusCode: 401,
		Errors: []string{
			"authorization token is missing",
			"validation failed for the request",
			"the authorization token has expired",
			"the authorization token is malformed",
		},
	}
	m["UnauthorizedErrors"] = s4

	s5 := make(map[string]any)
	s5["ErrorAdvise"] = annotation.ErrorAdvise{
		StatusCode: 400,
		Errors: []string{
			"invalid request payload",
		},
	}
	m["BadRequestErrors"] = s5

	s6 := make(map[string]any)
	s6["ErrorAdvise"] = annotation.ErrorAdvise{
		StatusCode: 405,
		Errors: []string{
			"route is not supported",
		},
	}
	m["MethodNotAllowedErrors"] = s6

	s7 := make(map[string]any)
	s7["ErrorAdvise"] = annotation.ErrorAdvise{
		StatusCode: 400,
		Errors: []string{
			"error in your SQL syntax",
		},
	}
	m["DatabaseErrors"] = s7

	return m
}

func (e ErrorAdviser) processErrorResponse(err error, errs []string, logger any) any {
	if logger != nil {
		if log, ok := util.GetDeclaredFieldValue(logger, "Log"); ok && !log.IsNil() {
			log.Call([]reflect.Value{reflect.ValueOf(util.StrFormat("[barmoury.ErrorAdviser] %s", errs[0])), reflect.ValueOf(err)})
		}
	}
	return model.NewApiResponseError(errs, "")
}

// @ErrorAdvise{ Errors: ["404 page not found", "found with the specified id"], StatusCode: 404 }
func (e ErrorAdviser) RouteNotFound(err error, opts ErrorAdviserOption) any {
	return e.processErrorResponse(err, []string{err.Error()}, opts)
}

// @ErrorAdvise{ Errors: ["the method is not supported for this route"], StatusCode: 501 }
func (e ErrorAdviser) MethodNotImplemented(err error, opts ErrorAdviserOption) any {
	return e.processErrorResponse(err, []string{err.Error()}, opts)
}

// @ErrorAdvise{ Errors: ["authorization token is missing", "validation failed for the request", "the authorization token has expired", "the authorization token is malformed"], StatusCode: 401 }
func (e ErrorAdviser) UnauthorizedErrors(err error, opts ErrorAdviserOption) any {
	return e.processErrorResponse(err, []string{err.Error()}, opts)
}

// @ErrorAdvise{ Errors: ["sql injection attack detected", "user details validation failed", "you do not have the required role"], StatusCode: 403 }
func (e ErrorAdviser) ForbiddenErrors(err error, opts ErrorAdviserOption) any {
	return e.processErrorResponse(err, []string{err.Error()}, opts)
}

// @ErrorAdvise{ Errors: ["invalid request payload"], StatusCode: 400 }
func (e ErrorAdviser) BadRequestErrors(err error, opts ErrorAdviserOption) any {
	return e.processErrorResponse(err, []string{err.Error()}, opts)
}

// @ErrorAdvise{ Errors: ["route is not supported"], StatusCode: 405 }
func (e ErrorAdviser) MethodNotAllowedErrors(err error, opts ErrorAdviserOption) any {
	return e.processErrorResponse(err, []string{err.Error()}, opts)
}

// @ErrorAdvise{ Errors: ["error in your SQL syntax"], StatusCode: 400 }
func (e ErrorAdviser) DatabaseErrors(err error, opts ErrorAdviserOption) any {
	return e.processErrorResponse(err, []string{"an error occur during persistence, check your request"}, opts)
}
