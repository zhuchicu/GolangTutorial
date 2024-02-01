# 03-OOP
- 示例详情
    - 如何实现 OOP
    - Class
    - Inherit
    - Polymorphism
- 扩展内容
    - Struct types
    - Embedded field
        - 嵌入字段的提升
        - 嵌入字段的可见性
        - 嵌入字段的名字屏蔽
        - 嵌入字段对方法集的影响
    - Interface
        - 空接口
        - 嵌入接口
        - 类型断言
        - 接口的应用示例1：非空接口调用所有实现
        - 接口的应用示例2：空接口的类型断言
    - FQA
        - Q：如何“实例化”接口类型？
        - Q：如何知道结构体实现了某个接口？
        - Q：结构体中的匿名字段有什么用？
        - Q：使用 `new` 和 `make` 实例化结构体类型有什么区别？
    - Reference

## 示例详情

### 如何实现 OOP

面向对象（OOP）的三大特征是：封装、继承和多态。

*   封装：Go 使用 Struct 对属性进行封装，结构体就像是类的一种简化形式；
*   继承：用 Struct 中的内嵌匿名类型（embedded field）的方法来实现继承；
*   多态：使用接口（interface）实现多态。

### Class

基本属性：

*   public：首字母大写，能被其它包（package）访问或调用;
*   private：首字母小写，只能在包内使用；
*   访问 private 字段：getter 方法只使用成员名；setter 使用 Set 前缀；
*   method Sets：接收器函数（recevier method）作为类方法 `func (recevier T) MethodName(参数列表) 返回值列表 {`。不管方法的 <b>recevier</b> 类型 T 是值或指针，都可以调用，不必严格符合 recevier 的类型。如果定义了 recevier 是<b>值类型</b>的方法，也会隐含定义 recevier 是<b>指针类型</b>的方法。在值类型为大型 struct 时，使用指针 recevier 会更加高效。

setter 和 getter 的命名规则：Go 并不默认支持 setter 和 getter。在 getter 方法名中加入 Get 既不符合习惯，也没有必要。例如，一个未导出字段 `owner`，其 getter 方法称为 `Owner()`，getter 导出时使用大写名称提供了区分字段和方法的钩子。如果需要，setter 函数可能会被称为 `SetOwner()`。

### Inherit

因为 Golang 不支持传统意义上的继承，因此需要一种手段（即嵌入类型），把父类型的字段和方法“注入”到子类型中。

### Polymorphism

多态，不同对象中同种行为的不同实现方式。在 Go 语言中可以使用接口实现这一特征，interface 是一组<b>抽象方法</b>的组合，通过 interface 来定义对象的一组行为。如果某个对象实现了某个接口的所有方法，则此对象就实现了此接口。

Go 通过 interface 实现了 duck-typing：“If it looks like a duck, swims like a duck, and quacks like a duck, then it probably is a duck.”

## 扩展内容

### Struct types

结构体由一系列字段（有类型的标识符）组成，域名可以显式指定（IdentifierList）或隐式指定（EmbeddedField）的。在结构体中，非空域名必须是唯一的。在声明域时，可以在尾部添加一个 string 类型的 tag，它主要用于<b>在反射接口（reflection interface）的访问中，作为结构体的类型识别</b>。空字符串 tag 等同于缺省。

- 使用 `new` 或 `make` 函数来实例化 struct 类型；
- 使用选择器（selector）表达式来访问 struct 字段，例如 `S.f`（Go 的设计者认为 `.` 比 `->` 更简洁、更易读）；
- 嵌入类型：将一个 struct 类型作为字段嵌入到另一个 struct 类型中；
- 定义接收器方法：`func (s *S) MethodName(params) return {` ；
- struct 类型可以实现接口，例如实现 `fmt.Stinger` 接口：`var i Stringer = S{}`。

