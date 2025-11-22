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

package httputil_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/go-spring/gs-http-gen/lib/httputil"
	"github.com/lvan100/golib/testing/assert"
)

func init() {
	httputil.DefaultClient = &LogHTTPClient{
		HTTPClient: httputil.DefaultClient,
	}
}

// ToString converts the given value to a string.
func ptr[T any](v T) *T {
	return &v
}

// LogHTTPClient is a HTTPClient implementation that logs all requests and responses.
type LogHTTPClient struct {
	httputil.HTTPClient
}

// JSON executes the given HTTP request using the provided Client.
func (c *LogHTTPClient) JSON(req *http.Request, meta httputil.RequestMeta) (*http.Response, []byte, error) {
	fmt.Printf("%#v\n", meta)
	return c.HTTPClient.JSON(req, meta)
}

// Stream executes the given HTTP request using the provided Client.
func (c *LogHTTPClient) Stream(req *http.Request, meta httputil.RequestMeta, s httputil.IStream) (*http.Response, error) {
	fmt.Printf("%#v\n", meta)
	return c.HTTPClient.Stream(req, meta, s)
}

type HelloClient struct {
	ServiceName string
}

type Item struct {
	ID int64 `json:"id"`
}

type Object struct {
	Item *Item  `json:"item"`
	Text string `json:"text"`
}

type HelloRequest struct {
	HelloRequestBody
	Int             int               `json:"int" query:"int"`
	String          string            `json:"string" query:"string"`
	IntPtr          *int              `json:"int_ptr" query:"int_ptr"`
	StringPtr       *string           `json:"string_ptr" query:"string_ptr"`
	IntSlice        []int             `json:"int_slice" query:"int_slice"`
	StringSlice     []string          `json:"string_slice" query:"string_slice"`
	ByteSlice       []byte            `json:"byte_slice" query:"byte_slice"`
	Object          *Object           `json:"object" query:"object"`
	ObjectSlice     []Object          `json:"object_slice" query:"object_slice"`
	IntStringMap    map[int]string    `json:"int_string_map" query:"int_string_map"`
	StringObjectMap map[string]Object `json:"string_object_map" query:"string_object_map"`
}

func (req *HelloRequest) QueryString() (string, error) {
	m := url.Values{}

	// Encode scalar values using the key format (e.g. a=1)
	m.Add("int", strconv.FormatInt(int64(req.Int), 10))
	m.Add("string", req.String)
	if req.IntPtr != nil {
		m.Add("int_ptr", strconv.FormatInt(int64(*req.IntPtr), 10))
	}
	if req.StringPtr != nil {
		m.Add("string_ptr", *req.StringPtr)
	}

	// Encode arrays using the repeated key format (e.g. a=1&a=2)
	for _, v := range req.IntSlice {
		m.Add("int_slice", strconv.FormatInt(int64(v), 10))
	}
	// Encode arrays using the repeated key format (e.g. a=1&a=2)
	for _, v := range req.StringSlice {
		m.Add("string_slice", v)
	}

	// Encode byte slices using base64 encoding (e.g., a=YWJj)
	if req.ByteSlice != nil {
		m.Add("byte_slice", base64.StdEncoding.EncodeToString(req.ByteSlice))
	}

	// Encode an array of objects using repeated keys with JSON values
	// e.g. items={"id":1,"name":"A"}&items={"id":2,"name":"B"}
	for _, v := range req.ObjectSlice {
		b, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		m.Add("object_slice", string(b))
	}

	// Encode maps or structs as JSON strings (e.g. data={"id":1,"name":"Alice"})
	if req.Object != nil {
		b, err := json.Marshal(req.Object)
		if err != nil {
			return "", err
		}
		m.Add("object", string(b))
	}

	// Encode maps or structs as JSON strings (e.g. data={"id":1,"name":"Alice"})
	if req.StringObjectMap != nil {
		b, err := json.Marshal(req.StringObjectMap)
		if err != nil {
			return "", err
		}
		m.Add("string_object_map", string(b))
	}

	// Encode maps or structs as JSON strings (e.g. data={"id":1,"name":"Alice"})
	if req.IntStringMap != nil {
		b, err := json.Marshal(req.IntStringMap)
		if err != nil {
			return "", err
		}
		m.Add("int_string_map", string(b))
	}

	return m.Encode(), nil
}

type HelloRequestBody struct{}

type HelloResponse struct {
	Message string `json:"message"`
}

