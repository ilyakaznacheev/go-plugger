package plugger

import (
	"net"
	"net/http"
)

// Server is a set of functions of a swagger-generated server
type Server interface {
	ConfigureAPI()
	ConfigureFlags()
	Logf(string, ...interface{})
	Fatalf(string, ...interface{})
	Serve() (err error)
	Listen() error
	Shutdown()
	GetHandler()
	SetHandler(http.Handler)
	UnixListener() (net.Listener, error)
	HTTPListener() (net.Listener, error)
	TLSListener() (net.Listener, error)
	AddMiddlewareFor(string, string, func(http.Handler) http.Handler)
}
