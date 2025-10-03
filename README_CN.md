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

> 本项目处于持续迭代阶段，功能和特性将不断完善。

`gs-http-gen` 是一款 **基于 IDL（接口定义语言）的 HTTP 代码生成工具**，
可根据统一的接口描述自动生成 **Go 语言** 服务端与 **其他语言** 客户端代码，服务端代码包括：

* 数据模型
* 验证逻辑
* HTTP 路由绑定
* 普通与流式（SSE）接口

通过声明式的 IDL 描述，开发者可以更专注于业务逻辑，显著减少样板代码编写和手动出错的风险。

此外，IDL 还可以作为 **跨团队、前后端统一的接口契约与文档**，帮助开发团队减少沟通成本，提升协作效率。

## 功能特性

### 🌟 IDL 驱动

* 使用简洁的接口定义语言描述服务接口与数据模型
* 支持：

    * 常量、枚举、结构体、`oneof` 类型
    * 泛型与类型嵌入（字段复用）
    * RPC 接口定义
    * 自定义注解（如 `json`、`go.type`、`enum_as_string` 等）

### ⚙️ 自动代码生成

根据 IDL 文件自动生成 Go 语言服务端及其他语言客户端代码：

* 数据模型结构体
* 参数与数据验证逻辑
* HTTP 请求参数绑定（路径、查询、头部、请求体）
* 普通与流式（SSE）接口实现
* 服务端接口定义与路由绑定
* 客户端调用代码

### 📦 丰富的数据类型支持

* 基本类型：`bool`、`int`、`float`、`string`
* 高级类型：`list`、`map`、`oneof`
* 可空字段：支持使用 `?` 表示可空字段
* 类型重定义与泛型模板

### 🔎 高效数据验证

* 无反射实现，高性能
* 支持基于表达式的验证规则
* 枚举类型自动生成 `OneOfXXX` 验证函数
* 支持自定义验证函数

### 🌐 HTTP 友好

* 自动绑定 HTTP 请求参数（路径、查询、头部、请求体）
* 支持 `form`、`json`、`multipart-form` 等格式
* 原生支持流式 RPC（SSE）接口

### 📝 注释与文档

* 支持单行与多行注释
* 未来计划支持 Markdown 格式注释

## 安装

- **推荐方式：**

使用 [gs](https://github.com/go-spring/gs) 集成开发工具。

- 单独安装本工具：

```bash
go install github.com/go-spring/gs-http-gen@latest
```

## 使用方法

### 第一步：定义 IDL 文件

创建 `.idl` 文件描述服务接口和数据模型。

> **语法说明：**
>
> * 文档由零个或多个定义组成，以换行或分号分隔，以 EOF 结束。
> * 标识符由字母、数字、下划线组成，且不能以数字开头。
> * 使用 `?` 表示字段可空。

示例：

```idl
// 常量定义
const int MAX_AGE = 150 // years
const int MIN_AGE = 18  // years

// 枚举定义
enum ErrCode {
    ERR_OK = 0
    PARAM_ERROR = 1003
}

enum Department {
    ENGINEERING = 1
    MARKETING = 2
    SALES = 3
}

// 数据结构
type Manager {
    string id
    string name (validate="len($) > 0 && len($) <= 64")
    int? age (validate="$ >= MIN_AGE && $ <= MAX_AGE")
    Department dept (enum_as_string)
}

type Response<T> {
    ErrCode errno (validate="OneOfErrCode($)")
    string errmsg
    T data
}

// 请求与响应
type ManagerReq {
    string id (path="id")
}

type GetManagerResp Response<Manager?>

// 普通 RPC 接口
rpc GetManager(ManagerReq) GetManagerResp {
    method="GET"
    path="/managers/{id}"
    summary="根据ID获取管理员信息"
}

// 流式处理
type StreamReq {
    string ID (json="id")
}

type StreamResp {
    string id
    string data
    Payload payload
}

oneof Payload {
    string text_data
    int? numberData (json="number_data")
    bool boolean_data (json="")
}

// 流式 RPC 接口
rpc Stream(StreamReq) stream<StreamResp> {
    method="GET"
    path="/stream/{id}"
    summary="流式传输数据"
}
```

### 第二步：生成代码

使用命令行工具生成代码：

```bash
# 仅生成服务端代码（默认）
gs-http-gen --server --output ./generated --go_package myservice

# 同时生成服务端和客户端代码
gs-http-gen --server --client --output ./generated --go_package myservice
```

**参数说明：**

| 参数             | 说明                    | 默认值     |
|----------------|-----------------------|---------|
| `--server`     | 生成服务端代码（HTTP 处理与路由绑定） | 否       |
| `--client`     | 生成客户端代码（HTTP 调用封装）    | 否       |
| `--output`     | 输出目录                  | `·`     |
| `--go_package` | 生成的 Go 包名             | `proto` |
| `--language`   | 目标语言（目前仅支持 `go`）      | `go`    |

### 第三步：使用生成的代码

示例：

```go
// 实现服务接口
type MyManagerServer struct{}

func (m *MyManagerServer) GetManager(ctx context.Context, req *proto.ManagerReq) *proto.GetManagerResp {
    // 普通响应
    return &proto.GetManagerResp{
        Data: &proto.Manager{
            Id:   "1",
            Name: "Jim",
            Dept: proto.Department_ENGINEERING,
        },
    }
}

func (m *MyManagerServer) Stream(ctx context.Context, req *proto.StreamReq, resp chan<- *proto.StreamResp) {
    // 流式响应
    for i := 0; i < 5; i++ {
        resp <- &proto.StreamResp{
            Id: strconv.Itoa(i),
            Payload: proto.Payload{
                TextData: "data",
            },
        }
    }
}

// 注册路由
mux := http.NewServeMux()
proto.InitRouter(mux, &MyManagerServer{})

http.ListenAndServe(":8080", mux)
```

## ⚠️ 注意事项

* 生成的代码不会自动强制字段必填，需在业务逻辑中自行保证。
* 不自动调用验证逻辑 `Validate()`，如需深度校验可自行组合。
* 建议统一管理生成的代码并保持与 IDL 一致，避免手动修改导致差异。

## 许可证

本项目采用 [Apache License 2.0](LICENSE) 许可证。
