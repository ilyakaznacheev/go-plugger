package plugger

import (
	"reflect"

	"github.com/go-chi/chi"
)

// Plug is a Swagger API wrapper to make it plugable
type Plug struct {
	s  Server
	sv reflect.Value

	api  API
	apiv reflect.Value

	r chi.Router
}

// NewPlug creates a new Swagger API plug
//
// srv is a code-generated go-swagger server.
// api is a code-generated go-swagger api
// opts is a set of api or server options of your choice
//
// Note that if you predefine any fields of
// the API or the Server, options will override them
//
// You can also still use the server structure to parse
// command-line flags using "github.com/jessevdk/go-flags"
// library as go-swagger does
func NewPlug(srv Server, api API, opts ...Option) *Plug {
	// the current go-swagger version doesn't provide
	// access to the exported fields via methods,
	// so the only way to do so dynamically is reflection.
	sv := reflect.ValueOf(srv)
	apiv := reflect.ValueOf(api)

	// we use Chi router to let the user set up
	// some middleware or do anything he or she
	// wants to do
	r := chi.NewRouter()
	r.Mount("/", api.Serve(nil))

	p := &Plug{
		s:    srv,
		sv:   sv,
		api:  api,
		apiv: apiv,
		r:    r,
	}

	// apply API options
	for _, opt := range opts {
		opt.applyAPI(p)
	}

	setServerAPI(srv, api)

	// apply server options
	for _, opt := range opts {
		opt.applyServer(p)
	}

	return p
}

// Serve the API
func (p *Plug) Serve() error {
	p.s.SetHandler(p.r)
	return p.s.Serve()
}

// Shutdown server and clean up resources
func (p *Plug) Shutdown() error {
	return p.s.Shutdown()
}

// Router returns a built-in router
//
// Note that router is bound to the root address.
// So if you specified any basePath in the swagger.yml
// you have to provide a full path to the router.
// E.g. if you set basePath: /api
// the route should start with /api as well.
func (p *Plug) Router() chi.Router {
	return p.r
}

// setServerAPI dynamically calls the method
// server.SetAPI(api)
func setServerAPI(srv Server, api API) {
	apiVal := reflect.ValueOf(api)

	reflect.ValueOf(srv).
		MethodByName("SetAPI").
		Call([]reflect.Value{apiVal})
}
