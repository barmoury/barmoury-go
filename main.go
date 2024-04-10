package main

import (
	"barmoury/api/model"
	"barmoury/trace"
	"fmt"
)

type User struct {
	//model model.Model
	Name string
}

func testModel() {
	model1 := new(model.Model)
	model2 := model.Model{}
	model3 := model.Model{Id: 40}
	//user1 := new(User)
	//user2 := User{Name: "Two"}
	model2.Resolve(nil, nil, nil)
	model2.Id = 30
	fmt.Println("The model1:", *model1.Resolve(model2, nil, nil))
	fmt.Println("The model1:", *model1.Resolve(model3, nil, nil))
	fmt.Println("The model2:", model2)
	fmt.Println("The model3:", model3)
}

func main() {
	//testModel()
	trace.Build("Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) FxiOS/8.1.1b4948 Mobile/14F89 Safari/603.2.4")
}
