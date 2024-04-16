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
	Id        uint      `json:"id" gorm:"primary_key" copy_property:"ignore"`
	Level     Level     `json:"level" binding:"required"`
	Group     string    `json:"group,omitempty"`
	Source    string    `json:"source"`
	TraceId   string    `json:"trace_id,omitempty"`
	SpanId    string    `json:"span_id,omitempty"`
	Content   string    `json:"content" binding:"required"`
	CreatedAt time.Time `json:"created_at" gorm:"<-:false"`
}

func (model *Log) Resolve(baseRequest any, queryArmoury eloquent.QueryArmoury, userDetails any) *Log {
	copier.Copy(model, baseRequest)
	timeo.Resolve(model)
	return model
}
