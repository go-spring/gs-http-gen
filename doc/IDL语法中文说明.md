# IDL 语法中文说明文档

## 概述

gs-http-gen 是一种专门设计的接口定义语言（IDL），用于描述数据模型和接口服务。
该语言支持定义常量、枚举、结构体、联合类型以及 RPC 服务，
旨在简化 API 的开发过程，并确保前后端的数据结构保持一致。

## IDL 项目结构

一个完整的 gs-http-gen IDL 项目由以下组成部分构成：

### 项目目录结构

一个标准的 IDL 项目应包含以下内容：

1. **meta.json** - 项目的元信息配置文件（必需）
2. **IDL 文件** - 扩展名为 `.idl` 的接口定义文件（至少一个）

典型的项目目录结构如下所示：

```
project/
├── meta.json           # 项目元信息配置文件
├── service.idl         # 主要服务定义
├── user.idl            # 用户模块定义
├── common.idl          # 通用类型定义
└── error.idl           # 错误码定义
```

### meta.json 配置文件

**meta.json** 是 IDL 项目的核心配置文件，位于项目的根目录，用于存储项目的元数据。
缺少此文件，项目将无法正确解析。

meta.json 文件的内容格式示例如下：

```json
{
  "name": "my-service",
  "version": "1.0.0",
  "description": "My service description"
}
```

主要字段说明：

- `name`: 项目名称，通常用于生成代码中的包名或命名空间
- `version`: 项目的版本号
- `description`: 项目的简要描述

### IDL 文件

IDL 文件以 `.idl` 扩展名命名，用于定义数据结构、接口和服务。一个项目可以包含多个 IDL
文件，所有文件共享同一个命名空间，可以相互引用其他文件中定义的类型、枚举和常量。

#### 命名空间与类型引用

- 所有 IDL 文件共享同一个全局命名空间
- 在任何 IDL 文件中定义的类型、枚举和常量，可以在其他文件中直接引用
- 类型名称在项目中必须唯一，不能在不同文件中定义相同的类型
- 系统会自动解析不同文件之间的类型引用关系

#### 文件组织建议

尽管 IDL 文件可以互相引用，但建议按功能将文件进行组织，常见的组织方式包括：

- **通用类型文件**：定义项目中通用的基础类型和枚举
- **业务模块文件**：按照业务领域划分，如用户模块、订单模块等
- **服务接口文件**：定义 RPC 接口及其请求和响应类型
- **错误码文件**：集中管理项目中的错误码枚举

## IDL 文件构成

一个 IDL 文件由以下几个主要部分组成：

- **注释**：用于功能说明和文档描述
- **关键字**：IDL 语言的保留字
- **常量定义**：定义固定不变的常量值
- **枚举定义**：定义枚举类型，通常用于状态码或选项
- **类型定义**：定义结构体、联合类型等复合数据类型
- **RPC 接口定义**：定义服务接口，支持普通 RPC 方法和流式接口

语句以换行符分隔，空行将被忽略。建议使用 4 个空格进行缩进，避免使用 Tab 键。

## 基本语法规则

### 注释

- **单行注释**：使用 `//` 或 `#` 开头
- **多行注释**：使用 `/* ... */` 包裹

### 空白字符处理

- 换行符、空格和制表符均视为空白字符
- 允许多个连续的换行符
- 文档可以以任意数量的换行符开始或结束
- 定义之间可以使用多个换行符进行分隔

### 关键字

以下是 IDL 语言的保留关键字，无法作为标识符使用：

- `extends` - 扩展错误码
- `const` - 定义常量
- `enum` - 定义枚举
- `type` - 定义类型
- `oneof` - 定义联合类型
- `rpc` - 定义普通 RPC 方法
- `sse` - 定义流式 RPC 方法
- `true` / `false` - 布尔值
- `optional` - 可选字段标记
- `required` - 必填字段标记

### 标识符规则

- 标识符必须以字母开头
- 可以包含字母、数字、下划线 `_` 和点号 `.`
- 区分大小写
- 推荐在以下场景使用首字母大写：

    - **常量名**：如 `MAX_SIZE`
    - **枚举名**：如 `Status`
    - **枚举字段**：如 `ACTIVE`
    - **类型名**：如 `User`
    - **方法名**：如 `GetUser`

**注意**：标识符也可以作为常量值或注解值。

### 格式与布局规则

- **缩进**：使用 4 个空格进行缩进，避免使用 Tab 键
- **换行**：每条语句使用换行符分隔，每行单独一条语句
- **空行**：空行用于分隔逻辑块，提升可读性
- **大括号**：

    - 左大括号 `{` 紧跟在定义后，不换行
    - 右大括号 `}` 独占一行

