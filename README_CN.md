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

## 项目简介

gs-http-gen 是一个基于 IDL（接口定义语言）的 HTTP 代码生成工具，专门用于 **Go-Spring** 框架。
它可以根据定义的接口描述自动生成 HTTP 服务端和客户端代码，包括数据模型、验证逻辑、路由绑定等。

该工具旨在简化 Go 语言 Web 服务的开发流程，通过声明式的 IDL 定义自动生成样板代码，
提高开发效率并减少手动编码错误。

更重要的是，IDL 不仅仅用于代码生成，它还可以作为前后端、跨部门之间接口信息的统一约定和文档。
通过使用标准化的 IDL 描述语言，可以清晰地了解接口的请求参数、响应格式、验证规则等关键信息，
有效减少沟通成本，确保接口的一致性和正确性。

## 功能特性

- **IDL 驱动**: 使用简洁的接口定义语言定义服务接口和数据模型
- **自动代码生成**: 根据 IDL 文件自动生成 Go 语言代码，包括：
    - 数据模型结构体
    - 数据验证逻辑
    - HTTP 路由绑定
    - 服务端接口定义
    - 客户端调用代码
- **多种数据类型支持**: 支持基本类型、结构体、枚举、列表、可选类型等
- **数据验证**: 内置丰富的数据验证规则，支持自定义验证
- **HTTP 参数绑定**: 自动将 HTTP 请求参数（路径、查询、头部、请求体）绑定到数据模型
- **类型嵌入**: 支持类型继承和字段复用，减少重复定义
- **灵活配置**: 支持生成服务端代码、客户端代码或两者兼有
- **枚举支持**: 支持枚举类型，并可选择以字符串形式序列化
- **流式处理**: 支持流式 RPC 接口生成
- **注释支持**: 支持在 IDL 中添加 Markdown 格式注释（暂未实现）

## 工具安装

**推荐**使用 gs 集成开发工具，参考 [https://github.com/go-spring/gs](https://github.com/go-spring/gs)。

单独安装本工具可以使用如下命令：

```bash
go install github.com/go-spring/gs-http-gen@latest
```

## 使用方法

### 第一步：定义IDL文件

首先创建IDL文件定义服务接口和数据模型：

```idl
// 定义常量
const int MAX_AGE = 150
const int MIN_AGE = 18

// 定义枚举
enum ErrCode {
    ERR_OK = 0
    PARAM_ERROR = 1003
}

enum Department {
    ENGINEERING = 1
    MARKETING = 2
    SALES = 3
}

// 定义数据结构
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

// 定义请求和响应
type ManagerReq {
    string id (path="id")
}

type GetManagerResp Response<Manager?>

// 定义流式处理
type StreamReq {
    string id
}

type StreamResp {
    string id
    string data
}

// 定义RPC接口
rpc GetManager(ManagerReq) GetManagerResp {
    method="GET"
    path="/managers/{id}"
    summary="根据ID获取管理员信息"
}

// 流式处理接口示例
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
gs-http-gen --server --output ./generated --package myservice

# 同时生成服务端和客户端代码
gs-http-gen --server --client --output ./generated --package myservice
```

命令行参数说明：

- `--server`: 生成服务端代码（HTTP处理、路由绑定等）
- `--client`: 生成客户端代码（HTTP调用封装）
- `--output`: 生成代码的输出目录（默认为当前目录）
- `--package`: Go包名（默认为"proto"）
- `--language`: 目标语言（目前仅支持"go"）

### 第三步：使用生成的代码

生成的代码包括数据模型、验证逻辑和HTTP处理逻辑：

```go
// 实现服务接口
type MyManagerServer struct{}

func (m *MyManagerServer) GetManager(ctx context.Context, req *proto.ManagerReq) *proto.GetManagerResp {
    // 实现业务逻辑
    data := proto.NewManager()
    data.SetName("Jim")
    res := proto.NewGetManagerResp()
    res.SetData(data)
    return res
}

func (m *MyManagerServer) Stream(ctx context.Context, req *proto.StreamReq, resp chan<- *proto.StreamResp) {
    // 实现流式处理逻辑
    for i := 0; i < 5; i++ {
        resp <- &proto.StreamResp{
            Id: strconv.Itoa(i),
        }
    }
}

// 注册路由
mux := http.NewServeMux()
proto.InitRouter(mux, &MyManagerServer{})
```

## 许可证

本项目采用Apache License 2.0许可证，详见[LICENSE](LICENSE)文件。