package plugger

import (
	"reflect"
)

// Plug is a Swagger API wrapper to make it pluggable
type Plug struct {
	s  Server
	sv reflect.Value
}

// NewPlug creates a new Swagger API plug
//
// a server is a code-generated go-swagger server.
//
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

func (p *Plug) Serve() error {
	return p.s.Serve()
}

func (p *Plug) Shutdown() error {
	return p.s.Shutdown()
}
