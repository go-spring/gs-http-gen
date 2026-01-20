# gs-http-gen IDL 语法中文说明文档

## 概述

gs-http-gen 使用一种专门设计的接口定义语言（IDL）来描述数据模型和接口服务。
该语言支持定义常量、枚举、结构体、联合类型以及RPC服务。

## 基本语法规则

### 注释

- 单行注释：使用 `//` 或 `#` 开头
- 多行注释：使用 `/* ... */` 包围

### 关键字

以下关键字是保留的，不能用作标识符：

- `extends` - 错误码扩展
- `const` - 定义常量
- `enum` - 定义枚举
- `type` - 定义类型
- `oneof` - 定义联合类型
- `rpc` - 定义普通RPC方法
- `sse` - 定义流式RPC方法
- `true`/`false` - 布尔值
- `optional` - 可选字段标记
- `required` - 必需字段标记

### 基础类型

- `bool` - 布尔类型
- `int` - 整数类型（note：注意没有无符号整型）
- `float` - 浮点数类型
- `string` - 字符串类型（note：只能使用双引号，单引号不行）
- `bytes` - 字节数组类型（note：传输时使用 base64 编码成字符串）

### 容器类型

- `map<K, V>` - 映射类型，键类型只能是 `int` 或 `string`
- `list<T>` - 列表类型
- map 和 list 都支持嵌套，例如：`list<list<int>>` 等

## 语法详细说明

### 1. 常量定义

使用 `const` 关键字定义常量，常量的类型只能是 bool、int、float、string 四种：

```
const string APP_NAME = "MyApp"
const int MAX_SIZE = 100
const float PI = 3.14159
const bool DEBUG = true
```

语法格式，标识符名称要求大写字母开头：

```
const <基础类型> <标识符名称> = <常量值>
```

### 2. 枚举定义

使用 `enum` 关键字定义枚举类型。枚举的每个字段必须赋一个整数值。

基本定义格式：

```
enum <枚举名> {
    <字段名> = <整数值>
    <字段名> = <整数值>
    ...
}
```

示例：

```
enum Color {
    RED = 1
    GREEN = 2
    BLUE = 3
}
```

#### 2.1 错误码枚举

有一种特殊的 enum 类型，叫错误码。标识一个 enum 类型是错误码的关键是在字段后面添加 `errmsg` 注解，例如：

```
enum ErrCode {
    ERR_OK = 0 (errmsg="success")
    PARAM_ERROR = 1003 (errmsg="parameter error")
    NOT_FOUND = 404 (errmsg="resource not found")
}
```

错误码注解（errmsg）用于为错误码提供人类可读的消息描述，在生成代码时会创建相应的错误消息映射。

#### 2.2 错误码扩展

错误码可以通过 `extends` 关键字进行扩展，扩展后的错误码会与基础错误码合并：

```
// 基础错误码定义
enum ErrCode {
    ERR_OK = 0 (errmsg="success")
    PARAM_ERROR = 1003 (errmsg="parameter error")
}

// 扩展错误码定义
enum extends ErrCode {
    USER_NOT_FOUND = 404
    PERMISSION_DENIED = 403
}
```

在代码生成过程中，扩展的错误码会与基础错误码合并，形成一个完整的错误码集合。扩展错误码需要注意：

1. 扩展的错误码类型必须已经定义
2. 扩展的错误码不能与已有错误码重复（字段名和值都不能重复）
3. 扩展错误码的值也建议保持单调递增
4. 合并后，基础错误码和扩展错误码会一起生成对应的Go代码

#### 2.3 枚举字段注解

枚举字段可以使用注解来添加额外信息：

- `errmsg` - 为错误码提供描述信息，仅用于错误码枚举

示例：

```
enum Status {
    ACTIVE = 1 (errmsg="active status")
    INACTIVE = 0 (errmsg="inactive status")
}
```

### 3. 注解（Annotations）

字段和接口可以使用注解来添加元数据。注解语法支持多行定义：

```
type User {
    string name (json="name", go.type="string")
    int age (json="age,omitempty")
}
```

注解语法：todo(需要补充支持多行)

```
(注解名 = 注解值, 注解名 = 注解值, ...)
```

注解提供了灵活的元数据机制，用于控制生成代码的行为，可以指定序列化方式、验证规则、API路由等信息。

### 4. 类型定义

#### 4.1 结构体类型

用 type 关键字定义结构体，字段支持 optional（可选）和 required（必需）：

