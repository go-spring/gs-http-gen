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

使用 `enum` 关键字定义枚举类型：

```
enum Color {
    RED = 1
    GREEN = 2
    BLUE = 3
}
```

有一种特殊的 enum 类型，叫错误码，标识一个 enum 类型是错误码的关键是在字段后面添加
errmsg 注解，例如：

```
enum extends ErrorCode {
    SUCCESS = 0 (errmsg="success")
}
```
 
然后我们还可以对错误码进行扩展，代码生成的时候基础错误码和扩展错误码会进行合并：

```
enum extends ErrorCode {
    ERROR_PARAM = 1
    ERROR_SERVER = 2
}
```

note: 枚举类型的值不能重复，而且最好保持单调递增。todo（好像缺少格式说明）

### 3. 类型定义

#### 3.1 结构体类型

用 type 关键字定义结构体，字段支持 optional（可选）和 required（必需）：

```
type User {
    required string name
    int age
    optional string email
    list<string> tags
}
```

todo（其实字段上可以添加很多注解，也可以介绍一下）

#### 3.2 泛型类型

支持简单的泛型定义，一般用于对返回值进行定义的场景，这时候一般只需要对 data 进行泛型定义：

```
type Response<T> {
    int code
    string message
    T data
}
```

#### 3.3 泛型实例化

对泛型进行实例化，注意语法格式，只能对 user_type 泛型进行实例化：

```
type UserResponse Response<User>
```

### 4. 字段定义

#### 4.1 普通字段

```
<可选修饰符> <类型> <字段名>
```

修饰符：

- `required` - 表示必需字段
- `optional` - 表示可选字段（默认）

#### 4.2 嵌入类型

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

### 5. 联合类型定义

使用 `oneof` 定义联合类型，todo（需要补充编码格式）：

```
oneof Value {
    User
    Manager
}
```

### 6. RPC 接口定义

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

todo（缺少很多注解，还是要参考语法和golang实现）

### 7. 注解（Annotations）

字段和接口可以使用注解来添加元数据：

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

### 10. 语句分隔

- 使用换行符作为语句分隔符
- 空行会被忽略

## 示例

以下是一个完整的IDL文件示例：

```
// 用户服务接口定义

// 用户状态枚举
enum UserStatus {
    ACTIVE = 1
    INACTIVE = 0
}

// 用户信息结构
type User {
    required string id
    required string name
    optional string email
    int age
    UserStatus status
    list<string> roles
    map<string, string> metadata
}

// 获取用户请求
type GetUserRequest {
    required string userId (json="userId")
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
    path = "/user/:userId"
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