### 基础类型

- **`bool`**：布尔类型
- **`int`**：整数类型，映射为 Go 的 `int64`（默认使用此类型）。如果需要自定义类型，可以使用 `go.type` 注解。
- **`float`**：浮点数类型，映射为 Go 的 `float64`（默认使用此类型）。同样，可以通过 `go.type` 注解进行自定义。
- **`string`**：字符串类型，仅支持双引号（单引号不合法）。
- **`bytes`**：字节数组类型，传输时使用 Base64 编码为字符串。

### 字面量格式

- **整数**：`42`，`-17`，`0x1A2B`
- **浮点数**：`3.14`，`.5`，`-2.7e10`
- **字符串**：`"hello"`，`"escaped \" quote"`（仅支持双引号）
- **布尔值**：`true`，`false`

### 容器类型

- **`map<K, V>`**：映射类型，键类型必须是 `int` 或 `string`，值类型可以是任意类型。
- **`list<T>`**：列表类型，元素类型可以是任意类型。
- **`bytes`**：字节数组类型（传输时以 Base64 编码为字符串）。

**注意**：`map` 的键类型仅支持 `int` 或 `string`，其他类型会在解析时抛出错误。

`map` 和 `list` 都支持嵌套，例如：`list<list<int>>`，`map<string,map<string,int>>` 等。

**重要限制**：字段类型不能直接使用泛型，只能使用预先定义好的泛型类型实例化后的具体类型。
只有容器类型（`list<T>` 和 `map<K,V>`）可以直接在字段中使用泛型参数。

### 示例

```idl
list<string> tags                     // 字符串列表
list<User> users                      // 用户对象列表
map<string, int> scores               // 字符串到整数的映射
map<int, User> userById               // 整数ID到用户的映射
list<map<string, User>> groups        // 用户组列表（每个组是用户映射）
map<string, list<User>> usersByDept  // 按部门分组的用户映射
```

## 语法详细说明

### 1. 常量

#### 1.1 常量定义

使用 `const` 关键字来定义常量，常量的类型可以是以下四种基本类型之一：
`bool`、`int`、`float` 或 `string`。示例如下：

```
const string APP_NAME = "MyApp"
const int MAX_SIZE = 100
const float PI = 3.14159
const bool DEBUG = true
```

常量的定义遵循以下语法规则：

```
const <基础类型> <标识符名称> = <常量值>
```

**注意**：常量的值不能是枚举类型。

#### 1.2 常量值类型

常量支持以下几种类型的字面量值：

- **整数字面量**：`42`, `-17`, `0x1A2B`
- **浮点数字面量**：`3.14`, `.5`, `-2.7e10`
- **字符串字面量**：`"hello"`, `"escaped \" quote"`
- **布尔值**：`true`, `false`
- **标识符**：通常用于引用枚举成员

**注意**：在注解中，你可以使用枚举作为常量值。

### 2. 注解（Annotations）

注解（Annotations）是IDL（接口定义语言）中的重要组成部分，主要用于为字段、类型、枚举和接口添加元数据。
注解为生成的代码行为提供灵活的控制机制，可以指定序列化规则、验证规则、API路由等。

注解的定义支持单行和多行两种格式：

- **单行定义**：适用于字段，格式为：

```
(注解名 = 注解值, 注解名 = 注解值, ...)
```

- **多行定义**：可以在字段、类型或接口中使用，示例如下：

```idl
# 类型定义中的注解
type User {
    string name (
      json="name"
      go.type="string"
    )
    int age (json="age,non-omitempty")
}

# 接口定义中的注解
rpc GetUser (GetUserRequest) GetUserResponse {
    method = "GET"
    path = "/user/:id"
}
```

注解可在以下场景中使用：

- **枚举定义**：为错误码等枚举项添加描述信息；
- **字段定义**：指定字段的JSON序列化规则、验证规则、参数绑定等；
- **RPC接口定义**：指定HTTP方法、路径、超时时间等路由信息。

#### 注解语法细节：

- **注解值类型**：注解的值可以是布尔值、整数、浮点数、字符串或标识符；
- **分隔符**：多个注解之间可以使用逗号或换行符进行分隔；
- **标识符形式**：注解可以仅包含标识符，不必是键值对的形式，例如：`(some_flag)`；
- **完整语法**：注解的完整语法为：`<标识符>(=<常量值>)?`。

#### RPC接口常用注解：

- **路由相关**：
    - `method` - 指定HTTP方法（例如：GET、POST、PUT、DELETE等）
    - `path` - 请求路径，支持RESTful路径参数
- **内容类型**：
    - `contentType` - 请求内容类型，支持 "form"（表单编码）或 "json"（JSON编码）