```
type User {
    required string name
    int age
    optional string email
    list<string> tags
}
```

#### 4.2 字段注解

在字段定义中可以使用注解来指定额外的元数据，常见的字段注解包括：

- `json` - 指定JSON序列化字段名和选项
- `form` - 指定表单绑定字段名
- `path` - 指定路径参数绑定
- `query` - 指定查询参数绑定
- `validate` - 指定验证表达式
- `deprecated` - 标记字段已弃用
- `enum_as_string` - 枚举作为字符串处理
- `compat_default` - 兼容性默认值

例如：

```
type User {
    string name (json="name", go.type="string")
    int age (json="age,omitempty")
    string email (validate="email", deprecated="true")
}
```

##### 4.2.1 Validate 注解

Validate 注解用于对字段值进行验证，支持多种内置验证函数和复杂的表达式。

内置验证函数包括：

- `len` - 验证长度（字符串、列表、映射等）
- `email` - 验证邮箱格式
- `url` - 验证URL格式
- 自定义验证函数

验证表达式支持的操作符：

- 比较操作符：`==`, `!=`, `<`, `<=`, `>`, `>=`
- 逻辑操作符：`&&`, `||`, `!`
- 算术操作符：`+`, `-`, `*`, `/`
- 函数调用：`len($)`、`email($)`

特殊变量 `$` 代表当前字段的值。

示例：

```
type User {
    required string name (validate="$ != '' && len($) <= 64")  // 非空且长度不超过64
    int age (validate="$ >= 0 && $ <= 150")                  // 年龄在0-150之间
    string email (validate="email($)")                        // 邮箱格式验证
    list<string> tags (validate="len($) <= 10")               // 标签数量不超过10个
}
```

如果字段是可选的（optional），验证仅在字段值存在时执行。

todo（其实字段上可以添加很多注解，也可以介绍一下）

#### 4.3 泛型类型

支持简单的泛型定义，一般用于对返回值进行定义的场景，这时候一般只需要对 data 进行泛型定义：

```
type Response<T> {
    int code
    string message
    T data
}
```

#### 4.4 泛型实例化

对泛型进行实例化，注意语法格式，只能对 user_type 泛型进行实例化：

```
type UserResponse Response<User>
```

### 5. 字段定义

#### 5.1 普通字段

```
<可选修饰符> <类型> <字段名>
```

修饰符：

- `required` - 表示必需字段
- `optional` - 表示可选字段（默认）

#### 5.2 嵌入类型

可以在结构体中嵌入其他类型，代码生成的时候是把嵌入的类型进行展开：

```
type Address {
    string street
    string city
}

type Person {
    Address  // 嵌入 Address 类型
    string name
}
```

### 6. 联合类型定义

使用 `oneof` 定义联合类型。联合类型在内部实现上会生成一个特殊结构，包含一个类型标识字段和多个可能的数据字段。
在序列化时，只有被选择的那个字段会被包含在输出中，同时通过类型标识字段指示当前激活的选项：

```
oneof Value {
    User
    Manager
}
```

生成的结构通常类似：

```
type Value struct {
    FieldType ValueType `json:"FieldType" form:"FieldType"`  // 类型标识字段
    User      *User     `json:"User,omitempty" form:"User"`  // 可能的选项1
    Manager   *Manager  `json:"Manager,omitempty" form:"Manager"`  // 可能的选项2
}
```

其中 FieldType 是一个枚举类型，用于标识当前值的类型。

### 7. RPC 接口定义

定义服务接口：

```
rpc GetUser (GetUserRequest) GetUserResponse {
    method = "GET"
    path = "/user/:id"
}
```

对于流式接口，使用 `sse` 关键字：

```
sse StreamEvents (StreamRequest) Event {
    method = "GET"
    path = "/events"
}
```

#### 7.1 RPC 注解

RPC接口可以使用注解来指定HTTP方法、路径、超时时间等路由信息，常见的RPC注解包括：

- `method` - HTTP方法
- `path` - 请求路径
- `content-type` - 内容类型
- `connTimeout` - 连接超时
- `readTimeout` - 读取超时
- `writeTimeout` - 写入超时

例如：

```
rpc GetUser (GetUserRequest) GetUserResponse {
    method = "GET"
    path = "/user/:id"
    connTimeout = "5s"
    readTimeout = "10s"
}
```

#### 7.2 RESTful Path 参数

