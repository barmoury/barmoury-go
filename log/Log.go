package log

import "time"

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
	Id        uint      `json:"id" gorm:"primary_key"`
	Level     Level     `json:"level" binding:"required"`
	Group     string    `json:"group" binding:"required"`
	Source    string    `json:"source" binding:"required"`
	TraceId   string    `json:"trace_id"`
	SpanId    string    `json:"span_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at" gorm:"<-:false"`
}