- **超时设置**（治理策略，必须设置）：
    - `connTimeout` - 连接超时，单位毫秒
    - `readTimeout` - 读取超时，单位毫秒
    - `writeTimeout` - 写入超时，单位毫秒
- **接口文档**：
    - `resp.go.type` - 指定RPC响应的Go类型

### 3. 枚举定义

使用 `enum` 关键字定义枚举类型。每个枚举项都必须赋予一个整数值。

基本定义格式如下：

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
    RED = 1 (desc="红色")
    GREEN = 2 (desc="绿色")
    BLUE = 3 (desc="蓝色")
}
```

#### 3.1 错误码枚举

一种特殊类型的枚举是错误码枚举。要定义一个错误码枚举，必须在枚举字段后添加 `errmsg` 注解，该注解为每个错误码提供可读的错误消息描述。

例如：

```
enum ErrCode {
    ERR_OK = 0 (errmsg="成功")
    PARAM_ERROR = 1003 (errmsg="参数错误")
    NOT_FOUND = 404 (errmsg="资源未找到")
}
```

`errmsg` 注解会在代码生成时为每个错误码生成相应的可读错误消息映射。

**注意**：普通枚举（如状态、分类等）应使用 `desc` 注解来描述枚举项含义，而错误码枚举必须使用 `errmsg` 注解。

#### 3.2 错误码扩展

错误码可以通过 `extends` 关键字进行扩展。扩展的错误码将与基础错误码合并，形成完整的错误码集合。

示例：

```
// 基础错误码定义
enum ErrCode {
    ERR_OK = 0 (errmsg="成功")
    PARAM_ERROR = 1003 (errmsg="参数错误")
}

// 扩展错误码定义
enum extends ErrCode {
    USER_NOT_FOUND = 404 (errmsg="用户未找到")
    PERMISSION_DENIED = 403 (errmsg="权限不足")
}
```

在代码生成过程中，扩展的错误码将与基础错误码合并，形成一个完整的错误码集。扩展错误码时需要注意以下几点：

1. 扩展的错误码类型必须事先定义。
2. 扩展的错误码字段名和值不能与已有的错误码重复。
3. 扩展错误码的值应保持单调递增，避免重复或回退的值。
4. 在代码生成时，基础错误码和扩展错误码会被一起处理，生成对应的Go代码。

#### 3.3 `enum_as_string` 特性

当结构体字段使用枚举类型时，可以通过 `enum_as_string` 注解将枚举值以字符串形式进行序列化和反序列化。

使用 `enum_as_string` 注解的字段，在生成代码时会生成两个版本的枚举类型：

1. 原始的枚举类型
2. 字符串形式的枚举类型（即 `AsString` 类型）

在序列化和反序列化过程中，这两个类型会自动处理转换，从而确保枚举值以字符串的形式存储和传输。

**示例：**

```idl
enum Department {
    ENGINEERING = 1 (desc="工程技术部")
    MARKETING = 2 (desc="市场部")
    SALES = 3 (desc="销售部")
}

