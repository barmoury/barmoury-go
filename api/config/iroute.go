package config

import (
	"strings"

	"github.com/barmoury/barmoury-go/util"
	"github.com/gin-gonic/gin"
)

type IRoute struct {
	Route  string
	Method string
}

func ShouldNotFilter(c *gin.Context, prefix string, openUrlPatterns []IRoute) bool {
	if c.Request.URL == nil || c.Request.URL.Path == "" {
		return false
	}
	m := c.Request.Method
	r := util.If(prefix != "", "/", "") + strings.Replace(c.Request.URL.Path, prefix, "", -1)
	for _, oup := range openUrlPatterns {
		if (oup.Method == "" || m == oup.Method) && util.PatternToRegex(oup.Route).MatchString(r) {
			return true
		}
	}
	return false
}
