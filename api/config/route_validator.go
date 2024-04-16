package config

import (
	"errors"

	"github.com/barmoury/barmoury-go/util"
	"github.com/gin-gonic/gin"
)

type RouteValidator struct {
	Prefix string
	Routes []IRoute
	Valid  func(*gin.Context) bool
}

func RegisterRouteValidators(engine *gin.Engine, validators []RouteValidator) {
	handler := func() gin.HandlerFunc {
		mappedRouteValidators := map[string]func(*gin.Context) bool{}
		for _, validator := range validators {
			prefix := util.If(validator.Prefix != "", validator.Prefix+"/", "")
			for _, route := range validator.Routes {
				p := util.ReplaceByRegex(prefix+route.Route, `([^:]\/)\/+`, `$1`)
				k := util.If(route.Method == "", "ANY", route.Method) + "<=#=>" + p
				mappedRouteValidators[k] = validator.Valid
			}
		}
		return func(c *gin.Context) {
			var ok bool
			var valid func(*gin.Context) bool
			if valid, ok = mappedRouteValidators[c.Request.Method+"<=#=>"+c.Request.URL.Path]; !ok {
				if valid, ok = mappedRouteValidators[c.Request.Method+"<=#=>"+c.Request.URL.Path]; !ok {
					return
				}
			}
			if !valid(c) {
				c.Error(errors.New("validation failed for the request"))
			}
		}
	}
	if registeredRoutes {
		engine.Use(handler())
		return
	}
	deferedHandler("ROUTE_VALIDATORS", handler())
}