type Manager {
    string name
    Department dept (enum_as_string)  // 将枚举字段以字符串形式序列化
}
```

在上述示例中，`dept` 字段在 JSON 序列化时将输出字符串表示（例如 `"ENGINEERING"`），而不是枚举的整数值（如 `1`）。

生成的 Go 代码将包括：

1. 原始的枚举类型定义
2. `AsString` 类型定义（用于字符串序列化）
3. `MarshalJSON` 方法（将枚举值序列化为字符串）
4. `UnmarshalJSON` 方法（将字符串反序列化为枚举值）

该特性非常适用于与前端进行数据交互或与外部 API 集成的场景，因为字符串形式的枚举值更加直观、易于理解和调试。

### 4. 类型定义

类型定义使用 `type` 关键字，并支持三种主要定义方式：

#### 4.1 普通结构体定义

普通结构体是最常见的类型定义方式，格式如下：

```
type User {
    required string name       // 必填字段
    int age                    // 可选字段
    optional string email      // 显式声明可选字段
    list<string> tags          // 字符串列表
}
```

#### 4.2 泛型结构体定义

IDL 语言支持简单的泛型定义，通常用于定义返回值的结构体。其语法格式如下：

```
type Response<T> {
    int code                    // 返回码
    string message               // 返回消息
    T data                       // 泛型数据
}
```

在定义泛型类型时，通过尖括号 `< >` 来指定泛型参数：

```
type <类型名><<泛型参数>> {
    <字段定义>
}
```

**注意**：字段类型不能直接使用泛型，只能使用预先定义好的泛型类型实例化后的具体类型。
只有容器类型（`list<T>` 和 `map<K,V>`）可以直接在字段中使用泛型参数。

#### 4.3 泛型实例化

对于已定义的泛型类型，可以进行实例化。实例化时，需要将泛型类型参数替换为具体的类型。例如：

```
type UserResponse Response<User>  // 实例化泛型类型，指定泛型为 User 类型
```

语法格式为：

```
type <新类型名> <已定义的类型名><<具体类型>>
```

#### 4.4 结构体字段

结构体字段可以分为普通字段、泛型字段和带有注解的字段。每种类型的字段都有其特定的用途和语法。

**4.4.1 普通字段**

普通字段的定义格式如下：

```
<修饰符> <类型> <字段名>
```

**修饰符**：

- `required`：表示该字段是必需的，不能省略。
- `optional`：表示该字段是可选的，如果未提供值，系统可以将其视为缺省值（默认为可选）。

**示例**：

```
type User {
    required string name  // 必填字段
    int age              // 可选字段（默认为可选）
    optional string email // 显式声明可选字段
}
```

**4.4.2 泛型字段**

在结构体中，可以使用泛型作为字段类型，允许更加灵活和可复用的定义。定义泛型字段时，可以指定具体的类型参数。

**示例**：

```
type Result<T> {
    bool success       // 泛型字段 success
    T data             // 泛型字段 data，类型为T
    string message     // 字段 message
}
```

**4.4.3 字段注解**

字段注解是用于为字段添加元数据信息，控制生成代码时的行为。注解可以指定如何序列化、如何验证字段值，或提供额外的元数据。常用的字段注解包括以下几类：

**自定义类型注解**：

- `go.type`：指定该字段在Go语言中的具体类型。例如，`go.type="int32"` 将该字段的类型映射为Go中的 `int32`。

**序列化注解**：

- `json`：指定字段在JSON序列化时的字段名和选项，例如 `json="name"` 会将该字段在JSON中的名称设为 `name`。
- `enum_as_string`：将枚举类型的字段序列化为字符串形式而非其整数值。
- `form`：指定该字段在表单数据序列化时的字段名及相关选项。

**参数绑定注解**：

- `path`：用于指定字段与路径参数的绑定关系。例如，`path="id"` 将字段与路径参数 `:id` 绑定。
- `query`：用于指定字段与查询参数的绑定关系。例如，`query="locale"` 将字段与查询参数 `locale` 绑定。

**验证注解**：

- `validate`：为字段指定验证规则。该规则可以是任意逻辑表达式，用于确保字段值的合法性。

**其他注解**：

- `deprecated`：标记该字段已弃用，通常在版本升级中使用，以提示开发者该字段不再推荐使用。
- `compat_default`：为字段提供兼容性默认值，即使在没有传递值的情况下，字段也会使用默认值。

**示例**：

```
type User {
    string name (go.type="string")  // 指定JSON字段名为 name，并在Go中使用 string 类型
    int age (json="age", go.type="int32")  // JSON序列化时，如果值为空则不序列化，并显式指定Go类型为 int32
    string description (json="desc,non-omitempty")  // 禁用omitempty行为，确保空字符串也会被序列化
    string email (validate="email($)", deprecated="true")  // 对邮箱字段进行格式验证，并标记为已弃用
    string userId (path="id")  // 将字段与路径参数 id 绑定
    string locale (query="locale")  // 将字段与查询参数 locale 绑定
}
```

#### 4.5 嵌入类型

在结构体中，可以将其他类型作为嵌入类型。嵌入类型的字段会在代码生成时自动展开并合并到当前结构体中。嵌入类型本身不会作为一个独立字段存在，而是直接将其所有字段合并进当前类型。

例如：

```idl
type Address {
    string street
    string city
}

