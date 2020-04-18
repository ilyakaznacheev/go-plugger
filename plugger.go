package plugger

import (
	"net/http"
	"reflect"
)

// Plug is a Swagger API wrapper to make it pluggable
type Plug struct {
	s  Server
	sv reflect.Value
}

// NewPlug creates a new Swagger API plug
func NewPlug(s Server, opts ...Option) *Plug {
	// the current go-swagger version doesn't provide
	// access to the exported fields via methods,
	// so the only way to do so dynamically is reflection.
	sv := reflect.ValueOf(s)

	p := &Plug{
		s:  s,
		sv: sv,
	}

	// apply options
	for _, opt := range opts {
		opt.apply(p)
	}

	return p
}

// SetMiddleware adds middleware for certain path
func (p *Plug) SetMiddleware(path string, mw func(http.Handler) http.Handler, methods ...string) *Plug {
	for _, m := range methods {
		p.s.AddMiddlewareFor(m, path, mw)
	}
	return p
}

func (p *Plug) Serve() error {
	return p.s.Serve()
}
