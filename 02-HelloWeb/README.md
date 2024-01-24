# 02-HelloWeb

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
