# Go-Plugger

This package is a wrapper for a [go-swagger](https://github.com/go-swagger/go-swagger) which makes is plugable as just a server.

> WARNING: This package is under development and is not ready for production use!

## Idea

The go-swagger generated server is designed to be a central part of the application. It isn't easy to say: "okay, lets switch from a hand-written server to a code-generated one".
But that's a common case.

This package is't designed to make go-swagger server look like a standard `net/http` server, but makes it plugable as easy as possible.