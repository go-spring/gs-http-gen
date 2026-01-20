# IDL 语法中文说明文档

## 概述

gs-http-gen 采用专门设计的接口定义语言（IDL）来描述数据模型和接口服务。
该语言支持定义常量、枚举、结构体、联合类型以及 RPC 服务，
旨在简化 API 开发流程并确保前后端数据结构的一致性。

## IDL 文件整体构成

一个 IDL 文件由以下几个主要部分组成：

- **注释**：用于功能说明和文档描述
- **关键字**：IDL 预定的保留字
- **常量定义**：定义固定不变的值
- **枚举定义**：定义枚举类型，通常用于状态码或选项
- **类型定义**：定义结构体、联合类型等复合数据类型
- **RPC 接口定义**：定义服务接口，支持普通接口和流式接口

语句使用换行符作为分隔符，空行会被忽略。缩进使用 4 个空格，不建议使用 Tab 键。

## 基本语法规则

### 注释

- 单行注释：使用 `//` 或 `#` 开头
- 多行注释：使用 `/* ... */` 包围

### 空白字符处理

- 换行符、空格和制表符被视为空白字符
- 多个连续换行符在语法解析中是允许的
- 文档可以以任意数量的换行符开始或结束
- 定义之间可以使用多个换行符分隔

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

### 标识符规则

- 以字母开头
- 可以包含字母、数字、下划线 `_` 和点号 `.`
- 区分大小写
- **首字母必须大写**，推荐在以下场景使用大写：
    - 常量名：必须使用帕斯卡命名法（PascalCase），如 `MAX_SIZE`
    - 枚举名：必须使用帕斯卡命名法（PascalCase），如 `Status`
    - 枚举字段：必须使用帕斯卡命名法（PascalCase），如 `ACTIVE`
    - 类型名：必须使用帕斯卡命名法（PascalCase），如 `User`
    - 方法名：必须使用帕斯卡命名法（PascalCase），如 `GetUser`

**注**：标识符也可以用作常量值或注解值。

### 格式与布局规则

- 缩进：使用 4 个空格进行缩进，不使用 Tab 键
- 换行：语句使用换行符分隔，每条语句独占一行
- 空行：空行用于分隔逻辑块，提高可读性
- 大括号：左大括号 `{` 不换行，紧跟在定义后，右大括号 `}` 独占一行

### 基础类型

- `bool` - 布尔类型
- `int` - 整数类型（映射到Go的int64，默认使用此类型，若需自定义可使用`go.type`注解）
- `float` - 浮点数类型（映射到Go的float64，默认使用此类型，若需自定义可使用`go.type`注解）
- `string` - 字符串类型（只能使用双引号，单引号不行）
- `bytes` - 字节数组类型（传输时使用 base64 编码成字符串）

### 字面量格式

- 整数字面量：`42`, `-17`, `0x1A2B`
- 浮点数字面量：`3.14`, `.5`, `-2.7e10`
- 字符串字面量：`"hello"`, `"escaped \" quote"`（只支持双引号）
- 布尔字面量：`true`, `false`

### 容器类型

- `map<K, V>` - 映射类型，键类型只能是 `int` 或 `string`，值类型可以是任意类型
- `list<T>` - 列表类型，元素类型可以是任意类型
- `bytes` - 字节数组类型（传输时使用 base64 编码成字符串）

**注意**：map 的键类型只能是 `int` 或 `string`，否则会在解析时抛出错误。

map 和 list 都支持嵌套，例如：`list<list<int>>`、`map<string,map<string,int>>` 等

示例：

```
list<string> tags                    // 字符串列表
list<User> users                     // 用户对象列表
map<string, int> scores             // 字符串到整数的映射
map<int, User> userById             // 整数ID到用户的映射
list<map<string, User>> groups      // 用户组列表（每个组是用户映射）
map<string, list<User>> usersByDept // 按部门分组的用户映射
```

## 语法详细说明

### 1. 常量

#### 1.1 常量

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

**注意**：在常量定义中，不能使用枚举类型。

#### 1.2 常量值

支持以下类型的常量值：

- 整数字面量：`42`, `-17`, `0x1A2B`
- 浮点数字面量：`3.14`, `.5`, `-2.7e10`
- 字符串字面量：`"hello"`, `"escaped \" quote"`
- 布尔值：`true`, `false`
- 标识符：通常用于引用枚举成员

**注意**：在注解中可以使用枚举作为常量值。

### 2. 注解（Annotations）

注解（Annotations）是IDL的重要组成部分，用于为字段、类型、枚举和接口添加元数据信息。
注解提供了灵活的元数据机制，用于控制生成代码的行为，可以指定序列化方式、验证规则、API路由等信息。

