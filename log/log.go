package log

import (
	"time"

	"github.com/barmoury/barmoury-go/api/timeo"
	"github.com/barmoury/barmoury-go/copier"
	"github.com/barmoury/barmoury-go/eloquent"
)

type Level string

const (
	INFO    Level = "INFO"
	WARN    Level = "WARN"
	ERROR   Level = "ERROR"
	TRACE   Level = "TRACE"
	FATAL   Level = "FATAL"
	PANIC   Level = "PANIC"
	VERBOSE Level = "VERBOSE"
)

type Log struct {
	Id        uint      `json:"id" gorm:"primary_key" copy_property:"ignore" request_param_filters:"field_is_snake_case,column_is_snake_case,operator=EQ"`
	Level     Level     `json:"level" binding:"required" request_param_filters:"field_is_snake_case,column_is_snake_case,operator=EQ|field_is_snake_case,column_is_snake_case,operator=LIKE"`
	Group     string    `json:"group,omitempty" request_param_filters:"field_is_snake_case,column_is_snake_case,operator=EQ|field_is_snake_case,column_is_snake_case,operator=LIKE"`
	Source    string    `json:"source" request_param_filters:"field_is_snake_case,column_is_snake_case,operator=EQ|field_is_snake_case,column_is_snake_case,operator=LIKE"`
	TraceId   string    `json:"trace_id,omitempty" request_param_filters:"field_is_snake_case,column_is_snake_case,operator=EQ|field_is_snake_case,column_is_snake_case,operator=LIKE"`
	SpanId    string    `json:"span_id,omitempty" request_param_filters:"field_is_snake_case,column_is_snake_case,operator=EQ|field_is_snake_case,column_is_snake_case,operator=LIKE"`
	Content   string    `json:"content" binding:"required" request_param_filters:"field_is_snake_case,column_is_snake_case,operator=EQ|field_is_snake_case,column_is_snake_case,operator=LIKE"`
	CreatedAt time.Time `json:"created_at" gorm:"<-:false" request_param_filters:"field_is_snake_case,column_is_snake_case,operator=RANGE"`
}

func (model *Log) Resolve(baseRequest any, queryArmoury eloquent.QueryArmoury, userDetails any) *Log {
	copier.Copy(model, baseRequest)
	timeo.Resolve(model)
	return model
}