type Person {
    Address  // 嵌入 Address 类型，无字段名
    string name
}
```

在代码生成时，`Person` 类型将包含 `Address` 类型的所有字段，且不再有 `Address` 字段名。生成的结构体类似如下：

```
type Person struct {
    street string // 来自 Address 类型
    city   string // 来自 Address 类型
    name   string
}
```

**嵌入类型合并规则**：

1. 嵌入类型的所有字段会被直接合并到当前类型中。
2. 如果嵌入的类型和当前类型有同名字段，则会导致编译错误。
3. 嵌入类型本身不会作为一个独立字段存在。
4. 在语法中，嵌入类型只需使用类型名，省略字段名和类型声明。

#### 4.6 联合类型定义

通过 `oneof`
关键字，可以定义联合类型。联合类型允许字段具有多个可能的类型，在生成的代码中，联合类型会被转换为一个特殊的结构，包含一个类型标识字段以及多个可能的字段。在序列化时，只有被选择的那个字段会出现在输出中，类型标识字段则表明当前激活的选项。

例如：

```idl
oneof Value {
    User
    Manager
}
```

生成的结构体通常如下所示：

```
type Value struct {
    FieldType ValueType `json:"FieldType" form:"FieldType"` // 类型标识字段
    User      *User     `json:"User,omitempty" form:"User"` // 可能的选项1
    Manager   *Manager  `json:"Manager,omitempty" form:"Manager"` // 可能的选项2
}
```

在这个例子中，`FieldType` 是一个枚举类型，用于标识当前值的类型（例如，`User` 或 `Manager`）。在序列化时，只有被选择的字段会被包括在内，而其他字段会被排除。

#### 4.7 序列化和反序列化

为了支持可选字段、必选字段以及默认值填充，生成的代码需要支持流式解析。这种解析方式允许在处理请求时对必选字段进行校验，并为缺失的字段填充默认值。

- **必选字段（`required`）**：必选字段在解析和验证时必须存在且非空。
- **可选字段（`optional`）**：可选字段在数据中可以不存在，通常会被生成为空指针类型。
- **默认值（`compat_default`）**：为字段提供兼容性默认值，以确保字段在缺失时能够自动填充。

这种机制确保了数据结构在序列化和反序列化过程中能够正确处理缺失字段和赋予字段合适的默认值。

### 5. Validate 注解

`Validate` 注解用于对字段值进行验证，支持内置验证函数以及复杂的表达式。它提供了强大的数据校验机制，能够确保数据在运行时的有效性和一致性。

#### 5.1 验证函数

`Validate` 注解内置了多种常用的验证函数：

- `len`：用于验证字段的长度（支持字符串、列表、映射等）

此外，用户还可以定义自定义的验证函数。需要注意的是，自定义验证函数的名称不得与内置验证函数重名，但可以在多个字段中重复使用。生成的代码中会自动为每个自定义验证函数创建对应的函数模板，开发者可在这些模板中实现具体的验证逻辑。

#### 5.2 支持的操作符

验证表达式支持以下操作符：

- **比较操作符**：`==`、`!=`、`<`、`<=`、`>`、`>=`
- **逻辑操作符**：`&&`（与）、`||`（或）、`!`（非）
- **算术操作符**：`+`、`-`、`*`、`/`
- **函数调用**：支持调用内置函数（如 `len($)`）、内置正则函数（如 `email($)`）等

验证表达式遵循标准的运算符优先级规则，用户可以使用括号来明确指定运算顺序。

**重要提示**：在validate表达式中，字符串字面量必须使用单引号（`'`）包围，而不是双引号（`"`
）。注意：这仅适用于validate表达式内部的字符串，注解值（如json="name"）仍然使用双引号。

验证表达式的语法结构包括：

1. **原子表达式**：包括标识符、字面量（如整数、浮点数、字符串）以及特殊值 `$`（代表当前字段）和 `nil`
2. **函数调用**：支持带参数的函数调用，如 `len($)`、`email($)`、`regexp($, 'pattern')` 等
3. **一元操作符**：支持逻辑非操作符 `!`，例如 `!expr`（表示“非”）
4. **关系操作符**：包括 `<`、`<=`、`>`、`>=`，用于比较大小
5. **相等性操作符**：包括 `==` 和 `!=`，用于比较是否相等或不等
6. **逻辑与操作符**：`&&`，支持左结合运算
7. **逻辑或操作符**：`||`，支持左结合运算
8. **括号表达式**：通过圆括号 `()` 显式指定计算顺序

运算符的优先级从高到低依次为：

- **最高优先级**：原子表达式、函数调用、括号表达式
- **中等优先级**：一元操作符 `!`
- **关系操作符**：`<`, `<=`, `>`, `>=`
- **相等性操作符**：`==`, `!=`
- **最低优先级**：逻辑与 `&&`、逻辑或 `||`

### 5.3 特殊变量

- **`$`**：代表当前字段的值。
- 在复杂的验证表达式中，可以使用 **`$`** 引用当前字段来进行验证。

### 5.4 验证规则

- **可选字段（optional）**：验证仅在字段值存在时执行。
- **必需字段（required）**：验证会在字段值存在时执行。
- 如果验证失败，系统将抛出相应的错误信息。

### 5.5 使用示例

以下是验证规则的示例，展示了如何在字段上应用验证表达式：

```idl
type User {
    required string name (validate="$ != '' && len($) <= 64")  // 非空且长度不超过64
    int age (validate="$ >= 0 && $ <= 150")                    // 年龄范围 0-150
    string email (validate="email($)")                          // 邮箱格式验证
    list<string> tags (validate="len($) <= 10")                 // 标签数量不超过10个
    string phone (validate="phone($)")                          // 手机号格式验证
    string username (validate="regexp($, '^[a-zA-Z][a-zA-Z0-9_]{2,19}$')")  // 用户名格式验证
}
```

