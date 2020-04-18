package plugger

import (
	"reflect"
	"time"

	"github.com/go-openapi/runtime/flagext"
	"github.com/jessevdk/go-flags"
)

type Option interface {
	apply(*Plug)
}

type funcOption struct {
	f func(*Plug)
}

func (fo *funcOption) apply(do *Plug) {
	fo.f(do)
}

func newOption(f func(*Plug)) *funcOption {
	return &funcOption{
		f: f,
	}
}

func setServerParam(srvVal reflect.Value, key string, value interface{}) {
	if rHost := reflect.Indirect(srvVal).FieldByName(key); rHost.CanSet() {
		rHost.Set(reflect.ValueOf(value))
	}
}

func newParamOption(key string, value interface{}) *funcOption {
	return newOption(func(p *Plug) {
		setServerParam(p.sv, key, value)
	})
}

// WithPort the port to listen on for insecure connections, defaults to a random value
func WithPort(port int) Option {
	return newParamOption("Port", port)
}

// WithHost the IP to listen on
func WithHost(host int) Option {
	return newParamOption("Host", host)
}

// WithEnabledListeners the listeners to enable, this can be repeated and defaults to the schemes in the swagger spec
func WithEnabledListeners(listeners []string) Option {
	return newParamOption("EnabledListeners", listeners)
}

// WithCleanupTimeout grace period for which to wait before killing idle connections
func WithCleanupTimeout(t time.Duration) Option {
	return newParamOption("CleanupTimeout", t)
}

// WithGracefulTimeout grace period for which to wait before shutting down the server
func WithGracefulTimeout(t time.Duration) Option {
	return newParamOption("GracefulTimeout", t)
}

// WithMaxHeaderSize controls the maximum number of bytes the server
// will read parsing the request header's keys and values,
// including the request line. It does not limit the size of the request body.
func WithMaxHeaderSize(size int) Option {
	return newParamOption("MaxHeaderSize", flagext.ByteSize(size))
}

// WithSocketPath the unix socket to listen on
func WithSocketPath(path string) Option {
	return newParamOption("SocketPath", flags.Filename(path))
}

// WithListenLimit limit the number of outstanding requests
func WithListenLimit(limit int) Option {
	return newParamOption("ListenLimit", limit)
}

// WithKeepAlive sets the TCP keep-alive timeouts on accepted connections.
// It prunes dead TCP connections ( e.g. closing laptop mid-download)
func WithKeepAlive(t time.Duration) Option {
	return newParamOption("KeepAlive", t)
}

// WithReadTimeout maximum duration before timing out read of the request
func WithReadTimeout(t time.Duration) Option {
	return newParamOption("ReadTimeout", t)
}

// WithWriteTimeout maximum duration before timing out write of the response
func WithWriteTimeout(t time.Duration) Option {
	return newParamOption("WriteTimeout", t)
}

// WithTLSHost the IP to listen on for tls
func WithTLSHost(host string) Option {
	return newParamOption("TLSHost", host)
}

// WithTLSPort the port to listen on for secure connections,
// defaults to a random value
func WithTLSPort(port int) Option {
	return newParamOption("TLSPort", port)
}

// WithTLSCertificate the certificate to use for secure connections
func WithTLSCertificate(cert string) Option {
	return newParamOption("TLSCertificate", flags.Filename(cert))
}

// WithTLSCertificateKeyTLSCertificateKey the private key to use for secure connections
func WithTLSCertificateKey(cert string) Option {
	return newParamOption("TLSCertificateKey", cert)
}

// WithTLSCACertificate the certificate authority file to be used with mutual tls auth
func WithTLSCACertificate(cert string) Option {
	return newParamOption("TLSCACertificate", cert)
}

// WithTLSListenLimit limit the number of outstanding requests
func WithTLSListenLimit(limit int) Option {
	return newParamOption("TLSListenLimit", limit)
}

// WithTLSKeepAlive sets the TCP keep-alive timeouts on accepted connections.
// It prunes dead TCP connections ( e.g. closing laptop mid-download)
func WithTLSKeepAlive(t time.Duration) Option {
	return newParamOption("TLSKeepAlive", t)
}

// WithTLSReadTimeout maximum duration before timing out read of the request
func WithTLSReadTimeout(t time.Duration) Option {
	return newParamOption("TLSReadTimeout", t)
}

// WithTLSWriteTimeout maximum duration before timing out write of the response
func WithTLSWriteTimeout(t time.Duration) Option {
	return newParamOption("TLSWriteTimeout", t)
}
