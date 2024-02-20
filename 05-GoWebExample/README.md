# 05-GoWebExample

官方推荐的 Web 教程 [GoWebExamples](https://gowebexamples.com/) 中引入了多个非 Go 官方的标准库：
- 比官方标准库提供更好的性能，但[开源项目 Gorilla 已停止维护](https://tehub.com/a/aVYcm4Jomq)，更多见 [gorilla-toolkit](https://github.com/gorilla#gorilla-toolkit)。
- 前缀为 `golang.org/x/...` 的包是 Go 项目的一部分，但位于 Go 主干之外。它们是在比 Go 内核更宽松的兼容性要求下开发的。一般来说，它们将支持前两个版本和提示（来源：[Go wiki](https://golang.google.cn/wiki/X-Repositories)）。
- Go 官方标准库 [database/sql](https://pkg.go.dev/database/sql) 提供了标准的、轻量的、面向行的接口，但不提供具体数据库驱动，只提供驱动接口和管理（这样是为了确保向前兼容，无法预知未来有哪些的数据库，且没有精力维护大量的驱动）。要使用数据库还需要引入想使用的特定数据库驱动，例如 [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)。


## WebServer Package

### Package net/http

官方标准包 [net/http](https://pkg.go.dev/net/http) 提供了所有关于 HTTP 协议相关的方法。

##### 创建 webserver 监听连接

注册 Request 处理器 [HandleFunc](https://pkg.go.dev/net/http#HandleFunc)，用于接收解析请求并编写响应内容。HTTP 服务器必须侦听端口才能将连接传递到请求处理程序，所以使用 [ListenAndServe](https://pkg.go.dev/net/http#ListenAndServe)（[Source](https://cs.opensource.google/go/go/+/refs/tags/go1.21.6:src/net/http/server.go;l=3237)，其中 [DefaultServeMux](https://cs.opensource.google/go/go/+/refs/tags/go1.21.6:src/net/http/server.go;l=2336) 就是 ServeMux）对指定的端口进行监听。
~~~go
// func HandleFunc(pattern string, handler func(ResponseWriter, *Request))  // 函数签名，注意区别方法签名
// handler 的 http.ResponseWriter：使用 text/htm 格式编写请求的响应内容
// handler 的 http.Request：包含 HTTP 请求的所有信息，比如 URL 或请求头的域
http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
})

// func ListenAndServe(addr string, handler Handler) error  // 函数签名
// handler 通常是 nil，在这种情况下使用 DefaultServeMux
http.ListenAndServe(":80", nil)
~~~

##### 访问静态资源

在上述示例的基础上，使用 [http.FileServer]() 提供 JavaScript、CSS 和图像等静态资源的访问。为了正确地提供文件，我们需要使用 [http.StripPrefix](https://pkg.go.dev/net/http#StripPrefix) 去掉 url 路径的一部分。通常这是我们的文件所在目录的名称。源码见 [HTTP Server](https://gowebexamples.com/http-server/)，关键代码如下：

~~~go
// func FileServer(root FileSystem) Handler  // 函数签名
// 参数 root 的定义为 type FileSystem interface { // 接口定义
fs := http.FileServer(http.Dir("static/"))

// func Handle(pattern string, handler Handler)        // 函数签名
// func StripPrefix(prefix string, h Handler) Handler  // 函数签名
http.Handle("/static/", http.StripPrefix("/static/", fs))
~~~

关于 <b>`http.Dir()` 函数</b>的问题，首先该类型定义为 [`type Dir string`](https://pkg.go.dev/net/http#Dir)，底层类型是 [`string`](https://pkg.go.dev/builtin#string)。Go 中有强制类型转换的语法 `T(expression)`，T 是目标类型，因此 `http.Dir()` 表示将某类型的变量强制转换为 `Dir` 类型（注：`string` 类型可强转 `Dir`，但 `Dir` 不可强转 `string`）。示例如下：

~~~go
package main
import (
    "fmt"
    "net/http"
    "reflect"
)

func main() {
    var s string = string("hh")
    fmt.Println(reflect.TypeOf(s).Name())  // 输出 string
    f := http.Dir("/")
    fmt.Println(reflect.TypeOf(f).Name())  // 输出 Dir
}
~~~

##### 访问静态资源的完整示例

上述静态资源访问的完整示例见 [Assets and Files](https://gowebexamples.com/static-files/)：

~~~bash
$ tree static-files
static-files.go
assets/
 └── css
    └── styles.css
~~~

~~~go
func main() {
    fs := http.FileServer(http.Dir("assets/"))                  // 相对位置，资源存储的本地路径
    http.Handle("/static/", http.StripPrefix("/static/", fs))   // 路径前缀，请求访问的 css/styles.css
    http.ListenAndServe(":8080", nil)
}
~~~

~~~bash
$ curl -s http://localhost:8080/static/css/styles.css
# 输出结果如下，即 styles.css 的内容
# body {
#    background-color: black;
# }
~~~

### Package gorilla/mux

官方标准包 net/http 提供了很多 HTTP 协议相关的方法，但还是存在一些局限性，比如将请求 url 分割成单个参数，而包 [gorilla/mux](https://pkg.go.dev/github.com/gorilla/mux) 能够提供这些能力（更多详细见官方文档 [Routing (using gorilla/mux)](https://gowebexamples.com/routes-using-gorilla-mux/)）。例如，<code>/books/<b>go-programming-blueprint</b>/page/<b>10</b></code> 请求链接中有两个动态参数，使用占位符表示 <code>/books/<b>{title}</b>/page/<b>{page}</b></code>。

- 使用命令 `$ go get -u github.com/gorilla/mux` 获取包；
- `mux.NewRouter()` 创建路由器，`mux.Vars()` 将当前请求的路由变量存储在 map 中（key 为占位符字符，如上文的 `["title"]`）；
- [`func (*Router) HandleFunc`](https://pkg.go.dev/github.com/gorilla/mux#Router.HandleFunc) 一些常用特性：`Methods("POST")` 将请求处理程序限制为特定的 HTTP 方法；`Host()` 限制为特定主机名或子域；`Schemes()` 限制为 http/https；
- [`Route.PathPrefix()`](https://pkg.go.dev/github.com/gorilla/mux#Router.PathPrefix) 使用 URL 路径前缀的匹配器注册新<b>路由 `Route`</b>（注意与 `Router` 的区别），[`Route.Subrouter().HandleFunc()`](https://pkg.go.dev/github.com/gorilla/mux#Route.Subrouter) 将请求处理程序限制为特定路径前缀。

### Package database/mysql

1. 使用 `$ go get -u github.com/go-sql-driver/mysql` 命令下载 MySql 驱动包，API 文档见 `pkg.go.dev` 中的 [github.com/go-sql-driver/mysql](https://pkg.go.dev/github.com/go-sql-driver/mysql)。
1. 可以使用 Docker 启动一个 Mysql 实例，具体见[文档](https://hub.docker.com/_/mysql)；
1. 导入 [`database/sql`](https://pkg.go.dev/database/sql) 包（以及使用空占位符用法 `import _ "go-sql-driver/mysql"`），使用 `sql.Open(driverName, dataSourceName)` 连接 Mysql（`dataSourceName` 的规则见 [go-sql-driver/DSN](https://github.com/go-sql-driver/mysql?tab=readme-ov-file#dsn-data-source-name)）；
1. 使用 `db.Exec()` 执行 SQL 语句来创建数据库的表，同时 SQL 插入一条有效数据（并获取插入的自增 ID）；
1. 查询数据库的方式有两种，`db.Query(query string, args ...any)` 一次返回多行数据用于遍历，而 `db.QueryRow` 只返回至多一行指定的数据；
1. 使用 `db.Exec()` 删除数据。

~~~go
// 数据库连接
db, err := sql.Open(driverName, dataSourceName string)
err := db.Ping()  // Ping 验证与数据库的连接是否仍然有效，必要时建立连接

// Exec 执行 SQL 语句进行创建、插入、删除，不返回行数据，但会返回插入的 ID
result, err := db.Exec(query string, args ...any)
userID, err := result.LastInsertId()  // 将插入的自增 ID 作为用户的主键
_, err := db.Exec(`DELETE FROM users WHERE id = ?`, 1) // 删除，check err

// 查询并返回结果
err := db.QueryRow(query, 1).Scan(&id, &username, &password, &createdAt)
rows, err := db.Query(query)  // 返回多行
var users []user
for rows.Next() {
    var u user
    err := rows.Scan(&u.id, &u.username, &u.password, &u.createdAt) // check err
    users = append(users, u)
}
err := rows.Err() // check err
~~~

##### 为什么导入 MySQL 驱动程序包使用空占位符？

导入包使用空占位符 `_` 表示不会使用包内的变量或方法，如果包内存在 `init()` 函数，则会调用。在上述应用中，导入的 MySQL 驱动程序包拥有 `init` 函数。以下是 MySQL 驱动程序代码的截取片段，完整源码见 [driver.go](https://github.com/go-sql-driver/mysql/blob/master/driver.go)：

~~~go
import (
    "database/sql"
    "database/sql/driver"
    // ...
)

// MySQLDriver is exported to make the driver directly accessible.
// In general the driver is used via the database/sql package.
type MySQLDriver struct{}

// This variable can be replaced with -ldflags like below:
// go build "-ldflags=-X github.com/go-sql-driver/mysql.driverName=custom"
var driverName = "mysql"

func init() {
    if driverName != "" {
        sql.Register(driverName, &MySQLDriver{})
    }
}
~~~

上述截取代码中，导入驱动包时会调用官方标准包中的 `sql.Register()` 函数注册 MySQL 驱动程序（即定义的 `MySQLDriver`）。

### Package html/template

包 [`html/template`](https://pkg.go.dev/html/template) 为 HTML 模板提供了丰富的模板语言。它主要用于 Web 应用程序，在客户端浏览器中以结构化方式显示数据。这个包还包含了 [`text/template`](https://pkg.go.dev/text/template) 包，可以共享其模板 API 来安全地解析和执行 HTML 模板。

##### 控制结构

要访问模板中的数据，最上面的变量是通过 `{{.}}` 访问的。大括号内的“点”称为管道和数据的根元素。模板语言包含一组丰富的控制结构来呈现 HTML：

~~~go
// Control Structure
{{/* a comment */}}              // 注释
{{.}}                            // 渲染 root 元素
{{.Title}}                       // 渲染嵌套在元素内的“Title”域
{{if .Done}} {{else}} {{end}}    // if 语句
{{range .Todos}} {{.}} {{end}}   // 循环语句，遍历“Todos”并使用 root 元素进行渲染
{{block "content" .}} {{end}}    // 定义一个名为“content”的块（block）
~~~

示例 `layout.html` 中，模板引擎能够正确地判断 `Done` 字段所在的层级，并正确地输出数据：

~~~html
<h1>{{.PageTitle}}</h1>
<ul>
    {{range .Todos}}
        {{if .Done}}
            <!-- 在官方教程的基础上，添加了不可交互的勾选框样式 -->
            <li class="done"><input type="radio" checked>{{.Title}}</li>
        {{else}}
            <li><input type="radio">{{.Title}}</li>
        {{end}}
    {{end}}
</ul>
~~~

~~~go
data := TodoPageData{
    PageTitle: "My TODO list",
    Todos: []Todo{
        {Title: "Task 1", Done: false},   // Done 的层级
        {Title: "Task 2", Done: true},
        {Title: "Task 3", Done: true},
    },
}
~~~

##### 从文件中解析模板

模板可以从字符串或磁盘上的文件中解析。通常情况下，模板是从磁盘中复制的，示例：

~~~go
// func ParseFiles(filenames ...string) (*Template, error)
tmpl, err := template.ParseFiles("layout.html")  // layout.html 与 Go 程序位于同一目录中
// or
// func Must(t *Template, err error) *Template  // 区别是旨在用于变量初始化
tmpl := template.Must(template.ParseFiles("layout.html"))
~~~

##### 接受请求并进行渲染输出

使用 `template.Execute()` 接受用于写出模板的 `io.Writer` 和用于将数据传递到模板的 `interface{}`：

~~~go
// func (t *Template) Execute(wr io.Writer, data any) error
func(w http.ResponseWriter, r *http.Request) {
    tmpl.Execute(w, "data goes here")
}
~~~

##### POST 请求提交表单数据

`template.Execute()` 将解析后的模板（`forms.html`）应用于指定的数据对象，并将输出写入 `wr`。<b>如果在执行模板或写入输出时发生错误，执行将停止，但部分结果可能已写入输出写入器</b>。模板可以安全地并行执行，但如果并行执行共享一个写入器，输出可能会交错。具体源码见 [Forms](https://gowebexamples.com/forms/)。

使用 `func(w http.ResponseWriter, r *http.Request)` 对用户请求进行解析：

- 判断请求方法类型（POST）：`r.METHOD != http.MethodPost`；
- 获取（POST）请求的参数值：`r.FormValue("email")`；

源码 [Forms](https://gowebexamples.com/forms/) 中代码片段解析：

~~~go
tmpl.Execute(w, struct{ Success bool }{true})   // struct{ Success bool }{true} 
// struct{ Success bool } 表示结构体的定义，一个匿名的结构体变量
// {true} 使用结构体的字面值形式，对匿名的结构体变量进行初始化操作
// 结构体 Success 字段的值为 true，对应 forms.html 中的声明
~~~

~~~html
<!-- forms.html -->
{{if .Success}}  <!-- 执行 Success==true 的块作用域 -->
    <h1>Thanks for your message!</h1>
{{else}}
    <!-- ... -->
{{end}}
~~~

### Package log

创建一个<b>日志中间件</b>（logging middleware，使用 [`log`](https://pkg.go.dev/log) 包）：中间件只需将 [`http.HandlerFunc`](https://pkg.go.dev/net/http#HandlerFunc) 作为参数之一，对其进行封装，然后返回一个新的 `http.HandlerFunc` 供服务器调用。封装如下（源码见 [basic middleware](https://gowebexamples.com/basic-middleware/)）：

~~~go
func logging(f http.HandlerFunc) http.HandlerFunc {        // 传入参类型与返回类型一致
    return func(w http.ResponseWriter, r *http.Request) {  // 匿名函数
        log.Println(r.URL.Path)                            // 多添加的执行
        f(w, r)                                            // 原有的执行流程
    }
}
~~~


##### 构建高级的中间件

在这里，我们定义了一种新的中间件（Middleware）类型，它可以让多个中间件更容易地连锁在一起。这个想法的灵感来自 Mat Ryers 关于构建应用程序接口的演讲。你可以在这里找到更详细的解释，包括演讲内容。


## 扩展内容


### Standard library

#### Package API 文档解析

以官方标准包 [net/http](https://pkg.go.dev/net/http) 为例，页面主要内容分为三类：API 文档、源文件和子目录。文档内容的层级结构以 [index](https://pkg.go.dev/net/http#pkg-index) 作为参考说明。官方文档在布局和排版时会以<b>相关性、可读性和一致性</b>为原则，即以逻辑而不是语言规则为主：

- 例如 Constants 中会将常用的 HTTP 方法名称作为一组常量，使用块语句（block）声明；
- 例如 Variables 部分会将 DefaultServeMux 单独成行；
- Functions 中以包中的“函数”为主，只要被纳入其中必然是“函数”而不是方法；
- Types 部分将会以自定义类型为核心，将相关的类型声明、“函数”和“方法”汇总在一起，例如 `type Request` 中的 [`func NewRequest`](https://pkg.go.dev/net/http#NewRequest) 使用“函数”而不是其接收器方法，但是由于该函数属于类型的实例化，所以将其合并在一起。同样的情况还有 `type Handler interface` 其包含的函数都能返回该接口类型，例如 `func FileServer(root FileSystem) Handler`，所以也汇总在一起。 
