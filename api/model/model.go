package model

import (
	"time"

	"github.com/barmoury/barmoury-go/api/timeo"
	"github.com/barmoury/barmoury-go/copier"
	"github.com/barmoury/barmoury-go/eloquent"
)

type DyModel[T any] struct {
	Id        T         `json:"id" gorm:"primary_key"`
	UpdatedAt time.Time `json:"updated_at" gorm:"<-:false"`
	CreatedAt time.Time `json:"created_at" gorm:"<-:false"`
}

type Model struct {
	DyModel[uint]
}

type PolyModel struct {
	DyModel[any]
}

// baseRequest any
// queryArmoury QueryAmoury
// userDetails any
func (model *Model) SuperResolve(baseRequest any, queryArmoury eloquent.QueryArmoury, userDetails any, tModel any) any {
	copier.Copy(tModel, baseRequest)
	timeo.Resolve(tModel)
	return tModel
}

type Request interface {
	//___BARMOURY_UPDATE_ENTITY_ID___ uint64
}
