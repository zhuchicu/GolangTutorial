# OOP

## 示例详情

### 如何实现 OOP

面向对象（OOP）的三大特征是：封装、继承和多态。

*   封装：Go 使用 Struct 对属性进行封装，结构体就像是类的一种简化形式；
*   继承：用 Struct 中的内嵌匿名类型（embedded field）的方法来实现继承；
*   多态：使用接口（interface）实现多态。

### Class：结构体

基本属性：

*   public：首字母大写，能被其它包（package）访问或调用;
*   private：首字母小写，只能在包内使用；
*   访问 private 字段：getter 方法只使用成员名；setter 使用 Set 前缀；
*   method Sets：接收器函数（recevier method）作为类方法 `func (recevier T) MethodName(参数列表) 返回值列表 {`。不管方法的 <b>recevier</b> 类型 T 是值或指针，都可以调用，不必严格符合 recevier 的类型。如果定义了 recevier 是<b>值类型</b>的方法，也会隐含定义 recevier 是<b>指针类型</b>的方法。在值类型为大型 struct 时，使用指针 recevier 会更加高效。


### Inherit：嵌入字段


因为 Golang 不支持传统意义上的继承，因此需要一种手段（即嵌入类型），把父类型的字段和方法“注入”到子类型中。如果在结构体 `S` 中声明了类型但没有名称的字段 `F`，那 `F` 就叫做<b>嵌入字段（embedded field）</b>。在引用时，字段的类型名会被当成该字段的名字。示例如下：

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

### Polymorphism：接口


## 扩展内容

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



