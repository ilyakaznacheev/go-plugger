package plugger

import (
	"net"
	"net/http"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
)

// Server is a set of functions of a swagger-generated server
type Server interface {
	ConfigureAPI()
	Logf(f string, args ...interface{})
	Fatalf(f string, args ...interface{})
	Serve() (err error)
	Listen() error
	Shutdown() error
	GetHandler() http.Handler
	SetHandler(handler http.Handler)
	UnixListener() (net.Listener, error)
	HTTPListener() (net.Listener, error)
	TLSListener() (net.Listener, error)
}

// API is a set of functions of a swagger-generated API
type API interface {
	SetDefaultProduces(mediaType string)
	SetDefaultConsumes(mediaType string)
	SetSpec(spec *loads.Document)
	DefaultProduces() string
	DefaultConsumes() string
	Formats() strfmt.Registry
	RegisterFormat(name string, format strfmt.Format, validator strfmt.Validator)
	Validate() error
	ServeErrorFor(operationID string) func(http.ResponseWriter, *http.Request, error)
	AuthenticatorsFor(schemes map[string]spec.SecurityScheme) map[string]runtime.Authenticator
	Authorizer() runtime.Authorizer
	ConsumersFor(mediaTypes []string) map[string]runtime.Consumer
	ProducersFor(mediaTypes []string) map[string]runtime.Producer
	HandlerFor(method, path string) (http.Handler, bool)
	Context() *middleware.Context
	Serve(builder middleware.Builder) http.Handler
	Init()
	RegisterConsumer(mediaType string, consumer runtime.Consumer)
	RegisterProducer(mediaType string, producer runtime.Producer)
	AddMiddlewareFor(method, path string, builder middleware.Builder)
}
