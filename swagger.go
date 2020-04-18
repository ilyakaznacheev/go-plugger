package plugger

import (
	"net"
	"net/http"
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
