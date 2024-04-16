package annotation

type RequestMethod string

const (
	PUT     RequestMethod = "PUT"
	GET     RequestMethod = "GET"
	HEAD    RequestMethod = "HEAD"
	POST    RequestMethod = "POST"
	PATCH   RequestMethod = "PATCH"
	TRACE   RequestMethod = "TRACE"
	DELETE  RequestMethod = "DELETE"
	OPTIONS RequestMethod = "OPTIONS"
)

type RequestMapping struct {
	Model         any
	Request       any
	BodyScheme    any
	QuerySchema   any
	ParamsSchema  any
	HeadersSchema any
	Name          string
	Value         string
	Method        RequestMethod
}
