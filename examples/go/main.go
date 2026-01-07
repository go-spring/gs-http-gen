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

package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"examples/ginsvr"
	"examples/proto"
	"examples/server"

	"github.com/go-spring/stdlib/httpsvr"
)

// init sets the working directory of the program to the directory
// where this source file resides. This ensures that relative paths
// used later in the program (e.g., for output) are resolved correctly.
func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot determine caller directory")
	}
	execDir := filepath.Dir(filename)
	err := os.Chdir(execDir)
	if err != nil {
		panic(err)
	}
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println("working directory:", workDir)
}

func main() {
	TestManager()
	//TestStream()
	//TestStreamV2()
}

func TestManager() {
	svr := ginsvr.NewGinServer(":9191")
	for _, r := range proto.Routers(&server.ManagerServer{}, ginsvr.NewGinRequestContext) {
		svr.HandleFunc(r)
	}
	go func() {
		fmt.Println(svr.ListenAndServe())
	}()
	time.Sleep(time.Millisecond * 300)

	resp, err := http.Get("http://localhost:9191/managers/123")
	if err != nil {
		panic(err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	fmt.Println(string(b))
	svr.Shutdown(context.Background())
}

func TestStream() {
	svr := httpsvr.NewSimpleServer(":9191")
	for _, r := range proto.Routers(&server.ManagerServer{}, ginsvr.NewGinRequestContext) {
		svr.Route(r)
	}
	go func() {
		fmt.Println(svr.ListenAndServe())
	}()
	time.Sleep(time.Millisecond * 300)

	resp, err := http.Get("http://localhost:9191/assistant")
	if err != nil {
		panic(err)
	}
	b, err := io.ReadAll(io.LimitReader(resp.Body, 1025))
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	fmt.Print(string(b))
	svr.Shutdown(context.Background())
}

func TestStreamV2() {
	svr := httpsvr.NewSimpleServer(":9191")
	for _, r := range proto.Routers(&server.ManagerServer{}, ginsvr.NewGinRequestContext) {
		svr.Route(r)
	}
	go func() {
		fmt.Println(svr.ListenAndServe())
	}()
	time.Sleep(time.Millisecond * 300)

	resp, err := http.Get("http://localhost:9191/assistantV2")
	if err != nil {
		panic(err)
	}
	b, err := io.ReadAll(io.LimitReader(resp.Body, 1025))
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	fmt.Print(string(b))
	svr.Shutdown(context.Background())
}

/*

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
		b, err := jsonflow.Marshal(v)
		if err != nil {
			return "", err
		}
		m.Add("object_slice", string(b))
	}

	// Encode maps or structs as JSON strings (e.g. data={"id":1,"name":"Alice"})
	if req.Object != nil {
		b, err := jsonflow.Marshal(req.Object)
		if err != nil {
			return "", err
		}
		m.Add("object", string(b))
	}

	// Encode maps or structs as JSON strings (e.g. data={"id":1,"name":"Alice"})
	if req.StringObjectMap != nil {
		b, err := jsonflow.Marshal(req.StringObjectMap)
		if err != nil {
			return "", err
		}
		m.Add("string_object_map", string(b))
	}

	// Encode maps or structs as JSON strings (e.g. data={"id":1,"name":"Alice"})
	if req.IntStringMap != nil {
		b, err := jsonflow.Marshal(req.IntStringMap)
		if err != nil {
			return "", err
		}
		m.Add("int_string_map", string(b))
	}

	return m.Encode(), nil
}

type HelloRequestBody struct{}

*/
