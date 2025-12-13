/*
 * Copyright 2025 The Go-Spring Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package httputil

import (
	"bytes"
	"encoding/json"
	"io"
	"maps"
	"net/http"
)

// Ptr returns a pointer to the given value.
func Ptr[T any](v T) *T {
	return &v
}

///////////////////////////////// interface ///////////////////////////////////

// RequestMeta holds contextual information for an HTTP request.
type RequestMeta struct {
	Target string
	Schema string
	Path   string
	Header http.Header
	Config map[string]string
}

// HTTPClient defines a customizable HTTP executor interface.
// Implementing this interface allows users to provide their own
// HTTP execution logic (for example, to add retry, logging, or tracing).
type HTTPClient interface {
	JSON(req *http.Request, meta RequestMeta) (*http.Response, []byte, error)
}

// DefaultClient is the default HTTPClient implementation.
var DefaultClient HTTPClient = &SimpleHTTPClient{http.DefaultClient}

// SimpleHTTPClient is the default implementation of HTTPClient,
// which delegates to the standard library http.Client.
// Target must be a static IP:port or domain name.
type SimpleHTTPClient struct {
	Client *http.Client
}

// do executes the given HTTP request using the embedded http.Client.
func (c *SimpleHTTPClient) do(r *http.Request, meta RequestMeta) (*http.Response, error) {
	r.Host = meta.Target
	r.URL.Host = meta.Target
	r.URL.Scheme = meta.Schema
	for k, values := range meta.Header {
		for _, v := range values {
			r.Header.Add(k, v)
		}
	}
	return c.Client.Do(r)
}

// JSON executes the HTTP request using the embedded http.Client.
// It reads the entire response body into memory and returns both
// the *http.Response and the body as a byte slice.
//
// The response body is also replaced with a reusable buffer so that
// it can be read again by the caller if needed.
//
// Note: For very large responses, this may be memory intensive.
func (c *SimpleHTTPClient) JSON(r *http.Request, meta RequestMeta) (*http.Response, []byte, error) {
	resp, err := c.do(r, meta)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	// Reset the response body to allow it to be read again later.
	resp.Body = io.NopCloser(bytes.NewBuffer(b))
	return resp, b, nil
}

////////////////////////////////// response ///////////////////////////////////

// RequestOption is a function type that modifies the RequestMeta.
type RequestOption func(info *RequestMeta)

// WithTarget sets the target to the RequestMeta.
func WithTarget(target string) RequestOption {
	return func(info *RequestMeta) {
		info.Target = target
	}
}

// WithSchema sets the schema to the RequestMeta.
func WithSchema(schema string) RequestOption {
	return func(info *RequestMeta) {
		info.Schema = schema
	}
}

// WithPath sets the path to the RequestMeta.
func WithPath(path string) RequestOption {
	return func(info *RequestMeta) {
		info.Path = path
	}
}

// WithHeader sets the given HTTP headers to the RequestMeta.
func WithHeader(header http.Header) RequestOption {
	return func(meta *RequestMeta) {
		if meta.Header == nil {
			meta.Header = http.Header{}
		}
		maps.Copy(meta.Header, header)
	}
}

// WithConfig sets the given configuration map to the RequestMeta.
func WithConfig(config map[string]string) RequestOption {
	return func(meta *RequestMeta) {
		if meta.Config == nil {
			meta.Config = map[string]string{}
		}
		maps.Copy(meta.Config, config)
	}
}

// buildMeta creates a RequestMeta with the given options.
func buildMeta(opts []RequestOption) RequestMeta {
	meta := RequestMeta{
		Header: http.Header{},
		Config: map[string]string{},
	}
	for _, opt := range opts {
		opt(&meta)
	}
	return meta
}

// JSONResponse executes the given HTTP request using the provided HTTPClient,
// reads the response body, and unmarshals it into a value of type RespType.
func JSONResponse[RespType any](r *http.Request, opts ...RequestOption) (*http.Response, RespType, error) {
	var ret RespType
	resp, b, err := DefaultClient.JSON(r, buildMeta(opts))
	if err != nil {
		return nil, ret, err
	}
	if err = json.Unmarshal(b, &ret); err != nil {
		return nil, ret, err
	}
	return resp, ret, nil
}
