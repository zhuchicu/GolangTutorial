# 05-GoWebExample

## WebServer Package


## 扩展内容


### Standard library

#### Package API 文档解析

以官方标准包 [net/http](https://pkg.go.dev/net/http) 为例，页面主要内容分为三类：API 文档、源文件和子目录。文档内容的层级结构以 [index](https://pkg.go.dev/net/http#pkg-index) 作为参考说明。官方文档在布局和排版时会以<b>相关性、可读性和一致性</b>为原则，即以逻辑而不是语言规则为主：

- 例如 Constants 中会将常用的 HTTP 方法名称作为一组常量，使用块语句（block）声明；
- 例如 Variables 部分会将 DefaultServeMux 单独成行；
- Functions 中以包中的“函数”为主，只要被纳入其中必然是“函数”而不是方法；
- Types 部分将会以自定义类型为核心，将相关的类型声明、“函数”和“方法”汇总在一起，例如 `type Request` 中的 [`func NewRequest`](https://pkg.go.dev/net/http#NewRequest) 使用“函数”而不是其接收器方法，但是由于该函数属于类型的实例化，所以将其合并在一起。同样的情况还有 `type Handler interface` 其包含的函数都能返回该接口类型，例如 `func FileServer(root FileSystem) Handler`，所以也汇总在一起。 
