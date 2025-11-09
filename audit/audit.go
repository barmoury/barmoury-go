package audit

import (
	"time"

	"github.com/barmoury/barmoury-go/api/timeo"
	"github.com/barmoury/barmoury-go/copier"
	"github.com/barmoury/barmoury-go/eloquent"
	"github.com/barmoury/barmoury-go/trace"
)

type Audit[T any] struct {
	Id          uint           `json:"id" gorm:"primary_key" copy_property:"ignore" request_param_filters:"field_is_snake_case,column_is_snake_case,operator=EQ"`
	Type        string         `json:"type" binding:"required" request_param_filters:"field_is_snake_case,column_is_snake_case,operator=EQ|field_is_snake_case,column_is_snake_case,operator=LIKE"`
	Group       string         `json:"group"  request_param_filters:"field_is_snake_case,column_is_snake_case,operator=EQ|field_is_snake_case,column_is_snake_case,operator=LIKE"`
	Status      string         `json:"status" request_param_filters:"field_is_snake_case,column_is_snake_case,operator=EQ|field_is_snake_case,column_is_snake_case,operator=LIKE"`
	Source      string         `json:"source" binding:"required" request_param_filters:"field_is_snake_case,column_is_snake_case,operator=EQ|field_is_snake_case,column_is_snake_case,operator=LIKE"`
	Action      string         `json:"action" binding:"required" request_param_filters:"field_is_snake_case,column_is_snake_case,operator=EQ|field_is_snake_case,column_is_snake_case,operator=LIKE"`
	ActorId     string         `json:"actor_id" request_param_filters:"field_is_snake_case,column_is_snake_case,operator=EQ|field_is_snake_case,column_is_snake_case,operator=LIKE"`
	ActorType   string         `json:"actor_type" request_param_filters:"field_is_snake_case,column_is_snake_case,operator=EQ|field_is_snake_case,column_is_snake_case,operator=LIKE"`
	IpAddress   string         `json:"ip_address" request_param_filters:"field_is_snake_case,column_is_snake_case,operator=EQ|field_is_snake_case,column_is_snake_case,operator=LIKE"`
	Environment string         `json:"environment" request_param_filters:"field_is_snake_case,column_is_snake_case,operator=EQ|field_is_snake_case,column_is_snake_case,operator=LIKE"`
	AuditId     string         `json:"audit_id" request_param_filters:"field_is_snake_case,column_is_snake_case,operator=EQ|field_is_snake_case,column_is_snake_case,operator=LIKE"`
	Device      trace.Device   `json:"device" gorm:"serializer:json" request_param_filters:"field_is_snake_case,column_is_snake_case,operator=LIKE|field_is_snake_case,column_is_snake_case,operator=OBJECT_LIKE"`
	Isp         trace.Isp      `json:"isp" gorm:"serializer:json" request_param_filters:"field_is_snake_case,column_is_snake_case,operator=LIKE|field_is_snake_case,column_is_snake_case,operator=OBJECT_LIKE"`
	Location    trace.Location `json:"location" gorm:"serializer:json" request_param_filters:"field_is_snake_case,column_is_snake_case,operator=LIKE|field_is_snake_case,column_is_snake_case,operator=OBJECT_LIKE"`
	Auditable   any            `json:"auditable" gorm:"serializer:json" request_param_filters:"field_is_snake_case,column_is_snake_case,operator=LIKE"`
	ExtraData   any            `json:"extra_data" gorm:"serializer:json" request_param_filters:"field_is_snake_case,column_is_snake_case,operator=LIKE|field_is_snake_case,column_is_snake_case,operator=OBJECT_LIKE"`
	CreatedAt   time.Time      `json:"created_at,omitempty" gorm:"<-:false" request_param_filters:"field_is_snake_case,column_is_snake_case,operator=RANGE"`
}

func (Audit[T]) TableName() string {
	return "audits"
}

func (model *Audit[T]) Resolve(baseRequest any, queryArmoury eloquent.QueryArmoury, userDetails any) *Audit[T] {
	copier.Copy(model, baseRequest)
	timeo.Resolve(model)
	return model
}
