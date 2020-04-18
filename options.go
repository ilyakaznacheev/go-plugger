package plugger

import (
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/docker/go-units"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/flagext"
	"github.com/go-openapi/runtime/security"
	"github.com/go-openapi/swag"
	"github.com/jessevdk/go-flags"
)

// Option is a parameter option that defines a server parameter
type Option interface {
	// applyAPI will set any swagger API parameters before server setup
	applyAPI(*Plug)
	// applyServer will set any swagger server parameters after server setup
	applyServer(*Plug)
}

type funcOption struct {
	api func(*Plug)
	srv func(*Plug)
}

func (fo *funcOption) applyAPI(do *Plug) {
	fo.api(do)
}

func (fo *funcOption) applyServer(do *Plug) {
	fo.srv(do)
}

func dummyFunc(*Plug) {}

func newOptionAPI(f func(*Plug)) *funcOption {
	return &funcOption{
		api: f,
		srv: dummyFunc,
	}
}

func newOptionServer(f func(*Plug)) *funcOption {
	return &funcOption{
		srv: f,
		api: dummyFunc,
	}
}

// WithConfiguredAPI runs generated API configuration function of swagger server.
// Normally it runs functions defined in `configure_*.go` file, where you has to
// define your handlers and settings in a classic go-swagger generated server.
//
// You may add this option if you want to setup some defaults defined there.
// But you don't need it if you set up the server configuration yourself
// using the Plugger
func WithConfiguredAPI() Option {
	return newOptionServer(func(p *Plug) {
		p.s.ConfigureAPI()
	})
}

func setDynParam(srvVal reflect.Value, key string, value interface{}) {
	if rHost := reflect.Indirect(srvVal).FieldByName(key); rHost.CanSet() {
		rHost.Set(reflect.ValueOf(value))
	}
}

func newParamServerOption(key string, value interface{}) *funcOption {
	return newOptionServer(func(p *Plug) {
		setDynParam(p.sv, key, value)
	})
}

func newParamAPIOption(key string, value interface{}) *funcOption {
	return newOptionAPI(func(p *Plug) {
		setDynParam(p.apiv, key, value)
	})
}

// WithPort the port to listen on for insecure connections, defaults to a random value
func WithPort(port int) Option {
	return newParamServerOption("Port", port)
}

// WithHost the IP to listen on
func WithHost(host string) Option {
	return newParamServerOption("Host", host)
}

// WithEnabledListeners the listeners to enable, this can be repeated and defaults to the schemes in the swagger spec
func WithEnabledListeners(listeners []string) Option {
	return newParamServerOption("EnabledListeners", listeners)
}

// WithCleanupTimeout grace period for which to wait before killing idle connections
func WithCleanupTimeout(t time.Duration) Option {
	return newParamServerOption("CleanupTimeout", t)
}

// WithGracefulTimeout grace period for which to wait before shutting down the server
func WithGracefulTimeout(t time.Duration) Option {
	return newParamServerOption("GracefulTimeout", t)
}

// WithMaxHeaderSize controls the maximum number of bytes the server
// will read parsing the request header's keys and values,
// including the request line. It does not limit the size of the request body.
func WithMaxHeaderSize(size int) Option {
	return newParamServerOption("MaxHeaderSize", flagext.ByteSize(size))
}

// WithSocketPath the unix socket to listen on
func WithSocketPath(path string) Option {
	return newParamServerOption("SocketPath", flags.Filename(path))
}

// WithListenLimit limit the number of outstanding requests
func WithListenLimit(limit int) Option {
	return newParamServerOption("ListenLimit", limit)
}

// WithKeepAlive sets the TCP keep-alive timeouts on accepted connections.
// It prunes dead TCP connections ( e.g. closing laptop mid-download)
func WithKeepAlive(t time.Duration) Option {
	return newParamServerOption("KeepAlive", t)
}

// WithReadTimeout maximum duration before timing out read of the request
func WithReadTimeout(t time.Duration) Option {
	return newParamServerOption("ReadTimeout", t)
}

// WithWriteTimeout maximum duration before timing out write of the response
func WithWriteTimeout(t time.Duration) Option {
	return newParamServerOption("WriteTimeout", t)
}

// WithTLSHost the IP to listen on for tls
func WithTLSHost(host string) Option {
	return newParamServerOption("TLSHost", host)
}

// WithTLSPort the port to listen on for secure connections,
// defaults to a random value
func WithTLSPort(port int) Option {
	return newParamServerOption("TLSPort", port)
}

// WithTLSCertificate the certificate to use for secure connections
func WithTLSCertificate(cert string) Option {
	return newParamServerOption("TLSCertificate", flags.Filename(cert))
}

// WithTLSCertificateKeyTLSCertificateKey the private key to use for secure connections
func WithTLSCertificateKey(cert string) Option {
	return newParamServerOption("TLSCertificateKey", cert)
}

// WithTLSCACertificate the certificate authority file to be used with mutual tls auth
func WithTLSCACertificate(cert string) Option {
	return newParamServerOption("TLSCACertificate", cert)
}

// WithTLSListenLimit limit the number of outstanding requests
func WithTLSListenLimit(limit int) Option {
	return newParamServerOption("TLSListenLimit", limit)
}

