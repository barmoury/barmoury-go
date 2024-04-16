package main

import (
	"fmt"

	"github.com/barmoury/barmoury-go/api/annotation"
	"github.com/barmoury/barmoury-go/api/config"
	"github.com/barmoury/barmoury-go/api/controller"
	"github.com/barmoury/barmoury-go/api/model"
	"github.com/barmoury/barmoury-go/meta"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type User struct {
	model.Model
}

type UserRequest struct {
	model.Request
}

type UserController struct {
	*controller.Controller[User, UserRequest]
}

func (u UserController) Annotations() (map[string]any, string) {
	m := make(map[string]any)
	m["RequestMapping"] = annotation.RequestMapping{
		Value: "/user",
	}
	//return nil, "@RequestMapping({ value: '/user', model: User, request: UserRequest })"
	return m, ""
}

func (c UserController) AttributesAnnotations() map[string]map[string]any {
	m := make(map[string]map[string]any)

	sm := make(map[string]any)
	sm["RequestMapping"] = annotation.RequestMapping{
		Value:  "/test/:id",
		Method: annotation.PUT,
	}
	m["User"] = sm

	return m
}

func (c UserController) User(g *gin.Context) {

}

type BactuatorControllerImpl struct {
	*controller.BactuatorController
}

func (u BactuatorControllerImpl) Annotations() (map[string]any, string) {
	m := make(map[string]any)
	m["RequestMapping"] = annotation.RequestMapping{
		Value: "/bactuator",
	}
	//return nil, "@RequestMapping({ value: '/user', model: User, request: UserRequest })"
	return m, ""
}

func main() {
	/*s := &UserController{
		Controller: &controller.Controller[User, UserRequest]{},
	}
	util.SetFieldValue(util.GetDeclaredSurefireFieldFromPtr(s, "Controller").Interface(), "Self", s)
	s.Setup(nil)
	s.Stat(nil)*/
	router := gin.Default()
	opts := config.RegisterControllers(router, config.RouterOption{
		Db: &gorm.DB{},
		Bacuator: &BactuatorControllerImpl{
			BactuatorController: &controller.BactuatorController{},
		},
	}, []meta.IAnotation{
		&UserController{
			Controller: &controller.Controller[User, UserRequest]{},
		},
	})
	config.RegisterErrorAdvisers(router, config.ErrorAdviserOption{}, []any{
		&config.ErrorAdviser{},
	})
	fmt.Println("THE NEW OPTIONS", opts)
}
