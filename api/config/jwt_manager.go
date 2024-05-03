package config

import (
	"errors"
	"strings"

	"github.com/barmoury/barmoury-go/api/model"
	"github.com/barmoury/barmoury-go/crypto"
	"github.com/barmoury/barmoury-go/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type JwtManagerOption struct {
	Prefix          string
	AuthorityPrefix string
	OpenUrlPatterns []IRoute
	Secrets         map[string]string
	Encryptor       crypto.IEncryptor[any]
	Validate        func(*gin.Context, string, model.UserDetails[any]) bool
}

func RegisterJwt(engine *gin.Engine, opts JwtManagerOption) {
	handler := func() gin.HandlerFunc {
		signer := func(authToken string, c *gin.Context) string {
			token, grp, err := findActiveToken(authToken, opts.Secrets)
			if err != "" {
				return err
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				return "invalid authorization token, claims not retrievable"
			}
			if opts.Encryptor != nil {
				for key, claim := range claims {
					if util.GetTypeName(claim) != "string" {
						continue
					}
					claims[key], _ = opts.Encryptor.Decrypt(claim.(string))
				}
			}
			var sub any
			var bas []string
			var bd map[string]interface{}
			if sub_, ok := claims["sub"]; ok {
				sub = sub_
				if bd_, ok := claims["BARMOURY_DATA"]; ok {
					bd = bd_.(map[string]interface{})
					if bas_, ok := claims["BARMOURY_AUTHORITIES"]; ok {
						for _, v := range bas_.([]interface{}) {
							bas = append(bas, v.(string))
						}
						ud := model.UserDetails[any]{
							Data:              bd,
							AuthoritiesValues: bas,
							Id:                sub.(string),
							AuthorityPrefix:   opts.AuthorityPrefix,
						}
						if opts.Validate != nil && !opts.Validate(c, grp, ud) {
							return "user details validation failed"
						}
						c.Set("user", ud)
						c.Set("authoritiesValues", ud.AuthoritiesValues)
						return ""
					}
				}
			}
			c.Set("user", claims)
			return ""
		}
		return func(c *gin.Context) {
			if len(opts.OpenUrlPatterns) > 0 && ShouldNotFilter(c, opts.Prefix, opts.OpenUrlPatterns) {
				return
			}
			atp := strings.Split(c.GetHeader("Authorization"), " ")
			if len(atp) < 2 || atp[1] == "" {
				panic(errors.New("authorization token is missing"))
			}
			s := signer(atp[1], c)
			if s != "" {
				c.Error(errors.New(s))
				c.Abort()
				return
			}
		}
	}
	if registeredRoutes {
		engine.Use(handler())
		return
	}
	deferedHandler("JWT_MANAGERS", handler())
}

func findActiveToken(authToken string, secrets map[string]string) (*jwt.Token, string, string) {
	for group, secret := range secrets {
		token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil {
			if strings.Contains(err.Error(), "expired") {
				return nil, group, "the authorization token has expired"
			} else if strings.Contains(err.Error(), "malformed") {
				return nil, group, "the authorization token is malformed"
			}
		} else {
			return token, group, ""
		}
	}
	return nil, "", "invalid authorization token"
}
