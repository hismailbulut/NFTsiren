package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"strings"
)

type Request struct {
	Method  string
	Path    []string
	Params  map[string]string
	Header  map[string]string
	Payload io.Reader
}

func (r *Request) String() string {
	return fmt.Sprintf("%s %s", r.Method, r.URL(&url.URL{}))
}

func NewRequest(method string, path []string) *Request {
	return &Request{Method: method, Path: path}
}

func (r *Request) URL(base *url.URL) string {
	u := base.JoinPath(r.Path...)
	q := u.Query()
	for k, v := range r.Params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func (r *Request) SetMethod(method string) *Request {
	r.Method = method
	return r
}

func (r *Request) SetPath(path []string) *Request {
	r.Path = path
	return r
}

func (r *Request) SetParam(key, value string) *Request {
	if r.Params == nil {
		r.Params = make(map[string]string)
	}
	r.Params[key] = value
	return r
}

func (r *Request) SetHeader(key, value string) *Request {
	if r.Header == nil {
		r.Header = make(map[string]string)
	}
	r.Header[key] = value
	return r
}

// shorthand
func (r *Request) SetAccept(accept string) *Request {
	return r.SetHeader("Accept", accept)
}

// shorthand
func (r *Request) SetAcceptJSON() *Request {
	return r.SetAccept("application/json")
}

// shorthand
func (r *Request) SetAcceptHTML() *Request {
	return r.SetAccept("text/html")
}

// shorthand
func (r *Request) SetContentType(contentType string) *Request {
	return r.SetHeader("Content-Type", contentType)
}

// shorthand
func (r *Request) SetContentTypeJSON() *Request {
	return r.SetContentType("application/json")
}

func (r *Request) SetPayload(payload io.Reader) *Request {
	r.Payload = payload
	return r
}

// shorthand
func (r *Request) SetPayloadBytes(payload []byte) *Request {
	return r.SetPayload(bytes.NewBuffer(payload))
}

// shorthand
func (r *Request) SetPayloadString(payload string) *Request {
	return r.SetPayload(strings.NewReader(payload))
}
