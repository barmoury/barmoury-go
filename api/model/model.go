package model

import (
	"time"

	"github.com/barmoury/barmoury-go/api/timeo"
	"github.com/barmoury/barmoury-go/copier"
	"github.com/barmoury/barmoury-go/eloquent"
)

type Model struct {
	Id        any       `json:"id" gorm:"primary_key"`
	UpdatedAt time.Time `json:"updated_at" gorm:"<-:false"`
	CreatedAt time.Time `json:"created_at" gorm:"<-:false"`
}

// baseRequest any
// queryArmoury QueryAmoury
// userDetails any
func (model *Model) Resolve(baseRequest any, queryArmoury eloquent.QueryArmoury, userDetails any) *Model {
	copier.Copy(model, baseRequest)
	timeo.Resolve(model)
	return model
}

type Request interface {
	//___BARMOURY_UPDATE_ENTITY_ID___ uint64
}