### 5.6 复杂验证表达式

你可以将多个验证条件组合在一起使用，进行更复杂的字段校验：

```idl
type Product {
    string code (validate="len($) >= 3 && len($) <= 20 && regexp($, '^\\w+$')")  // 3-20位字母、数字或下划线
    int price (validate="$ > 0 && $ <= 999999")                                  // 价格范围验证
    list<string> images (validate="len($) >= 1 && len($) <= 10")                 // 图片数量限制在1-10个之间
}
```

### 5.7 验证函数收集机制

系统会自动收集所有使用到的验证函数，包括内置的和自定义的函数。对于自定义验证函数，系统将在生成的代码中为其创建占位符，开发者需要在
`validate.go` 文件中实现具体的验证逻辑。

自定义验证函数的签名格式应为：
`func (参数类型) bool`

- 返回 `true` 表示验证通过
- 返回 `false` 表示验证失败

### 6. RPC 接口定义

RPC 接口定义分为两种形式：普通接口和 SSE（Server-Sent Events）流式接口。

#### 6.1 多文件项目结构

一个项目通常包含多个 IDL 文件（扩展名为 `.idl`）和一个 `meta.json` 配置文件。
系统通过 `ParseDir` 函数扫描整个目录，加载所有的 IDL 文件，并构建统一的项目结构。
每个 IDL 文件可以定义自己的常量、枚举、类型以及 RPC 接口，并且这些定义可以在其他 IDL 文件中引用。

#### 6.2 跨文件类型引用

系统支持在一个 IDL 文件中引用其他文件中定义的类型。
当解析器遇到当前文件中未定义的类型时，它会在项目中的其他文件中查找该类型的定义。如果无法找到该类型，系统会抛出错误，提示
`"type X is used but not defined"`。

#### 6.3 错误码扩展机制

错误码扩展（`enum extends`）功能允许跨文件进行操作。
当一个文件中的 `enum` 定义使用 `extends` 关键字扩展另一个文件中的 `enum` 时，系统会将扩展的字段与原 `enum` 定义合并。
在合并过程中，系统会检查字段名和值的唯一性，防止冲突。

#### 6.4 项目解析流程

项目的解析过程遵循以下步骤：

1. **加载阶段**：扫描目录中的所有 `.idl` 文件和 `meta.json` 配置文件。
2. **解析阶段**：逐一解析每个 IDL 文件，并将其转化为 `Document` 结构。
3. **验证阶段**：检查项目中所有类型的引用是否有效。
4. **合并阶段**：处理错误码扩展，将扩展的字段合并到基础的 `enum` 中。
5. **处理阶段**：处理泛型实例化、嵌入类型以及验证表达式等。
6. **路径处理阶段**：验证 RPC 路径参数是否与请求类型正确绑定。

#### 6.5 验证函数跨文件处理

系统会收集项目中所有 IDL 文件中定义的自定义验证函数，确保同名的验证函数在相同类型的字段上被正确应用。
这些验证函数会在生成的 `validate.go` 文件中创建占位符，供开发者实现具体的验证逻辑。

#### 6.6 项目配置文件

项目的根目录必须包含一个 `meta.json` 配置文件，该文件用于存储项目的元数据。
如果缺少此文件，项目将无法被正确解析。
`meta.json` 配置文件通常包含项目名称、版本号以及其他必要的元数据信息。

### 6.7 普通接口定义

普通接口的定义格式如下：

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

### 6.8 SSE流式接口定义

SSE（Server-Sent Events）流式接口是现代应用开发中不可或缺的一项功能，广泛用于实时数据推送和长连接通信。SSE接口的定义使用
`sse` 关键字，如下所示：

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

1. **实时数据推送**：如股票价格的变动、聊天消息的实时传输等。
2. **AI模型结果流式返回**：逐步返回AI生成的内容，特别是在需要长时间计算的场景中，例如生成图片、文本内容等。

#### 6.9 HTTP相关注解

在RPC接口中，可以使用注解来指定HTTP方法、路径、超时时间等路由信息。常见的HTTP相关注解包括：

- `method` - 指定HTTP方法（例如：GET、POST、PUT、DELETE等）
- `path` - 请求路径，支持RESTful路径参数
- `contentType` - 请求内容类型，支持 "form"（表单编码）或 "json"（JSON编码）
- `connTimeout` - 连接超时，单位毫秒
- `readTimeout` - 读取超时，单位毫秒
- `writeTimeout` - 写入超时，单位毫秒
- `resp.go.type` - 指定RPC响应的Go类型

示例：