RESTful路径参数用于定义动态路径，支持两种风格：

1. 冒号风格（Colon Style）：`:paramName` 或 `:paramName*`
2. 大括号风格（Brace Style）：`{paramName}` 或 `{paramName...}`

路径参数类型：

- 普通参数：`:id` 或 `{id}` - 匹配单个路径段
- 通配符参数：`:path*` 或 `{path...}` - 匹配多个路径段

示例：

```
// 冒号风格
rpc GetUser (GetUserRequest) GetUserResponse {
    method = "GET"
    path = "/user/:id"  // 路径参数 :id
}

rpc GetFile (GetFileRequest) GetFileResponse {
    method = "GET"
    path = "/files/:path*"  // 通配符参数 :path*
}

// 大括号风格
rpc GetUser (GetUserRequest) GetUserResponse {
    method = "GET"
    path = "/user/{id}"  // 路径参数 {id}
}

rpc GetFile (GetFileRequest) GetFileResponse {
    method = "GET"
    path = "/files/{path...}"  // 通配符参数 {path...}
}

// 复杂路径示例
rpc ComplexPath (ComplexPathRequest) ComplexPathResponse {
    method = "GET"
    path = "/org/{orgId}/repos/:repoId/branches/{branch...}"  // 混合参数
}
```

在请求类型中使用 path 绑定注解将字段与路径参数关联：

```
type GetUserRequest {
    string userId (path="id")  // 将 userId 字段绑定到路径参数 :id 或 {id}
    string locale (query="locale")  // 查询参数示例
}
```

注意：

- 路径参数对应的字段必须是必需的（required）
- 支持的参数类型包括：int、string 等基本类型
- 参数名必须与路径中定义的参数名匹配

todo（缺少很多注解，还是要参考语法和golang实现）

### 8. 常量值类型

支持以下类型的常量值，todo（这里应该说明是注解中的常量值吧）：

- 整数字面量：`42`, `-17`, `0x1A2B`
- 浮点数字面量：`3.14`, `.5`, `-2.7e10`
- 字符串字面量：`"hello"`, `"escaped \" quote"`
- 布尔值：`true`, `false`
- 标识符：通常用于引用枚举成员

### 9. 标识符规则 todo（9和10是不是应该放在前面部分）

- 以字母开头
- 可以包含字母、数字、下划线 `_` 和点号 `.`
- 区分大小写
- 常量名必须使用帕斯卡命名法（PascalCase），如 `MAX_SIZE`

### 10. 语句分隔

- 使用换行符作为语句分隔符
- 空行会被忽略

## 示例

以下是一个完整的IDL文件示例：

```
// 用户服务接口定义

// 用户状态枚举
enum UserStatus {
    ACTIVE = 1 (errmsg="active user")
    INACTIVE = 0 (errmsg="inactive user")
}

// 错误码定义
enum ErrCode {
    ERR_OK = 0 (errmsg="success")
    PARAM_ERROR = 1003 (errmsg="parameter error")
}

// 扩展错误码
enum extends ErrCode {
    USER_NOT_FOUND = 404 (errmsg="user not found")
    PERMISSION_DENIED = 403 (errmsg="permission denied")
}

// 用户信息结构
type User {
    required string id
    required string name (validate="$ != '' && len($) <= 64")
    optional string email (validate="email($)")
    int age (validate="$ >= 0 && $ <= 150")
    UserStatus status
    list<string> roles
    map<string, string> metadata
}

// 获取用户请求
type GetUserRequest {
    required string userId (path="id", validate="$ != ''")  // 路径参数绑定
    optional string locale (query="locale")  // 查询参数绑定
}

// 获取用户响应
type GetUserResponse {
    required User user
    int code
    string message
}

// 通用响应包装
type Response<T> {
    int code
    string message
    T data
}

// 用户服务接口
rpc GetUser (GetUserRequest) GetUserResponse {
    method = "GET"
    path = "/user/:id"  // RESTful路径参数
}

// 批量获取用户
rpc BatchGetUser (BatchGetUserRequest) Response<list<User>> {
    method = "POST"
    path = "/user/batch"
}
```

## 注意事项

1. 字段名和类型名区分大小写
2. 支持使用泛型来定义可重用的数据结构
3. 注解提供了灵活的元数据机制，用于控制生成代码的行为
4. RPC接口可以包含路径参数、查询参数和请求体参数
5. 支持通过注解指定HTTP方法、路径等路由信息