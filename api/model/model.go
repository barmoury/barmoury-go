package model

import (
	"time"

	"github.com/barmoury/barmoury-go/api/timeo"
	"github.com/barmoury/barmoury-go/copier"
)

type Model struct {
	Id        uint64    `json:"id" gorm:"primary_key"`
	UpdatedAt time.Time `json:"updated_at" gorm:"<-:false"`
	CreatedAt time.Time `json:"created_at" gorm:"<-:false"`
}

// baseRequest request
// queryArmoury QueryAmoury
// userDetails any
func (model *Model) Resolve(baseRequest any, queryArmoury any, userDetails any) *Model {
	copier.Copy(model, baseRequest)
	timeo.Resolve(model)
	return model
}

type Request struct {
	___BARMOURY_UPDATE_ENTITY_ID___ uint64
}
