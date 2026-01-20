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

* `name`: 项目名称，通常用于生成代码中的包名或命名空间
* `version`: 项目的版本号
* `description`: 项目的简要描述

### IDL 文件

IDL 文件以 `.idl` 扩展名命名，用于定义数据结构、接口和服务。一个项目可以包含多个 IDL
文件，所有文件共享同一个命名空间，可以相互引用其他文件中定义的类型、枚举和常量。

#### 命名空间与类型引用

* 所有 IDL 文件共享同一个全局命名空间
* 在任何 IDL 文件中定义的类型、枚举和常量，可以在其他文件中直接引用
* 类型名称在项目中必须唯一，不能在不同文件中定义相同的类型
* 系统会自动解析不同文件之间的类型引用关系

#### 文件组织建议

尽管 IDL 文件可以互相引用，但建议按功能将文件进行组织，常见的组织方式包括：

* **通用类型文件**：定义项目中通用的基础类型和枚举
* **业务模块文件**：按照业务领域划分，如用户模块、订单模块等
* **服务接口文件**：定义 RPC 接口及其请求和响应类型
* **错误码文件**：集中管理项目中的错误码枚举

## IDL 文件构成

一个 IDL 文件由以下几个主要部分组成：

* **注释**：用于功能说明和文档描述
* **关键字**：IDL 语言的保留字
* **常量定义**：定义固定不变的常量值
* **枚举定义**：定义枚举类型，通常用于状态码或选项
* **类型定义**：定义结构体、联合类型等复合数据类型
* **RPC 接口定义**：定义服务接口，支持普通 RPC 方法和流式接口

语句以换行符分隔，空行将被忽略。建议使用 4 个空格进行缩进，避免使用 Tab 键。

## 基本语法规则

### 注释

* **单行注释**：使用 `//` 或 `#` 开头
* **多行注释**：使用 `/* ... */` 包裹

### 空白字符处理

* 换行符、空格和制表符均视为空白字符
* 允许多个连续的换行符
* 文档可以以任意数量的换行符开始或结束
* 定义之间可以使用多个换行符进行分隔

### 关键字

以下是 IDL 语言的保留关键字，无法作为标识符使用：

* `extends` - 扩展错误码
* `const` - 定义常量
* `enum` - 定义枚举
* `type` - 定义类型
* `oneof` - 定义联合类型
* `rpc` - 定义普通 RPC 方法
* `sse` - 定义流式 RPC 方法
* `true` / `false` - 布尔值
* `optional` - 可选字段标记
* `required` - 必填字段标记

### 标识符规则

* 标识符必须以字母开头
* 可以包含字母、数字、下划线 `_` 和点号 `.`
* 区分大小写
* 推荐在以下场景使用首字母大写：

    * **常量名**：如 `MAX_SIZE`
    * **枚举名**：如 `Status`
    * **枚举字段**：如 `ACTIVE`
    * **类型名**：如 `User`
    * **方法名**：如 `GetUser`

**注意**：标识符也可以作为常量值或注解值。

### 格式与布局规则

* **缩进**：使用 4 个空格进行缩进，避免使用 Tab 键
* **换行**：每条语句使用换行符分隔，每行单独一条语句
* **空行**：空行用于分隔逻辑块，提升可读性
* **大括号**：

    * 左大括号 `{` 紧跟在定义后，不换行
    * 右大括号 `}` 独占一行

### 基础类型

* **`bool`**：布尔类型
* **`int`**：整数类型，映射为 Go 的 `int64`（默认使用此类型）。如果需要自定义类型，可以使用 `go.type` 注解。
* **`float`**：浮点数类型，映射为 Go 的 `float64`（默认使用此类型）。同样，可以通过 `go.type` 注解进行自定义。
* **`string`**：字符串类型，仅支持双引号（单引号不合法）。
* **`bytes`**：字节数组类型，传输时使用 Base64 编码为字符串。

### 字面量格式

* **整数**：`42`，`-17`，`0x1A2B`
* **浮点数**：`3.14`，`.5`，`-2.7e10`
* **字符串**：`"hello"`，`"escaped \" quote"`（仅支持双引号）
* **布尔值**：`true`，`false`

### 容器类型

* **`map<K, V>`**：映射类型，键类型必须是 `int` 或 `string`，值类型可以是任意类型。
* **`list<T>`**：列表类型，元素类型可以是任意类型。
* **`bytes`**：字节数组类型（传输时以 Base64 编码为字符串）。

**注意**：`map` 的键类型仅支持 `int` 或 `string`，其他类型会在解析时抛出错误。

`map` 和 `list` 都支持嵌套，例如：`list<list<int>>`，`map<string,map<string,int>>` 等。

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

* **整数字面量**：`42`, `-17`, `0x1A2B`
* **浮点数字面量**：`3.14`, `.5`, `-2.7e10`
* **字符串字面量**：`"hello"`, `"escaped \" quote"`
* **布尔值**：`true`, `false`
* **标识符**：通常用于引用枚举成员