~~~go
struct {}       // An empty struct.
struct {
    T1        // embedded fields, field name is T1
    *T2       // embedded fields, field name is T2，T2 本身也不能是指针类型
    P.T3      // embedded fields, field name is T3，P 是限定类型，即包名
    *P.T4     // embedded fields, field name is T4
    x, y int  // field names are x and y
}
struct {      // 非法的声明方式：字段名必须是唯一的
    T         // conflicts with embedded field *T and *P.T
    *T        // conflicts with embedded field T and *P.T
    *P.T      // conflicts with embedded field T and *T
}
~~~

在尾部添加一个 string 类型的 tag 示例：

~~~go
/* optional string literal tag */
struct {
    x, y float64 ""  // an empty tag string is like an absent tag
    name string  "any string is permitted as a tag"
    _    [4]byte "ceci n'est pas un champ de structure"
}
struct {   // TimeStamp protocol buffer
    microsec  uint64 `protobuf:"1"`  // define the protocol buffer field numbers
    serverIP6 uint64 `protobuf:"2"`  // 遵循反射包概述的约定
}
~~~


##### struct 的初始化

- 结构体字面量（struct literal）简单直接，不需要使用指针。
- `new` 可以对 struct 进行一些初始化操作，例如将字段设置为零值（zero value，即默认值）。

~~~go
type S struct {
    x int
    y string
}

s1 := S{x: 10, y: "Hello, world!"}  // struct literal 初始化
s1 := S{10, "Hello, world!"}        // 字段名可省略

s2 := new(S)
fmt.Println("x=", s2.x)  // x=0，int 类型的零值为 0
s2.y = "Hello, world!"
~~~

### Embedded field

如果在结构体 `S` 中声明了类型但没有名称的字段 `F`，那 `F` 就叫做<b>嵌入字段（embedded field）</b>。在引用时，字段的类型名会被当成该字段的名字。示例如下：

~~~go
type F struct { /* fields */ }  // 自定义类型（defined type）
func (f *F) Log() {}

type S struct {
    *F    // 嵌入字段，也可以是 F
}
~~~

嵌入字段 `F` 是一个自定义类型（defined type），`F` 是类型名 `T`，或一个指向非接口类型的指针类型名 `*T`（且 `T` 本身不能是一个指针类型）。可通过 `<Type>` 名称引用其方法和属性。示例如下：

~~~go
func main() {
    s := S{&F{/* ... */}}  // 按 F 的字段出现的顺序给出相应的初始化值
    s.F.Log()              // “*”只是类型修饰符，F 才是名称
    s.F.fields             // 也可以引用 F 的数据字段
}
~~~

注意：指针 `*` 只是类型修饰符，并不是类型名的一部分，如 `*<Type>` 和 `<Type>`，只能通过 `<Type>` 这个名字进行引用。


#### 嵌入字段的提升

上述示例中，`S` 对于 `F` 的字段和方法，还可以直接调用，而不需要给出 `F` 的类型名，如 `s.Log()` 和 `s.fields`，这种行为叫做<b>嵌入字段的提升（promoted embedded field）</b>，`F` 也叫<b>提升字段（Promoted field）</b>。具体示例如下：

~~~go
s.Log()    // s.F.Log()
s.fields   // s.F.fields
~~~


#### 嵌入字段的可见性

原因：<b>首字母小写的字段和方法是包私有的，而首字母大写的可以在任意地方被访问</b>。同一个包里的东西是彼此互相公开（public）的。示例如下：

~~~go
package a
type F struct{
    x int
}
func(f *F) log() {}
~~~

~~~go
package main
import (
    "a"
    "fmt"
}

type S struct {
    *a.F
}
func main() {
    s := S{&a.F{1}}
    // 以下都是错误的
    fmt.Println(s.x)    // s.x undefined (type S has no field or method x)
    fmt.Println(s.F.x)  // s.F.x undefined (cannot refer to unexported field or method a.(*F).x)
    s.log()   // s.log undefined (type S has no field or method log)
    s.F.log() // s.F.log undefined (cannot refer to unexported field or method a.(*F).log)
}
~~~