// Hello sends a GET request to the /v1/hello endpoint with the given request body.
func (c *HelloClient) Hello(ctx context.Context, req *HelloRequest, opts ...httputil.RequestOption) (*http.Response, *HelloResponse, error) {

	path := "/v1/hello"
	if s, err := req.QueryString(); err != nil {
		return nil, nil, err
	} else if s != "" {
		path += "?" + s
	}

	r, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Accept", "application/json")

	opts = append(opts, httputil.WithTarget(c.ServiceName))
	opts = append(opts, httputil.WithPath("/v1/hello"))
	opts = append(opts, httputil.WithSchema("http"))
	return httputil.JSONResponse[HelloResponse](r, opts...)
}

type StreamRequest struct {
	StreamRequestBody
	Int             int               `json:"int" query:"int"`
	String          string            `json:"string" query:"string"`
	IntPtr          *int              `json:"int_ptr" query:"int_ptr"`
	StringPtr       *string           `json:"string_ptr" query:"string_ptr"`
	IntSlice        []int             `json:"int_slice" query:"int_slice"`
	StringSlice     []string          `json:"string_slice" query:"string_slice"`
	ByteSlice       []byte            `json:"byte_slice" query:"byte_slice"`
	Object          *Object           `json:"object" query:"object"`
	ObjectSlice     []Object          `json:"object_slice" query:"object_slice"`
	IntStringMap    map[int]string    `json:"int_string_map" query:"int_string_map"`
	StringObjectMap map[string]Object `json:"string_object_map" query:"string_object_map"`
}

func (req *StreamRequest) QueryString() (string, error) {
	m := url.Values{}

	// Encode scalar values using the key format (e.g. a=1).
	m.Add("int", strconv.FormatInt(int64(req.Int), 10))
	m.Add("string", req.String)
	// Nil is omitted and not transmitted.
	if req.IntPtr != nil {
		m.Add("int_ptr", strconv.FormatInt(int64(*req.IntPtr), 10))
	}
	// Nil is omitted and not transmitted.
	if req.StringPtr != nil {
		m.Add("string_ptr", *req.StringPtr)
	}

	// Encode slices using JSON array format (e.g., a=[1,2]),
	// nil is omitted and not transmitted.
	if req.IntSlice != nil {
		b, err := json.Marshal(req.IntSlice)
		if err != nil {
			return "", err
		}
		m.Add("int_slice", string(b))
	}

	// Encode slices using JSON array format (e.g., a=[1,2]),
	// allowing nil to be encoded as null.
	{
		b, err := json.Marshal(req.StringSlice)
		if err != nil {
			return "", err
		}
		m.Add("string_slice", string(b))
	}

	// Encode byte slices using base64 encoding (e.g., a=YWJj).
	if req.ByteSlice != nil {
		m.Add("byte_slice", base64.StdEncoding.EncodeToString(req.ByteSlice))
	}

	// Encode an array of objects using repeated keys with JSON values
	// e.g. items={"id":1,"name":"A"}&items={"id":2,"name":"B"},
	// nil is omitted and not transmitted.
	for _, v := range req.ObjectSlice {
		b, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		m.Add("object_slice", string(b))
	}

	// Encode maps or structs as JSON strings (e.g. data={"id":1,"name":"Alice"}),
	// nil is omitted and not transmitted.
	if req.Object != nil {
		b, err := json.Marshal(req.Object)
		if err != nil {
			return "", err
		}
		m.Add("object", string(b))
	}

	// Encode maps or structs as JSON strings (e.g. data={"id":1,"name":"Alice"}),
	// nil is omitted and not transmitted.
	if req.StringObjectMap != nil {
		b, err := json.Marshal(req.StringObjectMap)
		if err != nil {
			return "", err
		}
		m.Add("string_object_map", string(b))
	}

	// Encode maps or structs as JSON strings (e.g. data={"id":1,"name":"Alice"}),
	// allowing nil to be encoded as null.
	{
		b, err := json.Marshal(req.IntStringMap)
		if err != nil {
			return "", err
		}
		m.Add("int_string_map", string(b))
	}

	return m.Encode(), nil
}

type StreamRequestBody struct {
	Prompt string `json:"prompt"`
}

type StreamResponse struct {
	Text string `json:"text"`
}

