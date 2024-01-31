# 02-HelloWeb

- 详细介绍
    -  如何运行
- 扩展内容
    -  Types of imports
        - import "path/to/package"
        - import \<alias> "path/to/package"
        - import 其他用法
    -  Type declarations
        - Type declarations
        - Client 对应的函数为什么叫 method，而不是 function？
            - 自定义类型的的方法集
    -  Function declarations
        -  任意参数的函数
        -  省略掉参数标识符的函数
    -  Method declarations
    - FQA
        - Q：为什么要设计接收器函数 `func (receiverType) methodName` 与普通函数 `func methodName`？
        - Q：自定义类型 T 的值接收器方法与指针接收器方法有什么区别？
            - 接收器方法集遵守的规则
        - Q：为什么 Go 语言把类型放在后面？

## 详细介绍

### 如何运行

在文件目录下运行 `$ go run helloweb.go`，访问 GET：<http://localhost:1234/helloweb>，或者 `$ curl localhost:1234/helloweb`。

## 扩展内容

### Types of imports

Go 程序是通过 import 关键字将一组包链接在一起。包的导入形式有五种：<b>单行或多行直接导入、别名导入、点语法导入、占位符导入和嵌套结构导入</b>。详细示例见官方文档 [Import in GoLang](https://golangdocs.com/import-in-golang)。

#### import "path/to/package"

`import "xx"` 的意思是 `import path/to/package`，指导入 package 所在的目录的路径，而不是 package 包名，所以 package 和 package 所在的目录名也可以不一致。如 `import "fmt"` 指标准库 package fmt 所在的路径 `"$GOROOT/src/fmt/"`，而不是指 package fmt 本身包名（只是默认规范）。

```go
import (               
   "fmt"              // $GOROOT/src 标准库包路径
   "my/testpackage"   // $GOPATH/my/testpackage 第三方依赖包
   "./api"            // ./api
)
```

#### import \<alias> "path/to/package"

`import alias "path/to/package"`，是指使用别名代替包的真实名称。如果引用多个（不同路径但）重名的包，别名能避免冲突。

```go
import f "fmt"

func main() {
   f.Println("Hello, World")        // Hello, World
}
```

#### import 其他用法

<b>`import _ "os"`</b> 空占位符用法（blank import）。在 Go 程序中，导入的包必须要被使用，否则编译时会报 `unused import error` 错误。blank import 能让编译器忽略该包未使用（仅执行包的初始化函数）。注意：当使用该包时，要去掉该占位符 `_`，否则编译器会报 undefined。

<b>`import "math/rand"`</b> 嵌套结构，仅导入大包中真正被使用的子包，如 `rand.Int()`。<b>`import . "fmt"`</b> 点语法导入，是将 fmt 包的 namespace 合并到当前程序  中，这样在调用导入包方法时（`fmt.Println("hello")`），不需要使用 package name，而是直接调用（ `Println("xx")`）。这样做的缺点：可能导致命名空间冲突。该方法是 `import <alias> "path/to/package"` 导入的省略方式。

### Type declarations


~~~go
type Client struct {
    Transport RoundTripper
    CheckRedirect func(req *Request, via []*Request) error
    Jar CookieJar
    Timeout time.Duration // Go 1.3
}

func (c *Client) Do(req *Request) (*Response, error)
~~~

net/http 包中 type [Client](http://docscn.studygolang.com/pkg/net/http/#Client) 声明如上：

- type 是什么类型？如何使用？
- Client 对应的函数为什么叫 method，而不是 function？


#### Type declarations

类型声明有两种形式（<b>别名声明是有赋值表达式，类型定义是为了复用</b>）：

- 别名声明（Alias declarations）：标识符与给定的类型进行绑定，在标识符的作用域范围内的，它就是类型的别名。ENBF：`AliasDecl = identifier "=" Type .`
- 类型定义（[Type definitions](http://docscn.studygolang.com/ref/spec#Type_declarations)）：创建一个与标识符绑定的新类型，它不同于那些给定的底层类型和操作，这种新类型称为自定义类型（defined type）。ENBF：`TypeDef = identifier Type .`

别名声明示例：

~~~go
type (
    nodeList = []*Node  // nodeList and []*Node are identical types
    Polar    = polar    // Polar and polar denote identical types
)
~~~

类型定义的几种形式示例：

~~~go
type ConnState int
type (                          // 多列合并的书写方式，拆成单行与 ConnState 形式一致
    Point struct{ x, y float64 }  // Point and struct{ x, y float64 } are different types
    polar Point                   // polar and Point denote different types
)
type Client struct {            // 结构体是最常用的复合类型，还包括数组、切片、映射和函数
    Transport RoundTripper
    CheckRedirect func(req *Request, via []*Request) error
    Jar CookieJar
    Timeout time.Duration // Go 1.3
}
type CloseNotifier interface {  // 接口类型：一组方法签名的集合
    CloseNotify() <-chan bool
}
~~~

上述几种类型定义的形式，需要与没有 `type` 关键字的定义形式进行辨析：它们之间的主要区别在于<b>作用域</b>。使用 `type` 关键字声明的具有以下特点（对应没有关键字的声明）：

* 可以在包的任何地方使用。（只能在声明它的函数或方法中使用）
* 可以被其他包导入。（不能被其他包导入）
* 可以被用作其他类型的字段类型。（不能被用作其他类型的字段类型）

这些使用 `type` 关键字声明的类型，称为<b>已定义类型（defined type）</b>，即<b>用户自定义类型</b>。其中特别说明一下 `type ConnState int` 声明中的 `ConnState` 类型，与 `int` 类型的关系：

- 使用<b>类型断言</b> `x.(T)` 可知两者属于不同的类型；
- `ConnState` 的<b>底层类型（underlying type）</b>是 `int`；
- `ConnState` 变量和 `int` 变量的值是相同的，使用类型断言可以将 `ConnState` 变量转换为 `int` 类型，但不可将 `int`->`ConnState`。

~~~go
// 创建一个 ConnState 变量
state := ConnState(net.DialTCP("tcp", nil, nil))
// 使用类型断言将 ConnState 变量转换为 int 类型
var assertedIntValue int = state.(int)
~~~

#### Client 对应的函数为什么叫 method，而不是 function？

自定义类型（defined type）可以有相应的<b>关联方法（Method）</b>，也称方法集（method set）。

~~~go
type Mutex struct         { /* Mutex fields */ }

// 两个关联方法，也称方法集（method set）
func (m *Mutex) Lock()    { /* Lock implementation */ }
func (m *Mutex) Unlock()  { /* Unlock implementation */ }
~~~

##### 自定义类型的的方法集

自定义类型的的方法集将不会被继承，除非其被作为<b>复合类型的元素</b>，或是接口类型时。（复合类型包括数组、切片、映射和函数，而结构体是最常见的）

~~~go
// 不继承绑定到给定类型，即这两个变量无法调用方法集
type NewMutex Mutex
type PtrMutex *Mutex

// 复合类型会继承自定义类型的方法
type PrintableMutex struct {  // struct 是常见的复合类型
    Mutex   // 作为 embedded field 元素
}

type Block interface {
    BlockSize() int
    Encrypt(src, dst []byte)
    Decrypt(src, dst []byte)
}
type MyBlock Block  // 接口类型 MyBlock 会继承自定义类型 Block 的方法集
~~~

### Function declarations

Go 函数名必须以字母或下划线开头，后面可以跟任意数量的字母、数字或下划线，且不能与 Go 语言的关键字相同。具体格式如：`func 函数名称(参数列表) 返回值列表`。

- 返回值类型可以是任何类型，包括基本类型、复合类型、指针类型、函数类型等。如果没有返回值，则返回值类型可以省略；
- 如果给返回类型进行了命名，return 语句可以为空（注意不是省略），且被返回类型所命令的变量不要重复定义；
- 函数可以嵌套定义，即在一个函数体内定义另一个函数。即匿名函数，直接调用或赋予给变量。也称闭包 closures，[Function literals](http://docscn.studygolang.com/ref/spec#Function_literals)；
- 函数可以重载，即同一个函数名可以对应不同的参数列表和返回值类型。
- 形参为 `...` 的函数称为 variadic，表示可以传入 0 到多个参数（官方文档 [Passing arguments to ... parameters](http://docscn.studygolang.com/ref/spec#Passing_arguments_to_..._parameters)）。
- 若函数的参数，在函数体内不存在引用，则可以省略掉参数标识符。

函数的类型声明（[Function types](http://docscn.studygolang.com/ref/spec#Function_types)），示例如下：

~~~go
func()                                                        // 无参，无返回类型的闭包
func(x int) int
func(a, _ int, z float32) bool                                // 多个参数，int 类型有 a 和占位标识符，
func(a, b int, z float32) (bool)                              // 返回类型为 bool
func(prefix string, values ...int)                            // 第二个形参 values 表示可以传入多个 int 类型的参数
func(a, b int, z float64, opt ...interface{}) (success bool)  // 第四个形参 opt 表示可以传入多个 interface{} 类型的参数，返回类型为 bool，且为其命名为 success
func(int, int, float64) (float64, *[]int)                     // 返回类型有两个，且使用括号（块语句）将它们圈在一起
func(n int) func(p *T)                                        // 返回类型为函数

f := func(x, y int) int { return x + y }  // 匿名函数，直接调用或赋予给变量
~~~

若函数签名（function's signature）中声明了结果参数，则函数体的语句列表必须以终止语言 return 结束，即<b>声明了函数返回值</b>。函数声明是可以省略函数体，这样的声明为 Go 在外部实现函数提供了签名，例如在汇编（assembly routine）程序中实现。详细见 [Function_declarations](http://docscn.studygolang.com/ref/spec#Function_declarations)

```go
func min(x int, y int) int {  // function's signature, int 即为 result parameters
    if x < y {
        return x
    }
    return y                  // terminating statement
}

// 没有函数体，implemented externally
func flushICache(begin, end uintptr)  // function's signature
```

#### 任意参数的函数

~~~go
func sum(s string, args ...int)  {
    var x int
    for _, n := range args {  // _ 为将索引忽略
        x += n
    }
    fmt.Println(s, x)
}
~~~

##### 省略掉参数标识符的函数

例如函数签名为 `func f(int)`，该函数的形参类型为 int，但是却省略了名称，这是为什么？若函数的参数，在函数体内不存在引用，则可以省略掉参数标识符。
那既然是不需要，那为什么还要在函数签名中定义呢？某些函数需要使用统一的函数类型，但不是所有的函数都需要对应的参数。如果在函数内部不使用这个参数，却在形参中定义了这个形参变量的话，编译的时候会提示变量未使用。在编译的时候只是检查形参的类型，所以定义函数的时候只需指定形参的类型就可以了，可以省去形参名。例如各种窗口消息的响应函数，函数要定义 WPARAM、LPARAM 类型的两个参数，但不是每个消息都要用到这两个参数的，例如 WM_CLOSE 消息，这是就可以把参数只写类型而不写名称。

示例如下：

~~~go
type IWriter interface {
    Write(p []byte) error
}
type DiscardWriter struct{}                 // 自定义类型
func (DiscardWriter) Write([]byte) error {  // 只为丢弃数据，不需要具体数据
    return nil
}
~~~

### Method declarations

[Method_declarations](http://docscn.studygolang.com/ref/spec#Method_declarations)。方法（Method）是指拥有接收器（receiver）的函数。方法声明是将一个标识符（方法名）与一个方法绑定，同时将该方法与接收器的基本类型进行关联。接收器，是通过在方法名前添加一个额外的参数部分来指定。接收器必须是 T 类型的单参数，且 T 必须是自定义类型（defined type），或一个自定义类型指针。自定义类型 T 称为接收器的基本类型（base type）。接收器的基本类型要满足：不是指针或接口类型，且必须与方法定义在同一个包内。该方法被绑定到接收器的基本类型中，方法名仅在类型 T 或 \*T 的选择器（selector）表达式中可见。

在方法签名中，一个非空（non-blank）的接收器标识符必须是唯一的。如果接收器的值在方法体中不存在引用，则可以在方法签名中省略接收器的标识符。若接收器的基本类型是结构体（struct type），则非空的方法名和结构体字段名必须是不同的（方法名在某些情况下可为空字符）。

```go
// point（base type）已经在包中声明
// 基础类型 point 绑定了两个方法
func (p *Point) Length() float64 {   // method signature
    return math.Sqrt(p.x * p.x + p.y * p.y)
}

func (p *Point) Scale(factor float64) {
    p.x *= factor
    p.y *= factor
}
```

### FQA

#### Q：为什么要设计接收器函数 `func (receiverType) methodName` 与普通函数 `func methodName`？

* **接收器函数**（也称为方法）可以访问结构体或接口的字段和方法，而没有接收器的函数不能。
* 接收器函数可以修改结构体或接口的状态，而没有接收器的函数不能。
* 接收器函数可以作为结构体或接口的类型方法被调用，而没有接收器的函数不能。

简单的说：代码更易读、更易维护、提高代码的安全性、拥有面向对象的特性

* 提高代码的可读性和可维护性：接收器函数可以将与特定类型相关的方法组织在一起，使代码更易于阅读和维护。
* 提高代码的复用性：接收器函数可以被其他类型复用，这有助于减少代码重复。
* 提高代码的安全性：接收器函数可以对结构体或接口的状态进行封装，这有助于提高代码的安全性。

#### Q：自定义类型 T 的值接收器方法与指针接收器方法有什么区别？

例如：

~~~go
type T struct { /* T fields */ }
func (t T) MethodName { /* method implementation */ }  // 值接收器方法
func (t *T) MethodName { /* method implementation */ } // 指针接收器方法
~~~

- 接收器为 `T` 的方法是<b>值接收器方法</b>， 可以访问结构体本身的字段和方法，但不能修改结构体本身的状态。（只读）
- 接收器为 `*T` 的方法是<b>指针接收器方法</b>，可以访问和修改结构体本身的字段和方法。（读写）

##### 接收器方法集遵守的规则

- 实例 `o` 的类型是 `T`，则 o 的方法集包含接收器是 `T` 的所有方法；
- 实例 `o` 的类型是 `*T`，则 o 的方法集包含接收器是 `T` 和 `*T` 的所有方法；

注意：在某些情况下，会发现实例 o 是值（而不是指针）时，依旧是可以调用接收器为 `*T` 的方法，这似乎不符合上述规则。 
原因：值类型的变量调用接收器的指针类型的方法时，golang 会进行对该变量的取地址操作，从而产生出一个指针，之后再用这个指针调用方法。

前提是这个变量要能取地址。如果不能取地址，比如传入 interface 时的值是不可取地址的。示例如下：

~~~go
package main
import "fmt"

type I interface {
    Method()
}

type A struct {}
func (a A) Method() {  // 值接收器方法
    fmt.Println("A.Method")
}

type B struct {}
func (b *B) Method() {  // 指针接收器方法
    fmt.Println("B.Method")
}

func main() {
    var o1 I = A{}
    o1.Method()

    var o2 I = B{}   // err: B does not implement I (Method method has pointer receiver)
    var o2 I = &B{}  // suc:
    o2.Method()
}
~~~

#### Q：为什么 Go 语言把类型放在后面？

官方解释 [Go's Declaration Syntax](https://go.dev/blog/declaration-syntax)。分为类型前置和类型后置两种。变量类型后置、函数返回值后置，带来的代码可读性提高，类型推导更简单。

在类型后置上来讲，Go 官方核心思想是：这种声明方式（从左到右的风格）的一个优点是，当类型变得更加复杂时，它的效果非常好（One merit of this left-to-right style is how well it works as the types become more complex）。

Go 的变量名总是在前，在人的代码阅读上可以保持从左到右阅读，不需要像 C 语言一样在一大堆声明中用技巧找变量名对应的类型。C 语言的顺时针读法 [The Clockwise/Spiral Rule](http://c-faq.com/decl/spiral.anderson.html)。C# 设计组成员对类型前置、后置的设计教训，[Sharp Regrets: Top 10 Worst C# Features](https://www.informit.com/articles/article.aspx?p=2425867)。在设计时，C# 本来计划把类型注释放在右边。但考虑到类 C 语言，因此遵循了其他语言的惯例。