#### 嵌入字段的名字屏蔽

名字屏蔽：与当前类型的字段或者方法同名时，会屏蔽嵌入类型的。若要访问嵌入类型中的字段或方法，不使用字段提升即可。因此<b>不要让多个嵌入类型包含同名字段或方法</b>。示例如下：

~~~go
package main
import "fmt"

type F struct {
    x int
}
func (f *F) Log() {
    fmt.Println("F Method: Log")
}

type S struct {
    F
    x int
}
func (s *S) Log() {
    fmt.Println("S Method: Log")
}

func main() {
    s := S{F{1}, 2}
    s.Log()
    s.F.Log()                   // 不使用字段提升即可
    fmt.Println("S.x=", s.x)
    fmt.Println("F.x=", s.F.x)  // 不使用字段提升即可
}
~~~

当 `S` 中有多个嵌入字段，且都存在的同名的字段或时，`S` 使用字段提升方式进行访问时，会发生“调用时发生了二义性错误”。

#### 嵌入字段对方法集的影响

类型 `T` 的方法集有“值接收器”和“指针接收器”的区别，它们遵循以下原则：

- 实例 `o` 的类型是 `T`，则 `o` 的方法集包含接收器是 `T` 的所有方法；
- 实例 `o` 的类型是 `*T`，则 `o` 的方法集包含接收器是 `T` 和 `*T` 的所有方法；

嵌入类型也分为值和指针，它与普通变量的方法集规律一致。当结构体 `S` 与自定义类型 `T` 结合时，`S` 的方法集中包含的提升方法（即嵌入字段的方法）遵循以下原则：

- `S` 含嵌入字段 `T`（值类型嵌入）：`S` 和 `*S` 的方法集都将包含拥有“字段提升”的接收器方法 `T`；`*S` 的方法集将包含拥有“字段提升”的接收器方法 `*T`。
- `S` 含嵌入字段 `*T`（指针类型嵌入）：`S` 和 `*S` 的方法集都将包含拥有“字段提升”的接收器方法 `T` 或 `*T`。

示例如下：

~~~go
type T struct { /* T fields */ }
func (t *T) PointerMethod() { /* implementation */ }
func (t T) ValueMethod() { /* implementation */ }

type IPtrMethod interface {
    PointerMethod()
}
type IValMethod interface {
    ValueMethod()
}
~~~

~~~go
type S struct { T }   // 值类型嵌入

func main() {
    var ptrIPtr IPtrMethod = &S{}  // 类型为 *S 的实例 ptrIPtr 拥有 *T 的方法
    var ptrIVal IValMethod = &S{}  // 类型为 *S 的实例 ptrIVal 拥有 T 的方法
    fmt.Println(ptrIPtr, ptrIVal)
    
    // S does not implement IPtrMethod
    var valIPtr IPtrMethod = S{}  // err：类型为 S 的实例 valIPtr 未拥有 *T 的方法
    var valIVal IValMethod = S{}  // 类型为 S 的实例 valIVal 拥有 T 的方法
    fmt.Println(valIPtr, valIVal)
}
~~~

~~~go
type S struct { *T }  // 指针类型嵌入

func main() {
    var ptrIPtr IPtrMethod = &S{}  // 类型为 *S 的实例 ptrIPtr 拥有 *T 的方法
    var ptrIVal IValMethod = &S{}  // 类型为 *S 的实例 ptrIVal 拥有 T 的方法
    fmt.Println(ptrIPtr, ptrIVal)
    
    var valIPtr IPtrMethod = S{}  // 类型为 S 的实例 valIPtr 拥有 *T 的方法
    var valIVal IValMethod = S{}  // 类型为 S 的实例 valIVal 拥有 T 的方法
    fmt.Println(valIPtr, valIVal)
}
~~~