// Stream sends a POST request to the /v1/stream endpoint with the given request body.
func (c *HelloClient) Stream(ctx context.Context, req *StreamRequest, opts ...httputil.RequestOption) (
	*http.Response, *httputil.Stream, error) {

	path := "/v1/stream"
	if s, err := req.QueryString(); err != nil {
		return nil, nil, err
	} else if s != "" {
		path += "?" + s
	}

	// Encode the request body using JSON.
	body, err := json.Marshal(req)
	if err != nil {
		return nil, nil, err
	}
	buf := bytes.NewBuffer(body)

	r, err := http.NewRequestWithContext(ctx, "POST", path, buf)
	if err != nil {
		return nil, nil, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "text/event-stream")

	opts = append(opts, httputil.WithTarget(c.ServiceName))
	opts = append(opts, httputil.WithPath("/v1/hello"))
	opts = append(opts, httputil.WithSchema("http"))
	return httputil.StreamResponse(r, opts...)
}

func TestHello(t *testing.T) {
	server := &http.Server{Addr: ":9090", Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.RawQuery)
		_ = r.Header.Write(os.Stdout)
		fmt.Println()
		_, _ = w.Write([]byte(fmt.Sprintf(`{"message": "hello %s"}`, r.URL.Query().Get("string"))))
	})}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = server.ListenAndServe()
	}()
	time.Sleep(time.Millisecond * 100)

	h := http.Header{}
	h.Set("X-Request-ID", "12345678")

	client := &HelloClient{
		ServiceName: "127.0.0.1:9090",
	}

	_, data, err := client.Hello(context.Background(), &HelloRequest{
		Int:         5,
		String:      "world",
		IntPtr:      ptr(10),
		StringPtr:   ptr("message"),
		IntSlice:    []int{1, 2, 3},
		StringSlice: []string{"a", "b", "c"},
		ByteSlice:   []byte("hello world"),
		Object: &Object{
			Item: &Item{ID: 1010},
			Text: "message",
		},
		ObjectSlice: []Object{
			{
				Item: &Item{ID: 1010},
				Text: "message",
			},
			{
				Item: &Item{ID: 1010},
				Text: "message",
			},
		},
		IntStringMap: map[int]string{1: "one", 2: "two"},
		StringObjectMap: map[string]Object{
			"one": {
				Item: &Item{ID: 1010},
				Text: "message",
			},
			"two": {
				Item: &Item{ID: 1010},
				Text: "message",
			},
		},
	}, httputil.WithHeader(h))
	assert.Error(t, err).Nil()
	assert.That(t, data).Equal(&HelloResponse{Message: "hello world"})

	_ = server.Shutdown(context.Background())
	wg.Wait()
}

func TestStream(t *testing.T) {
	server := http.Server{Addr: ":9090", Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.RawQuery)
		_ = r.Header.Write(os.Stdout)
		fmt.Println()
		for i := range 5 {
			_, _ = w.Write([]byte(fmt.Sprintf("id: %v\n", i)))
			_, _ = w.Write([]byte(fmt.Sprintf("data: %s\n", `{"message": "hello world"}`)))
			_, _ = w.Write([]byte("\n\n"))
			w.(http.Flusher).Flush()
			time.Sleep(time.Millisecond * 500)
		}
		fmt.Println()
		fmt.Println("server done")
	})}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = server.ListenAndServe()
	}()
	time.Sleep(time.Millisecond * 100)

	h := http.Header{}
	h.Set("X-Request-ID", "12345678")

	client := &HelloClient{
		ServiceName: "127.0.0.1:9090",
	}

	ctx := context.Background()
	_, resp, err := client.Stream(ctx, &StreamRequest{
		StreamRequestBody: StreamRequestBody{
			Prompt: "hello world",
		},
		Int:       5,
		String:    "world",
		IntPtr:    ptr(10),
		StringPtr: ptr("message"),
		IntSlice:  []int{1, 2, 3},
		ByteSlice: []byte("hello world"),
		Object: &Object{
			Item: &Item{ID: 1010},
			Text: "message",
		},
		ObjectSlice: []Object{
			{
				Item: &Item{ID: 1010},
				Text: "message",
			},
			{
				Item: &Item{ID: 1010},
				Text: "message",
			},
		},
		StringObjectMap: map[string]Object{
			"one": {
				Item: &Item{ID: 1010},
				Text: "message",
			},
			"two": {
				Item: &Item{ID: 1010},
				Text: "message",
			},
		},
	}, httputil.WithHeader(h))
	defer resp.Close()
	assert.Error(t, err).Nil()

	for resp.Next(ctx, time.Second) {
		fmt.Println(resp.Event())
		//resp.Close()
	}
	fmt.Println(resp.Error())
	fmt.Println("done")
}