```
rpc GetUser (GetUserRequest) GetUserResponse {
    method = "GET"
    path = "/user/:id"
    contentType = "json"
    connTimeout = "100"
    readTimeout = "300"
    writeTimeout = "300"
    summary = "获取用户信息"
}
```

#### 6.10 RESTful 路径规则与参数绑定

RESTful路径参数用于定义动态路径，支持两种风格：

1. **冒号风格（Colon Style）**：`:paramName` 或 `:paramName*`
2. **大括号风格（Brace Style）**：`{paramName}` 或 `{paramName...}`

路径参数的类型包括：

- **静态段（Static）**：固定的路径段，如 `users`、`books` 等
- **普通参数（Param）**：`:id` 或 `{id}` - 匹配单个路径段
- **通配符参数（Wildcard）**：`:path*` 或 `{path...}` - 匹配多个路径段

**路径排序规则**：

服务器会自动对路由路径进行排序，确保更具体的路径优先匹配。排序规则如下：

- 静态路径段优先于参数路径段，参数路径段优先于通配符路径段（静态段 < 参数段 < 通配符段）
- 静态段按字典顺序排序
- 路径长度较长的优先于较短的
- 在相同长度的情况下，路径按字典顺序排序

### 路径匹配优先级

路径匹配规则决定了在多个匹配项中，哪个路径会首先被匹配。以下是不同路径匹配类型的优先级顺序：

```
// 排序后的匹配优先级：
GET /user/profile     # 最具体，优先级最高
GET /user/:id         # 参数化路径，优先级中等
GET /files/:path*     # 通配符路径，优先级最低
```

### 路径参数命名规则

在路径中使用的参数必须遵循以下命名规则：

- 参数名必须以字母开头，不能以数字开头。
- 支持字母、数字、下划线 (`_`) 和连字符 (`-`) 的组合，例如：`:user_id`, `:user-name`。

#### 示例：

```
// 冒号风格路径参数
rpc GetUser (GetUserRequest) GetUserResponse {
    method = "GET"
    path = "/user/:id"  // 路径参数 :id
}

rpc GetFile (GetFileRequest) GetFileResponse {
    method = "GET"
    path = "/files/:path*"  // 通配符路径参数 :path*
}

// 大括号风格路径参数
rpc GetUser (GetUserRequest) GetUserResponse {
    method = "GET"
    path = "/user/{id}"  // 路径参数 {id}
}

rpc GetFile (GetFileRequest) GetFileResponse {
    method = "GET"
    path = "/files/{path...}"  // 通配符路径参数 {path...}
}

// 混合路径示例
rpc ComplexPath (ComplexPathRequest) ComplexPathResponse {
    method = "GET"
    path = "/org/{orgId}/repos/:repoId/branches/{branch...}"  // 混合参数
}

// 静态路径段示例
rpc StaticPath (StaticPathRequest) StaticPathResponse {
    method = "GET"
    path = "/api/v1/users/:id/settings"  // 包含多个静态路径段
}
```

### 参数绑定

在定义请求类型时，可以使用注解将字段与路径参数或查询参数进行关联：

```
// 请求类型示例
type GetUserRequest {
    string userId (path="id")  // 将 userId 字段绑定到路径参数 :id 或 {id}
    string locale (query="locale")  // 将 locale 字段绑定到查询参数 ?locale=value
}
```

#### 注意事项：

- **路径参数**：绑定到路径中的参数必须是必需的 (`required`)，且其类型通常为基本类型，如 `int` 或 `string`。
- **查询参数**：查询参数通常是可选的 (`optional`)，可以为字段设置默认值。它们非常适合用于传递过滤条件、分页信息等。

### 路径参数命名规则总结

- 路径参数必须与路径中的参数名保持一致。
- 在路径中，**静态段**优先于**参数化路径段**，而**通配符参数**的优先级最低。
- 参数命名必须遵循字母开头、支持字母、数字、下划线和连字符的规则。

## 示例

以下是一个完整的电商系统的IDL文件示例：