总结：结构体 `S` 的对象 `o` 是值或指针类型时，理论上都能调用其嵌入字段的方法集（即接收器类型为 `T` 或 `*T`），除了一种特殊情况外，就是<b>当 `o` 为值类型时，不能调用其指针类型嵌入的方法集，即接收器类型为 `*T` 的方法</b>。

### Interface

接口就是一组方法签名的集合。接口类型指定了一个方法集，并把这个方法集称为自己接口。<b>接口类型变量</b>可以存储任何实现了该接口（也即实现了接口声明的所有方法）的类型变量，未初始化的接口变量值为 nil。接口类型中的方法集，可以通过<b>直接显式声明方法</b>，也可以<b>通过接口类型名称嵌入其他接口的方法集</b>形式。显式声明的方法必须是唯一且非空的。

示例如下：

~~~go
interface {   // A simple File interface.
    Read([]byte) (int, error)     // method signature
    Write([]byte) (int, error)    // 显示声明
    Close() error
}

interface {
    String() string
    String() string  // illegal: String not unique
    _(x int)         // illegal: method must have non-blank name
}
~~~

#### 空接口

空接口特性如下：

*   每一个接口都包含两个属性：值和类型 `(type, value)`；
*   一个接口可以被多个类型实现（一个类型也可以实现多个不同的接口）；
*   所有类型都至少实现了空接口 `interface{}`；
*   `interface{}` 可以存储任意类型数值（类似于 C 语言的 `void*` 类型），但不代表任意类型也能承接 `interface{}` 类型的值；
*   `interface{}` 作为参数时，可以接受任何值，参数类型是 `interface{}`；
*   `interface{}` 作参来接收任意类型时，需要[“断言”判断类型](https://golang.google.cn/ref/spec#Type_assertions)，失败会 panic；
*   `interface{}` 承载数组和切片后，该对象无法再进行切片；

~~~go
type I interface {
    M()
}

// ===================================================
type S struct {
    x int
}
func (s *S) M() {}
var i I = &S{x: 10} // 接口值 i 包含类型 *S 和值 10

// ===================================================
func f(i I) {  // 接口类型值作为函数参数
    i.M()      // 调用接口方法
}
func g() I {  // 作为函数的返回值
    return &S{x: 10}
}

// ===================================================
// Type_assertions
var x interface{} = 7          // x has dynamic type int and value 7
i := x.(int)                   // i has type int and value 7
~~~

#### 嵌入接口


如果接口类型 `T` 中有一个嵌入接口 `I`，那么 `T` 的方法集则是 `T` 的显式声明方法集与 `T` 的嵌入接口所含方法集的<b>联合集（union）</b>。方法的联合集只会包含所有（导出和非导出）方法集中一个方法，且所有方法集中同名方法的签名必须是一致的。注意：接口类型的嵌入接口不能是自身，也不能存在循环嵌入的情况。

如下示例中，`ReadWriter` 的方法联合集中只会存在一个 `Close()` 方法：

~~~go
type Reader interface {
    Read(p []byte) (n int, err error)
    Close() error
}
type Writer interface {
    Write(p []byte) (n int, err error)
    Close() error
}
type ReadWriter interface {  // ReadWriter's methods are Read, Write, and Close.
    Reader  // includes methods of Reader in ReadWriter's method set
    Writer  // includes methods of Writer in ReadWriter's method set
}
~~~

#### 类型断言

类型断言的表达式：`x.(T)`，表示 x 不是 nil，并且存储在 x 中的值是类型 T。更准确地说，如果 T 不是接口类型，`x.(T)` 断言 x 的动态类型与类型 T 相同。在这种情况下，T 必须实现 x 的(接口)类型；否则类型断言无效，因为 x 不可能存储类型T  的值。如果 T 是接口类型，则 `x.(T)` 断言 x 的动态类型实现接口 T。

类型断言表达式的返回值为 true 或 false。虽然 x 的动态类型仅在运行时（Run-time）可知，但在正确的程序中，`x.(T)` 的类型就是 T，否则运行时发生 panic。

```go
var x interface{} = 7          // x has dynamic type int and value 7
i := x.(int)                   // i has type int and value 7

type I interface { m() }
func f(y I) {
    s := y.(string)        // illegal: string does not implement I (missing method m)
    r := y.(io.Reader)     // r has type io.Reader and the dynamic type of y must implement both I and io.Reader
    …
}
```

在赋值或初始化语句中使用类型断言表达式，它将产生一个额外的无类型化布尔值（untyped）。当断言成立时，ok 的值为 true，否则 ok 为 false 且 v 的值为类型 T 的 zero value（布尔值为 false，数字类型为 0，字符串为 `""`，指针、函数、接口、切片、信道和映射为 nil）。同时这种特殊格式将不会在运行时产生 panic。示例如下：

```go
v, ok = x.(T)                 // 赋值
v, ok := x.(T)                // 短声明格式
var v, ok = x.(T)             // 声明
var v, ok interface{} = x.(T) // dynamic types of v and ok are T and bool
```

#### 接口的应用示例1：非空接口调用所有实现

通过接口类型调用多个实现了 Stringer 接口的类型的 String() 方法：

~~~go
type Stringer interface {
    String() (string)
}

type T struct {}               // 类型 S 与 T 相同的定义
func (t T) String() string {   // 也实现了接口的 String()
    return "T"
}

func PrintAll(s []Stringer) {
    for _, i := range s {
		fmt.Println("接口调用", i.String())
	}
}

func main() {
    stringers := []Stringer{T{}, S{}}  // 同时声明多个接口类型变量
    PrintAll(stringers)
}
~~~

上述示例中，也可以直接使用类型 T 或 S 的变量调用接口方法。

#### 接口的应用示例2：空接口的类型断言

~~~go
type I interface{
    print()
}
type S struct {
    x int
}
func (t S) print() {
    fmt.Println("x=", t.x)
}

// 参数为空接口类型，能够接受任意类型，按需要进行“类型断言”
func f(i interface{}) {
  j, _ := i.(I)
  j.print()
}

func main() {
  var i I = S{10}
  f(i) // 输出：x=10
}
~~~

### FQA

#### Q：如何“实例化”接口类型？

Go 中的接口类型不能直接实例化，因为接口类型本身并不包含任何数据或方法。接口的应用：

~~~go
package main
import "fmt"

type Stringer interface {
    String() string
}

type Person struct {
    name string
}
func (p Person) String() string {  // 实现了接口的方法 String
    return "Name: " + p.name
}

func main() {
    p := Person{name: "John Doe"}
    var s Stringer = p       // 使用 Person 结构体来实例化 Stringer 接口类型
    fmt.Println(s.String())  // 调用 Stringer 接口类型的方法
    // 输出：Name: John Doe
}
~~~

上述示例中的 `s.String()`，和 `p.String()` 的主要区别在于：

- `s.String()` 是通过接口类型 `Stringer` 调用 `String()` 方法，而 `p.String()` 是通过结构体类型 `Person` 调用 `String()` 方法。
- `s.String()` 可以调用<b>任何实现了 `Stringer` 接口的类型的 `String()` 方法</b>，而 `p.String()` 只可以调用 `Person` 结构体的 `String()` 方法。


#### Q：如何知道结构体实现了某个接口？

例如，结构体 `S` 实现了某接口方法，示例如下：

~~~go
package main
import (
    "fmt"     // fmt.Stringer 接口
    "custom"  // 自定义包 custom.Stringer 接口
)

type S struct { /* fields */ }
func (s S) String() string { return "S" }

func main() {
    var s fmt.Stringer = S{}     // 使用显式类型断言来指定实现的是什么接口
    var s custom.Stringer = S{}
}
~~~

- 接口的来源（即声明）只有两种：当前包内和导入的外部包。
- 实现接口时，都会指明是哪个接口。根据来源查看接口声明，确认接口内所有的方法是否都实现了。
- 若导入的外部包中存在同名的接口，则会使用显式类型断言来指定实现，例如上述示例。


#### Q：结构体中的匿名字段有什么用？

- 代码复用：嵌入字段实现“继承”效果（利用了字段提升）；
- 嵌入接口：嵌入了接口 `I` 的结构体 `S`，可以作为任何需要 `I` 接口的函数或方法的参数；
- 内存对齐： 从而提高内存访问的效率。

结构体中嵌入接口作为匿名字段，示例如下：

~~~go
type Writer interface {
    Write(p []byte) (n int, err error)
}

type File struct { Writer } // Writer 接口作为匿名字段
func (f *File) Write(p []byte) (n int, err error) {
    // implements...
    return len(p), nil
}

func main() {
    file := &File{}
    n, err := file.Write([]byte("Hello, world!"))
    if err != nil { /* 处理错误 */ }
    fmt.Println(n) // 打印写入的字节数，输出：13
    
    // =========================================================
    // 将 File 结构体作为任何需要 Writer 接口的函数或方法的参数
    file := &File{
        Writer: os.Stdout, // 将标准输出作为 Writer。
    }
    fmt.Fprintf(file, "Hello, world!")  // 将数据写入文件，将数据写入标准输出。
}

~~~

##### 匿名字段对齐结构体中的字段

Go 中结构体中的字段必须以“字节对齐”的方式存储。这意味着每个字段的起始地址必须是 2 的幂的倍数。例如，如果一个字段的大小是 1 字节，那么它的起始地址必须是 1、2、4、8 等。如果一个结构体中的字段没有对齐，那么编译器就会在字段之间插入<b>填充字节（padding bytes）</b>来对齐字段。这些填充字节不会被使用，它们只是为了确保字段的起始地址是正确的。

例如，下面示例中的 `_` 变量是一个匿名字段，它的大小为 4 字节。它被用来对齐结构体中的其他字段，以确保它们在内存中以正确的顺序和位置存储：

~~~go
struct {        // A struct with 6 fields.
    x, y int    // 占 4 字节，x 起始地址 1，y 是 4
    u float32   // 占 4 字节，起始地址 8
    _ float32   // padding bytes，4 字节
    A *[]int    // 指针类型占 8 字节，起始地址 16=2^4
    F func()    // 函数类型占 8 字节
}
// 若没有匿名字段，A 的起始地址将为 12，导致字节不对齐
~~~

注意：`_` 变量不能被赋值，也不能被访问。它只是一个占位符，用来对齐结构体中的其他字段。

在实际编程中，开发者通常不需要考虑字节对齐。编译器会自动处理字节对齐，以确保结构体中的字段以正确的顺序和位置存储。在某些情况下需要考虑字节对齐。例如，当开发者需要与其他语言编写的代码进行交互时，或者当程序员需要对结构体中的数据进行底层操作时。

#### Q：使用 `new` 和 `make` 实例化结构体类型有什么区别？

对于引用类型的变量，不光要声明它，还要为它分配内容空间。`func new(Type) *Type` 返回的永远是类型的指针，指向分配类型的内存地址。

~~~go
var i *int = new(int)
*i = 10
fmt.Println(*i)
~~~

`func make(t Type, size ...IntegerType) Type` 只能用于内置结构体类型的实例化，即 `chan`、`map`、`slice`。它返回的类型就是这三个类型本身，而不是它们的指针类型，因为这三种类型就是引用类型，所以就没有必要返回他们的指针了。

~~~go
sli := make([]int, 3)
ch := make(chan int, 1) // 创建有 1 个缓冲的 channel
m := make(map[string]float32, 100)  // 初始容量 capacity=100
~~~

二者都是内存的分配（堆上），`new` 用于类型的内存分配，并且内存置为零；`make` 只用于 `slice`、`map` 以及 `channel` 的初始化（非零值）。

### Reference

1. [golang拾遗：嵌入类型](https://www.cnblogs.com/apocelipes/p/14090671.html)

