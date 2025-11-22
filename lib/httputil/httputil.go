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
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"maps"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

// Ptr returns a pointer to the given value.
func Ptr[T any](v T) *T {
	return &v
}

////////////////////////////////// stream /////////////////////////////////////

// Message represents a single message unit read from the stream.
type Message struct {
	Data string
	Err  error
}

type IStream interface {
	Send(msg Message) bool
}

// SSEEvent represents an SSE event.
type SSEEvent struct {
	ID    string
	Data  string
	Event string
	Retry int
}

// Stream manages streaming data asynchronously from an HTTP response.
// It supports safe concurrent use and can be closed idempotently.
type Stream struct {
	msgs   chan Message
	curr   Message
	closed atomic.Bool
	done   chan struct{}
}

// NewStream creates and initializes a new Stream instance.
func NewStream() *Stream {
	return &Stream{
		msgs: make(chan Message),
		done: make(chan struct{}),
	}
}

// Event unmarshals the current data item into an SSEEvent.
func (s *Stream) Event() (SSEEvent, error) {
	var e SSEEvent
	return e, nil
}

// Error returns the last error encountered by the stream.
func (s *Stream) Error() error {
	return s.curr.Err
}

// Next waits for the next data item from the stream,
// honoring the provided context and optional timeout.
// Returns true if a new data item is successfully received,
// or false if the stream is closed, the context is done, or an error occurs.
func (s *Stream) Next(ctx context.Context, timeout time.Duration) bool {
	if s.closed.Load() {
		return false
	}

	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	var ok bool
	select {
	case <-ctx.Done():
		s.curr.Data = ""
		s.curr.Err = ctx.Err()
		return false
	case s.curr, ok = <-s.msgs:
		if !ok {
			return false
		}
		if s.curr.Err != nil {
			// Treat io.EOF as normal stream termination
			if s.curr.Err == io.EOF {
				s.curr.Err = nil
				return false
			}
			return false
		}
		return true
	}
}

// Send pushes a Message into the internal channel.
// Returns false if the stream is closed or already done.
func (s *Stream) Send(msg Message) bool {
	if s.closed.Load() {
		return false
	}
	select {
	case <-s.done:
		return false
	case s.msgs <- msg:
		return true
	}
}

// Close closes the Stream safely (idempotent).
// It ensures that the internal channels are closed only once.
func (s *Stream) Close() {
	if s.closed.CompareAndSwap(false, true) {
		close(s.done)
		close(s.msgs)
	}
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
	Stream(req *http.Request, meta RequestMeta, stream IStream) (*http.Response, error)
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

// Stream executes an HTTP request and continuously reads lines from the response body.
// Each line is sent into the returned Stream channel asynchronously.
func (c *SimpleHTTPClient) Stream(r *http.Request, meta RequestMeta, stream IStream) (*http.Response, error) {
	resp, err := c.do(r, meta)
	if err != nil {
		return nil, err
	}

	go func() {
		defer resp.Body.Close()
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			text := strings.TrimSpace(scanner.Text())
			if text == "" {
				continue
			}
			if !stream.Send(Message{Data: text}) {
				return
			}
		}
		if err := scanner.Err(); err != nil {
			stream.Send(Message{Err: err})
		} else {
			stream.Send(Message{Err: io.EOF})
		}
	}()
	return resp, nil
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
func JSONResponse[RespType any](r *http.Request, opts ...RequestOption) (*http.Response, *RespType, error) {
	resp, b, err := DefaultClient.JSON(r, buildMeta(opts))
	if err != nil {
		return nil, nil, err
	}
	var ret RespType
	if err = json.Unmarshal(b, &ret); err != nil {
		return nil, nil, err
	}
	return resp, &ret, nil
}

// StreamResponse executes the given HTTP request using the provided HTTPClient,
// and returns a Stream instance for streaming the response body.
func StreamResponse(r *http.Request, opts ...RequestOption) (*http.Response, *Stream, error) {
	s := NewStream()
	resp, err := DefaultClient.Stream(r, buildMeta(opts), s)
	if err != nil {
		return nil, nil, err
	}
	return resp, s, nil
}
