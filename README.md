# gs-http-gen

<div>
   <img src="https://img.shields.io/github/license/go-spring/gs-http-gen" alt="license"/>
   <img src="https://img.shields.io/github/go-mod/go-version/go-spring/gs-http-gen" alt="go-version"/>
   <img src="https://img.shields.io/github/v/release/go-spring/gs-http-gen?include_prereleases" alt="release"/>
   <a href="https://codecov.io/gh/go-spring/gs-http-gen" > 
      <img src="https://codecov.io/gh/go-spring/gs-http-gen/graph/badge.svg?token=SX7CV1T0O8" alt="test-coverage"/>
   </a>
   <a href="https://deepwiki.com/go-spring/gs-http-gen"><img src="https://deepwiki.com/badge.svg" alt="Ask DeepWiki"></a>
</div>

[English](README.md) | [中文](README_CN.md)

## Project Overview

**gs-http-gen** is an HTTP code generation tool based on **IDL (Interface Definition Language)**, designed for the *
*Go-Spring** framework. It can automatically generate HTTP server and client code from interface definitions, including
data models, validation logic, and route bindings.

The tool aims to simplify the development workflow of Go web services. By using declarative IDL definitions, it
generates boilerplate code automatically, improving development efficiency and reducing human errors.

More importantly, IDL is not only used for code generation—it also serves as a unified contract and documentation for
APIs across frontend-backend teams and departments. With standardized IDL definitions, key details such as request
parameters, response formats, and validation rules become clear, reducing communication costs and ensuring API
consistency and correctness.

## Features

- **IDL-driven**: Define service interfaces and data models using a simple interface definition language.
- **Automatic code generation**: Generate Go code from IDL files, including:
  - Data model structs
    - Data validation logic
    - HTTP route bindings
    - Server interface definitions
    - Client call code
- **Rich type support**: Supports basic types, structs, enums, lists, optional types, etc.
- **Data validation**: Built-in validation rules with support for custom validators.
- **HTTP parameter binding**: Automatically bind HTTP request parameters (path, query, header, body) to data models.
- **Type embedding**: Supports type inheritance and field reuse to reduce redundancy.
- **Flexible configuration**: Generate server code, client code, or both.
- **Enum support**: Enum types with optional string serialization.
- **Streaming support**: Generate streaming RPC interfaces.
- **Annotation support**: Add Markdown-style comments in IDL (not yet implemented).

## Installation

**Recommended**: Install via the **gs** integrated development tool,
see [https://github.com/go-spring/gs](https://github.com/go-spring/gs).

To install this tool separately:

```bash
go install github.com/go-spring/gs-http-gen@latest
```

## Usage

### Step 1: Define an IDL file

First, create an IDL file to define service interfaces and data models:

```idl
// Define constants
const int MAX_AGE = 150
const int MIN_AGE = 18

// Define enums
enum ErrCode {
    ERR_OK = 0
    PARAM_ERROR = 1003
}

enum Department {
    ENGINEERING = 1
    MARKETING = 2
    SALES = 3
}

// Define data structures
type Manager {
    string id
    string name (validate="len($) > 0 && len($) <= 64")
    int? age (validate="$ >= MIN_AGE && $ <= MAX_AGE")
    Department dept
}

type Response<T> {
    ErrCode errno = ErrCode.ERR_OK (validate="OneOfErrCode($)")
    string errmsg
    T data
}

// Define request and response
type ManagerReq {
    string id (path="id")
}

type GetManagerResp Response<Manager?>

// Define streaming types
type StreamReq {
    string id
}

type StreamResp {
    string id
    string data
}

// Define RPC interface
rpc GetManager(ManagerReq) GetManagerResp {
    method="GET"
    path="/managers/{id}"
    summary="Get manager info by ID"
}

// Example of streaming interface
rpc Stream(StreamReq) stream<StreamResp> {
    method="GET"
    path="/stream/{id}"
    summary="Stream data transfer"
}
```

### Step 2: Generate code

Use the CLI tool to generate code:

```bash
# Generate server code only (default)
gs-http-gen --server --output ./generated --package myservice

# Generate both server and client code
gs-http-gen --server --client --output ./generated --package myservice
```

Command-line options:

* `--server`: Generate server code (HTTP handlers, route bindings, etc.)
* `--client`: Generate client code (HTTP call wrappers)
* `--output`: Output directory for generated code (default: current directory)
* `--package`: Go package name (default: "proto")
* `--language`: Target language (currently only `"go"` supported)

### Step 3: Use the generated code

The generated code includes data models, validation logic, and HTTP handlers:

```go
// Implement service interface
type MyManagerServer struct{}

func (m *MyManagerServer) GetManager(ctx context.Context, req *proto.ManagerReq) *proto.GetManagerResp {
    // Business logic
    data := proto.NewManager()
    data.SetName("Jim")
    res := proto.NewGetManagerResp()
    res.SetData(data)
    return res
}

func (m *MyManagerServer) Stream(ctx context.Context, req *proto.StreamReq, resp chan<- *proto.StreamResp) {
    // Streaming logic
    for i := 0; i < 5; i++ {
        resp <- &proto.StreamResp{
            Id: strconv.Itoa(i),
        }
    }
}

// Register routes
mux := http.NewServeMux()
proto.InitRouter(mux, &MyManagerServer{})
```

## License

This project is licensed under the **Apache License 2.0**. See the [LICENSE](LICENSE) file for details.