**注意**：在注解中，你可以使用枚举作为常量值。

### 2. 注解（Annotations）

注解（Annotations）是IDL（接口定义语言）中的重要组成部分，主要用于为字段、类型、枚举和接口添加元数据。
注解为生成的代码行为提供灵活的控制机制，可以指定序列化规则、验证规则、API路由等。

注解的定义支持单行和多行两种格式：

* **单行定义**：适用于字段，格式为：

```
(注解名 = 注解值, 注解名 = 注解值, ...)
```

* **多行定义**：可以在字段、类型或接口中使用，示例如下：

```idl
# 类型定义中的注解
type User {
    string name (
      json="name"
      go.type="string"
    )
    int age (json="age,omitempty")
}

# 接口定义中的注解
rpc GetUser (GetUserRequest) GetUserResponse {
    method = "GET"
    path = "/user/:id"
}
```

注解可在以下场景中使用：

* **枚举定义**：为错误码等枚举项添加描述信息；
* **字段定义**：指定字段的JSON序列化规则、验证规则、参数绑定等；
* **RPC接口定义**：指定HTTP方法、路径、超时时间等路由信息。

#### 注解语法细节：

* **注解值类型**：注解的值可以是布尔值、整数、浮点数、字符串或标识符；
* **分隔符**：多个注解之间可以使用逗号或换行符进行分隔；
* **标识符形式**：注解可以仅包含标识符，不必是键值对的形式，例如：`(some_flag)`；
* **完整语法**：注解的完整语法为：`<标识符>(=<常量值>)?`。

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
    RED = 1
    GREEN = 2
    BLUE = 3
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
    ENGINEERING = 1
    MARKETING = 2
    SALES = 3
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

#### 4.2 哈希键冲突检查

系统会对每个类型中的字段进行哈希键冲突检查，以确保生成的代码能够正确工作。具体来说：

- 每个字段的JSON标签和表单标签都会生成唯一的哈希键（使用FNV1a64算法）
- 系统会检查同一类型中是否存在相同的哈希键
- 如果发现重复的哈希键，系统会报告错误："type X has duplicate hash key for field Y and Z"
- 这种检查确保了基于哈希的字段调度机制能够正常工作，提高了JSON解析性能

#### 4.3 泛型结构体定义

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

#### 4.4 泛型实例化

对已定义的泛型进行实例化，注意语法格式：

```
type UserResponse Response<User>
```

语法格式：

```
type <新类型名> <已定义的类型名><<具体类型>>
```

#### 4.5 结构体字段

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

- `json` - 指定JSON序列化字段名和选项，支持 `non-omitempty` 选项禁用omitempty行为
- `enum_as_string` - 枚举作为字符串处理
- `form` - 指定表单序列化字段名和选项

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
    string description (json="desc,non-omitempty")  // 禁用omitempty行为，即使为空也会序列化
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

验证表达式支持的完整语法包括：

1. **原子表达式**：标识符、字面量（整数、浮点数、字符串）、特殊值 `$`（当前字段）和 `nil`
2. **函数调用**：支持带参数的函数调用，如 `len($)`、`email(value)` 等
3. **一元操作符**：支持逻辑非操作 `!expr`
4. **关系操作符**：`<`, `<=`, `>`, `>=`
5. **相等性操作符**：`==`, `!=`
6. **逻辑与操作符**：`&&` （支持左结合）
7. **逻辑或操作符**：`||` （支持左结合）
8. **括号表达式**：使用圆括号 `()` 来明确指定计算顺序

运算符优先级从高到低依次为：

- 原子表达式、函数调用、括号表达式
- 一元操作符 `!`
- 关系操作符 `<`, `<=`, `>`, `>=`
- 相等性操作符 `==`, `!=`
- 逻辑与 `&&`
- 逻辑或 `||`

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

#### 5.7 验证函数收集机制

系统会自动收集所有使用的验证函数，包括内置函数和自定义函数。
对于自定义验证函数，系统会在生成代码时创建对应函数的占位符，
开发者需要在生成的 `validate.go` 文件中实现具体的验证逻辑。

自定义验证函数的签名格式为：`func (参数类型) bool`，
返回 `true` 表示验证通过，`false` 表示验证失败。

### 6. RPC 接口定义

RPC接口定义有两种形式：普通接口和SSE（Server-Sent Events）流式接口。

#### 6.1 多文件项目结构

在一个项目中，可以有多个IDL文件（.idl扩展名）和一个meta.json配置文件。
系统通过 `ParseDir` 函数扫描整个目录，加载所有IDL文件并构建统一的项目结构。
每个IDL文件可以定义自己的常量、枚举、类型和RPC接口，这些定义可以在其他IDL文件中被引用。

#### 6.2 跨文件类型引用

系统支持在任一IDL文件中引用其他文件中定义的类型。
当解析器遇到未在当前文件中定义的类型时，会搜索项目中的其他文件寻找相应定义。
如果找不到，则会报错："type X is used but not defined"。