// WithTLSKeepAlive sets the TCP keep-alive timeouts on accepted connections.
// It prunes dead TCP connections ( e.g. closing laptop mid-download)
func WithTLSKeepAlive(t time.Duration) Option {
	return newParamServerOption("TLSKeepAlive", t)
}

// WithTLSReadTimeout maximum duration before timing out read of the request
func WithTLSReadTimeout(t time.Duration) Option {
	return newParamServerOption("TLSReadTimeout", t)
}

// WithTLSWriteTimeout maximum duration before timing out write of the response
func WithTLSWriteTimeout(t time.Duration) Option {
	return newParamServerOption("TLSWriteTimeout", t)
}

// WithBasicAuthenticator generates a runtime.Authenticator from the supplied basic auth function.
// It has a default implementation in the security package, however you can replace it for your particular usage.
func WithBasicAuthenticator(
	f func(security.UserPassAuthentication) runtime.Authenticator) Option {
	return newParamServerOption("BasicAuthenticator", f)
}

// WithAPIKeyAuthenticator generates a runtime.Authenticator from the supplied token auth function.
// It has a default implementation in the security package, however you can replace it for your particular usage.
func WithAPIKeyAuthenticator(
	f func(string, string, security.TokenAuthentication) runtime.Authenticator) Option {
	return newParamServerOption("APIKeyAuthenticator", f)
}

// WithBearerAuthenticator BearerAuthenticator generates a runtime.Authenticator from the supplied bearer token auth function.
// It has a default implementation in the security package, however you can replace it for your particular usage.
func WithBearerAuthenticator(
	f func(string, security.ScopedTokenAuthentication) runtime.Authenticator) Option {
	return newParamServerOption("BearerAuthenticator", f)
}

// WithJSONConsumer JSONConsumer registers a consumer for the following mime types:
//   - application/json
func WithJSONConsumer(c runtime.Consumer) Option {
	return newParamServerOption("JSONConsumer", c)
}

// WithBinProducer BinProducer registers a producer for the following mime types:
//   - application/octet-stream
func WithBinProducer(p runtime.Producer) Option {
	return newParamServerOption("BinProducer", p)
}

// WithHTMLProducer HTMLProducer registers a producer for the following mime types:
//   - text/html
func WithHTMLProducer(p runtime.Producer) Option {
	return newParamServerOption("HTMLProducer", p)
}

// WithJSONProducer JSONProducer registers a producer for the following mime types:
//   - application/json
func WithJSONProducer(p runtime.Producer) Option {
	return newParamServerOption("JSONProducer", p)
}

// WithServeError ServeError is called when an error is received, there is a default handler
// but you can set your own with this
func WithServeError(f func(http.ResponseWriter, *http.Request, error)) Option {
	return newParamServerOption("ServeError", f)
}

// WithPreServerShutdown PreServerShutdown is called before the HTTP(S) server is shutdown
// This allows for custom functions to get executed before the HTTP(S) server stops accepting traffic
func WithPreServerShutdown(f func()) Option {
	return newParamServerOption("PreServerShutdown", f)
}

// WithServerShutdown ServerShutdown is called when the HTTP(S) server is shut down and done
// handling all active connections and does not accept connections any more
func WithServerShutdown(f func()) Option {
	return newParamServerOption("ServerShutdown", f)
}

// WithCommandLineOptionsGroups Custom command line argument groups with their descriptions
func WithCommandLineOptionsGroups(g []swag.CommandLineOptionsGroup) Option {
	return newParamServerOption("CommandLineOptionsGroups", g)
}

// WithLogger User defined logger function
func WithLogger(f func(string, ...interface{})) Option {
	return newParamServerOption("Logger", f)
}

// WithAPIDefaults sets default values to API fields
func WithAPIDefaults() Option {
	return newOptionAPI(func(p *Plug) {
		setDynParam(p.apiv, "BasicAuthenticator", security.BasicAuth)
		setDynParam(p.apiv, "APIKeyAuthenticator", security.APIKeyAuth)
		setDynParam(p.apiv, "BearerAuthenticator", security.BearerAuth)
		setDynParam(p.apiv, "JSONConsumer", runtime.JSONConsumer())
		setDynParam(p.apiv, "BinProducer", runtime.ByteStreamProducer())
		setDynParam(p.apiv, "JSONProducer", runtime.JSONProducer())
		setDynParam(p.apiv, "ServeError", errors.ServeError)
		setDynParam(p.apiv, "HTMLProducer", runtime.ProducerFunc(func(w io.Writer, data interface{}) error {
			return errors.NotImplemented("html producer has not yet been implemented")
		}))
	})
}

// WithServerDefaults sets default values to server fields
func WithServerDefaults() Option {
	return newOptionServer(func(p *Plug) {
		setDynParam(p.sv, "CleanupTimeout", 10*time.Second)
		setDynParam(p.sv, "GracefulTimeout", 15*time.Second)
		setDynParam(p.sv, "MaxHeaderSize", units.MiB)
		setDynParam(p.sv, "SocketPath", flags.Filename("/var/run/backend-storage.sock"))
		setDynParam(p.sv, "KeepAlive", 3*time.Minute)
		setDynParam(p.sv, "ReadTimeout", 30*time.Second)
		setDynParam(p.sv, "WriteTimeout", 60*time.Second)
	})
}