注解支持单行定义，这种格式只能用在字段上：

```
(注解名 = 注解值, 注解名 = 注解值, ...)
```

注解语法支持多行定义，例如：

```
# 类型上的注解
type User {
    string name (
      json="name"
      go.type="string"
    )
    int age (json="age,omitempty")
}

# 接口中的注解
rpc GetUser (GetUserRequest) GetUserResponse {
    method = "GET"
    path = "/user/:id"
}
```

注解可以在以下场景中使用：

- 枚举定义中：为错误码添加描述信息
- 字段定义中：指定JSON序列化规则、验证规则、参数绑定等
- RPC接口定义中：指定HTTP方法、路径、超时时间等路由信息

**注解语法细节**：

- 注解值可以是布尔值、整数、浮点数、字符串或标识符
- 注解之间可以用逗号或换行符分隔
- 注解可以是单独的标识符，不一定要是键值对形式（例如：`(some_flag)`）
- 注解的完整语法为：`<标识符>(=<常量值>)?`

### 3. 枚举定义

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

#### 3.1 错误码枚举

有一种特殊的 enum 类型，叫错误码。
标识一个 enum 类型是错误码的关键是在字段后面添加 `errmsg` 注解，例如：

```
enum ErrCode {
    ERR_OK = 0 (errmsg="success")
    PARAM_ERROR = 1003 (errmsg="parameter error")
    NOT_FOUND = 404 (errmsg="resource not found")
}
```

错误码注解（errmsg）用于为错误码提供人类可读的消息描述，在生成代码时会创建相应的错误消息映射。

#### 3.2 错误码扩展

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

#### 3.3 enum_as_string 特性

当在结构体字段中使用枚举类型时，可以使用 `enum_as_string` 注解将枚举值作为字符串进行序列化和反序列化。

使用 `enum_as_string` 注解的字段在生成代码时会创建两个类型：
原始枚举类型和字符串形式的枚举类型（AsString类型），并在序列化/反序列化时自动处理类型转换。

示例：

```
enum Department {
    ENGINEERING = 1
    MARKETING = 2
    SALES = 3
}

type Manager {
    string name
    Department dept (enum_as_string)  // 枚举将以字符串形式序列化
}
```

在这个例子中，`dept` 字段在JSON序列化时将输出字符串形式（如 "ENGINEERING"），而不是整数值（如 1）。

生成的Go代码将包含：

1. 原始枚举类型定义
2. AsString 类型定义（用于字符串序列化）
3. MarshalJSON 方法（将枚举值序列化为字符串）
4. UnmarshalJSON 方法（将字符串反序列化为枚举值）

这种特性特别适用于需要与前端交互或外部API集成的场景，因为字符串形式的枚举更容易理解和调试。

### 4. 类型定义

类型定义使用 `type` 关键字，有三种定义方式：

#### 4.1 普通结构体定义

最常见的是普通结构体，定义方式如下：

```
type User {
    required string name
    int age
    optional string email
    list<string> tags
}
```

#### 4.2 泛型结构体定义

支持简单的泛型定义，一般用于对返回值进行定义的场景：

```
type Response<T> {
    int code
    string message
    T data
}
```

语法格式，尖括号内定义泛型参数：

```
type <类型名><<泛型参数>> {
    <字段定义>
}
```

#### 4.3 泛型实例化

对已定义的泛型进行实例化，注意语法格式：

```
type UserResponse Response<User>
```

语法格式：

```
type <新类型名> <已定义的类型名><<具体类型>>
```

#### 4.4 结构体字段

结构体字段包括普通字段、泛型字段和带有注解的字段：

**4.4.1 普通字段**

普通字段定义格式：

```
<可选修饰符> <类型> <字段名>
```

修饰符：

- `required` - 表示必需字段
- `optional` - 表示可选字段（默认）

示例：

```
type User {
    required string name  // 必填字段
    int age              // 可选字段（默认）
    optional string email // 显式声明可选字段
}
```

**4.4.2 泛型字段**

在结构体中可以使用泛型参数作为字段类型：

```
type Result<T> {
    bool success
    T data
    string message
}
```

**4.4.3 字段注解**

在字段定义中可以使用注解来指定额外的元数据。以下是常用注解的分类：

**自定义类型注解**：

- `go.type` - 指定Go语言的具体类型，例如 `go.type="int32"` 将int类型映射到Go的int32

**序列化注解**：

- `json` - 指定JSON序列化字段名和选项
- `enum_as_string` - 枚举作为字符串处理

**参数绑定注解**：