#### 6.3 错误码扩展机制

错误码扩展（enum extends）功能支持跨文件操作。
当一个文件中的enum定义使用 `extends` 关键字扩展另一个文件中定义的enum时，
系统会将扩展的字段合并到原enum定义中。合并过程中会检查字段名和值的唯一性，防止冲突。

#### 6.4 项目解析流程

系统按以下步骤解析多文件项目：

1. **加载阶段**：扫描目录中的所有 `.idl` 文件和 `meta.json` 文件
2. **解析阶段**：逐个解析IDL文件为Document结构
3. **验证阶段**：检查所有类型引用是否有效
4. **合并阶段**：处理错误码扩展，将扩展的枚举字段合并到基枚举
5. **处理阶段**：处理泛型实例化、嵌入类型、验证表达式等
6. **路径处理阶段**：验证RPC路径参数与请求类型绑定的一致性

#### 6.5 验证函数跨文件处理

系统会收集所有IDL文件中定义的自定义验证函数，并确保相同名称的验证函数应用于相同类型的字段。
验证函数会在生成的 `validate.go` 文件中创建占位符，供开发者实现具体的验证逻辑。

#### 6.6 项目配置文件

项目根目录需要包含 `meta.json` 配置文件，用于存储项目的元信息。
此文件是必需的，没有此文件的项目将无法正确解析。
配置文件通常包含项目名称、版本和其他元数据信息。

#### 6.7 普通接口定义

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

#### 6.8 SSE流式接口定义

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

#### 6.9 HTTP相关注解

RPC接口可以使用注解来指定HTTP方法、路径、超时时间等路由信息，常见的HTTP相关注解包括：

- `method` - HTTP方法（GET、POST、PUT、DELETE等）
- `path` - 请求路径，支持RESTful路径参数
- `contentType` - 内容类型，支持 "form"（表单编码）或 "json"（JSON编码）
- `connTimeout` - 连接超时
- `readTimeout` - 读取超时
- `writeTimeout` - 写入超时
- `summary` - 接口摘要说明
- `resp.go.type` - 指定RPC响应的Go类型

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

#### 6.10 RESTful Path 规则和参数绑定

RESTful路径参数用于定义动态路径，支持两种风格：

1. 冒号风格（Colon Style）：`:paramName` 或 `:paramName*`
2. 大括号风格（Brace Style）：`{paramName}` 或 `{paramName...}`

路径参数类型：

- 静态段（Static）：固定的路径段，如 `users`, `books` 等
- 普通参数（Param）：`:id` 或 `{id}` - 匹配单个路径段
- 通配符参数（Wildcard）：`:path*` 或 `{path...}` - 匹配多个路径段

路径排序机制：服务器端会自动对路由路径进行排序，确保更具体的路径优先匹配。排序规则如下：

- 静态段优先于参数段，参数段优先于通配符段（Static(0) < Param(1) < Wildcard(2)）
- 静态段按字典序排序
- 路径长度较长的优先于较短的
- 相同长度的路径按路径字符串排序

示例：

```
// 排序后的匹配优先级：
GET /user/profile     # 最具体，优先级最高
GET /user/:id         # 参数路径，优先级中等
GET /files/:path*     # 通配符路径，优先级最低
```

路径参数命名规则：

- 参数名必须以字母开头，不能以数字开头
- 支持字母、数字、下划线和连字符，如 `:user_id`, `:user-name`

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

// 混合路径示例
rpc ComplexPath (ComplexPathRequest) ComplexPathResponse {
    method = "GET"
    path = "/org/{orgId}/repos/:repoId/branches/{branch...}"  // 混合参数
}

// 静态路径段示例
rpc StaticPath (StaticPathRequest) StaticPathResponse {
    method = "GET"
    path = "/api/v1/users/:id/settings"  // 包含多个静态段
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

// 表单提交接口示例
type FormData {
    string name (form="username", validate="$ != ''")
    int age (form="user_age", validate="$ > 0")
    string email (form="user_email", validate="email($)")
}

rpc SubmitForm (FormData) Response {
    method = "POST"
    path = "/form/submit"
    contentType = "form"  // 使用表单编码
    summary = "提交表单数据"
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
15. 多文件项目必须包含 meta.json 配置文件，否则解析会失败
16. 跨文件引用的类型必须确保唯一性，不能有重复的类型名称
17. 不同文件中的验证函数如果名称相同，则必须应用于相同类型的字段
18. 在多文件环境中，错误码扩展功能可在不同文件间合并枚举定义
19. 每个IDL文件应专注于定义相关的数据结构和接口，保持良好的模块化设计
20. 系统会自动为每个字段生成基于字段名的哈希键，用于高性能JSON解析
21. 同一类型中的字段名不能过于相似以至于产生相同的哈希键，否则会导致解析错误
22. compat_default 注解只能用于 required 字段，用于提供向后兼容的默认值
