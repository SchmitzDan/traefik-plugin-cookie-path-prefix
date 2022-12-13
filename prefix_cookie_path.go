// Package
package traefik_plugin_cookie_path_prefix

import (
	"context"
	"net/http"
	"os"
)

const setCookieHeader string = "Set-Cookie"

type Config struct {
	Prefix string `json:"prefix,omitempty" toml:"prefix,omitempty" yaml:"prefix,omitempty"`
}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

type PathPrefixer struct {
	next   http.Handler
	name   string
	prefix string
}

func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	os.Stdout.WriteString("Plugin cookie path prefix started")

	//TODO: Check config?

	return &PathPrefixer{
		name:   name,
		next:   next,
		prefix: config.Prefix,
	}, nil
}

func (p *PathPrefixer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	myWriter := &responseWriter{
		writer: rw,
		prefix: p.prefix,
	}

	p.next.ServeHTTP(myWriter, req)
}

type responseWriter struct {
	writer http.ResponseWriter
	prefix string
}

func (r *responseWriter) Header() http.Header {
	return r.writer.Header()
}

func (r *responseWriter) Write(bytes []byte) (int, error) {
	return r.writer.Write(bytes)
}

func (r *responseWriter) WriteHeader(statusCode int) {

	//workaround to get the cookies
	headers := r.writer.Header()
	req := http.Response{Header: headers}
	cookies := req.Cookies()

	//Delete set-cookie headers
	r.writer.Header().Del(setCookieHeader)

	//Add new cookie with modifies path
	for _, cookie := range cookies {
		if cookie.Path == "/" {
			//prevent trailing /
			cookie.Path = "/" + r.prefix
		} else {
			cookie.Path = "/" + r.prefix + cookie.Path
		}
		http.SetCookie(r, cookie)
	}
}