- `path` - 指定路径参数绑定
- `query` - 指定查询参数绑定

**验证注解**：

- `validate` - 指定验证表达式

**其他注解**：

- `deprecated` - 标记字段已弃用
- `compat_default` - 兼容性默认值

示例：

```
type User {
    string name (json="name", go.type="string")
    int age (json="age,omitempty", go.type="int32")  // 显式指定Go类型为int32
    string email (validate="email($)", deprecated="true")
    string userId (path="id")  // 绑定路径参数
    string locale (query="locale")  // 绑定查询参数
}
```

**4.4.4 嵌入类型**

可以在结构体中嵌入其他类型，代码生成时会把嵌入的类型进行展开：

```
type Address {
    string street
    string city
}

type Person {
    Address  // 嵌入 Address 类型，无字段名
    string name
}
```

嵌入类型在生成代码时的合并规则：

1. 嵌入类型的所有字段会被直接合并到当前类型中
2. 如果嵌入类型中有同名字段，会导致编译错误
3. 嵌入类型本身不会作为一个独立字段存在
4. 嵌入类型语法：直接使用类型名，无需字段名和类型声明

#### 4.5 序列化和反序列化

为了支持可选和必选字段，以及默认值填充，生成的代码必须支持流式解析。
这允许在解析过程中校验必选字段是否存在，以及为字段设置默认值。

- 必选字段（`required`）：在解析和验证阶段必须存在且非空
- 可选字段（`optional`）：可以不存在，通常表示为指针类型
- 默认值（`compat_default`）：为字段提供兼容性默认值

#### 4.6 联合类型定义

使用 `oneof` 定义联合类型。
联合类型在内部实现上会生成一个特殊结构，包含一个类型标识字段和多个可能的数据字段。
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

### 5. Validate 注解

Validate 注解用于对字段值进行验证，支持多种内置验证函数和复杂的表达式。
它是一种强大的数据校验机制，可以在运行时确保数据的有效性。

#### 5.1 验证函数

内置验证函数包括：

- `len` - 验证长度（字符串、列表、映射等）

除了内置验证函数外，还可以定义自定义验证函数。
自定义验证函数的名称不能与内置函数重名，但可以被多次使用，只要它们应用于相同类型的字段。
自定义验证函数会在生成的代码中创建对应的函数模板，供开发者实现具体的验证逻辑。

#### 5.2 支持的操作符

验证表达式支持的操作符：

- 比较操作符：`==`, `!=`, `<`, `<=`, `>`, `>=`
- 逻辑操作符：`&&`, `||`, `!`
- 算术操作符：`+`, `-`, `*`, `/`
- 函数调用：`len($)`、`email($)`、`regexp($, "pattern")` 等

验证表达式遵循标准的运算符优先级规则，括号可用于显式指定计算顺序。

#### 5.3 特殊变量

- `$` 代表当前字段的值
- 在复杂表达式中，可以使用 `$` 引用当前字段进行验证

#### 5.4 验证规则

- 如果字段是可选的（optional），验证仅在字段值存在时执行
- 如果字段是必需的（required），验证将在值存在时执行
- 验证失败时会抛出特定的错误信息

#### 5.5 使用示例

```
type User {
    required string name (validate="$ != '' && len($) <= 64")  // 非空且长度不超过64
    int age (validate="$ >= 0 && $ <= 150")                  // 年龄在0-150之间
    string email (validate="email($)")                        // 邮箱格式验证
    list<string> tags (validate="len($) <= 10")               // 标签数量不超过10个
    string phone (validate="phone($)")                        // 手机号格式验证
    string username (validate="regexp($, \"^[a-zA-Z][a-zA-Z0-9_]{2,19}$\")")  // 用户名格式验证
}
```

#### 5.6 复杂验证表达式

可以组合多个验证条件：

```
type Product {
    string code (validate="len($) >= 3 && len($) <= 20 && regexp($, \"^\\w+$\")")  // 3-20位字母数字下划线
    int price (validate="$ > 0 && $ <= 999999")              // 价格范围验证
    list<string> images (validate="len($) >= 1 && len($) <= 10")  // 图片数量限制
}
```

### 6. RPC 接口定义

RPC接口定义有两种形式：普通接口和SSE（Server-Sent Events）流式接口。

在项目中，多个IDL文件可以协同工作。系统会自动解析跨文件的类型引用，确保所有使用的类型都正确定义。
此外，错误码可以通过 `extends` 关键字在不同文件间进行扩展合并。

#### 6.1 普通接口定义

普通接口定义格式：

```
rpc <接口名> (<请求类型>) <响应类型> {
    <注解定义>
}
```