```
// 核心功能接口定义

// 订单状态枚举
enum OrderStatus {
    PENDING_PAYMENT = 1 (desc="待支付")
    PAID = 2 (desc="已支付")
}

// 错误码定义
enum ErrCode {
    ERR_OK = 0 (errmsg="success")
    PARAM_ERROR = 1003 (errmsg="parameter error")
}

// 扩展错误码
enum extends ErrCode {
    USER_NOT_FOUND = 404 (errmsg="user not found")
}

// 用户信息结构
type User {
    required string id
    required string username (validate="$ != '' && len($) >= 3")
    required string email (validate="email($)")
}

// 商品信息结构
type Product {
    required string id
    required string name
    float price (validate="$ > 0")
}

// 订单信息结构
type Order {
    required string id
    string userId (json="user_id")
    list<OrderItem> items
    float totalPrice (json="total_price")
    OrderStatus status
}

// 订单项
type OrderItem {
    string productId (json="product_id")
    float price
    int quantity
}

// 用户注册请求
type RegisterRequest {
    required string username (validate="$ != '' && len($) >= 3")
    required string email (validate="email($)")
    required string password (validate="len($) >= 6")
}

// 用户登录请求
type LoginRequest {
    required string email (validate="email($)")
    required string password (validate="len($) >= 6")
}

// 获取用户信息请求
type GetUserInfoRequest {
    required string userId (path="id", validate="$ != ''")
}

// 获取商品详情请求
type GetProductDetailRequest {
    required string productId (path="id", validate="$ != ''")
}

// 用户注册响应包装
type RegisterResponse {
    int code
    string message
    User data
}

// 用户登录响应包装
type LoginResponse {
    int code
    string message
    User data
}

// 获取用户信息响应包装
type GetUserInfoResponse {
    int code
    string message
    User data
}

// 获取商品详情响应包装
type GetProductDetailResponse {
    int code
    string message
    Product data
}

// 用户服务接口
rpc Register (RegisterRequest) RegisterResponse {
    method = "POST"
    path = "/auth/register"
    contentType = "json"
    connTimeout = "100"
    readTimeout = "300"
    writeTimeout = "300"
    summary = "用户注册"
}

rpc Login (LoginRequest) LoginResponseWrapper {
    method = "POST"
    path = "/auth/login"
    contentType = "json"
    connTimeout = "100"
    readTimeout = "300"
    writeTimeout = "300"
    summary = "用户登录"
}

rpc GetUserInfo (GetUserInfoRequest) GetUserInfoResponseWrapper {
    method = "GET"
    path = "/user/:id"
    contentType = "json"
    connTimeout = "100"
    readTimeout = "300"
    writeTimeout = "300"
    summary = "获取用户信息"
}

// 商品服务接口
rpc GetProducts (GetProductsRequest) GetProductsResponse {
    method = "GET"
    path = "/products"
    contentType = "json"
    connTimeout = "100"
    readTimeout = "300"
    writeTimeout = "300"
    summary = "获取商品列表"
}

// 商品分页结果类型
type ProductPageResult {
    list<Product> items
    int total
    int page
    int size
}

// 获取商品列表响应
type GetProductsResponse {
    int code
    string message
    ProductPageResult data
}

rpc GetProductDetail (GetProductDetailRequest) GetProductDetailResponse {
    method = "GET"
    path = "/product/:id"
    contentType = "json"
    connTimeout = "100"
    readTimeout = "300"
    writeTimeout = "300"
    summary = "获取商品详情"
}

// 购物车服务接口
rpc AddToCart (AddToCartRequest) AddToCartResponse {
    method = "POST"
    path = "/cart/add"
    contentType = "json"
    connTimeout = "100"
    readTimeout = "300"
    writeTimeout = "300"
    summary = "添加商品到购物车"
}

// 订单服务接口
rpc CreateOrder (CreateOrderRequest) CreateOrderResponseWrapper {
    method = "POST"
    path = "/order/create"
    contentType = "json"
    connTimeout = "100"
    readTimeout = "300"
    writeTimeout = "300"
    summary = "创建订单"
}

// 获取订单详情响应包装
type GetOrderDetailResponse {
    int code
    string message
    Order data
}

rpc GetOrderDetail (GetOrderDetailRequest) GetOrderDetailResponse {
    method = "GET"
    path = "/order/:id"
    contentType = "json"
    connTimeout = "100"
    readTimeout = "300"
    writeTimeout = "300"
    summary = "获取订单详情"
}

// 支付回调响应包装
type PaymentNotifyResponse {
    int code
    string message
    // 空数据或成功信息
}

// 支付回调接口
rpc PaymentNotify (PaymentNotifyRequest) PaymentNotifyResponse {
    method = "POST"
    path = "/payment/notify"
    contentType = "form"
    connTimeout = "100"
    readTimeout = "300"
    writeTimeout = "300"
    summary = "支付结果通知回调"
}

// SSE流式接口 - 订单状态变更通知
// SSE流式响应包装
type OrderStatusChangeResponse {
    int code
    string message
    Order data
}

sse OrderStatusChange (GetOrderDetailRequest) OrderStatusChangeResponse {
    method = "GET"
    path = "/order/:id/stream"
    contentType = "text/event-stream"
    connTimeout = "100"
    readTimeout = "300"
    writeTimeout = "300"
    summary = "订单状态变更实时推送"
}
```
