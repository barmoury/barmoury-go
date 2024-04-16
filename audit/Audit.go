package audit

import (
	"time"

	"github.com/barmoury/barmoury-go/api/timeo"
	"github.com/barmoury/barmoury-go/copier"
	"github.com/barmoury/barmoury-go/eloquent"
	"github.com/barmoury/barmoury-go/trace"
)

type Audit[T any] struct {
	Id          uint           `json:"id" gorm:"primary_key" copy_property:"ignore"`
	Type        string         `json:"type" binding:"required"`
	Group       string         `json:"group" binding:"required"`
	Status      string         `json:"status"`
	Source      string         `json:"source" binding:"required"`
	Action      string         `json:"action" binding:"required"`
	ActorId     string         `json:"actor_id"`
	ActorType   string         `json:"actor_type"`
	IpAddress   string         `json:"ip_address"`
	Environment string         `json:"environment"`
	AuditId     string         `json:"audit_id"`
	Device      trace.Device   `json:"device" gorm:"serializer:json"`
	Isp         trace.Isp      `json:"isp" gorm:"serializer:json"`
	Location    trace.Location `json:"location" gorm:"serializer:json"`
	Auditable   any            `json:"auditable" gorm:"serializer:json"`
	ExtraData   any            `json:"extra_data" gorm:"serializer:json"`
	CreatedAt   time.Time      `json:"created_at,omitempty" gorm:"<-:false"`
}

func (Audit[T]) TableName() string {
	return "audits"
}

func (model *Audit[T]) Resolve(baseRequest any, queryArmoury eloquent.QueryArmoury, userDetails any) *Audit[T] {
	copier.Copy(model, baseRequest)
	timeo.Resolve(model)
	return model
}