示例：

```
rpc GetUser (GetUserRequest) GetUserResponse {
    method = "GET"
    path = "/user/:id"
}
```

#### 6.2 SSE流式接口定义

SSE（Server-Sent Events）流式接口是现代AI编程中必须的一种能力，用于实现实时数据推送和长连接通信。
SSE接口定义使用 `sse` 关键字：

```
sse <接口名> (<请求类型>) <响应类型> {
    <注解定义>
}
```

示例：

```
sse StreamEvents (StreamRequest) Event {
    method = "GET"
    path = "/events"
}
```

SSE接口特别适用于以下场景：

1. 实时数据推送（如股票价格、聊天消息）
2. AI模型结果流式返回（逐步返回AI生成的内容）

#### 6.3 HTTP相关注解

RPC接口可以使用注解来指定HTTP方法、路径、超时时间等路由信息，常见的HTTP相关注解包括：

- `method` - HTTP方法（GET、POST、PUT、DELETE等）
- `path` - 请求路径，支持RESTful路径参数
- `content-type` - 内容类型
- `connTimeout` - 连接超时
- `readTimeout` - 读取超时
- `writeTimeout` - 写入超时
- `summary` - 接口摘要说明

例如：

```
rpc GetUser (GetUserRequest) GetUserResponse {
    method = "GET"
    path = "/user/:id"
    connTimeout = "5s"
    readTimeout = "10s"
    summary = "获取用户信息"
}
```

#### 6.4 RESTful Path 规则和参数绑定

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

**参数绑定**：

在请求类型中可以使用注解将字段与路径或查询参数关联：

```
type GetUserRequest {
    string userId (path="id")  // 将 userId 字段绑定到路径参数 :id 或 {id}
    string locale (query="locale")  // 将 locale 字段绑定到查询参数 ?locale=value
}
```

注意：

- 路径参数对应的字段必须是必需的（required）
- 支持的参数类型包括：int、string 等基本类型
- 参数名必须与路径中定义的参数名匹配

Query参数绑定的特点：

- 查询参数通常是可选的
- 可以设置默认值
- 适合传递过滤条件、分页参数等

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

// 部分枚举
enum Department {
    ENGINEERING = 1
    MARKETING = 2
    SALES = 3
}

// 用户信息结构
type User {
    required string id
    required string name (validate="$ != '' && len($) <= 64")  // 非空且长度不超过64
    optional string email (validate="email($)")
    int age (validate="$ >= 0 && $ <= 150", go.type="int32")  // 显式指定Go类型为int32
    UserStatus status
    Department dept (enum_as_string)  // 使用 enum_as_string 特性
    list<string> roles
    map<string, string> metadata
}

// 地址信息结构
type Address {
    string street
    string city
}

// 个人信息结构（嵌入地址）
type Person {
    Address  // 嵌入 Address 类型
    string name
    int age
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

// 流式事件请求
type StreamRequest {
    required string clientId (path="id")
    optional string eventType (query="type")
}

// 流式事件响应
type StreamResponse {
    string eventId
    string eventType
    string data
    int timestamp
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
    summary = "获取用户信息"
}

// 批量获取用户
rpc BatchGetUser (BatchGetUserRequest) Response<list<User>> {
    method = "POST"
    path = "/user/batch"
    summary = "批量获取用户信息"
}

// SSE流式事件接口
sse StreamEvents (StreamRequest) StreamResponse {
    method = "GET"
    path = "/events/:id"
    summary = "流式事件推送，适用于实时数据更新"
}
```

## 注意事项

1. 字段名和类型名区分大小写
2. 支持使用泛型来定义可重用的数据结构
3. 注解提供了灵活的元数据机制，用于控制生成代码的行为
4. RPC接口可以包含路径参数、查询参数和请求体参数
5. 支持通过注解指定HTTP方法、路径等路由信息
6. map 的键类型只能是 `int` 或 `string`，否则会在解析时抛出错误
7. 基础类型不能作为单独的字段类型使用，只能作为容器类型（如list、map）的元素类型或常量类型
8. 在枚举定义中使用 `extends` 关键字时，扩展的枚举类型必须已经定义
9. 路径参数对应的字段必须是必需的（required）
10. 生成的代码使用基于哈希的字段调度机制进行高性能解析
11. 项目支持多IDL文件协作，系统会自动处理跨文件的类型引用
12. 自定义验证函数需要在生成的 validate.go 文件中手动实现具体的验证逻辑
13. 嵌入类型时，嵌入的类型所有字段会被直接合并到当前类型中，需避免字段名冲突
14. 泛型实例化时，系统会将泛型参数替换为具体类型
