# Go-Plugger

This package is a wrapper for a [go-swagger](https://github.com/go-swagger/go-swagger) which makes is plugable as just a server.

> WARNING: This package is under development and is not yet ready for use!

## Idea

The go-swagger generated server is designed to be a central part of the application. It isn't easy to say: "okay, lets switch from a hand-written server to a code-generated one".
But that's a common case.

This package is't designed to make go-swagger server look like a standard `net/http` server, but makes it plugable as easy as possible.

## TODO:

- [ ] Add middleware routing
- [ ] Add easy logging
- [ ] Add instrumentation support (tracing)
- [ ] Add instrumentation support (monitoring)
- [ ] Cover all with tests
