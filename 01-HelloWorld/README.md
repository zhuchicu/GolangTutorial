# 01-HelloWorld

## 具体步骤

### 在本地部署运行

1. [Install](https://golang.org/dl/) Golang（BSD-style license）。使用 `go version` 命令查看版本。
2. 在用户主目录中新建 go 文件（以 `.go` 结尾），在 `helloworld.go` 文件中编写代码。[Editors and IDEs for Go](https://github.com/golang/go/wiki/IDEsAndTextEditorPlugins)，例如 [vscode-go](https://github.com/golang/vscode-go)，[sublime-build](https://github.com/golang/sublime-build)，[vim-go](https://github.com/fatih/vim-go)；
3. 使用 `go mod init` 命令来创建 `go.mod` 文件，它用于开启代码依赖模块的跟踪。创建成功后，它将位于你代码的根目录位置。
5. 在文件目录下使用 `go run` 命令运行。

#### 命令行操作

~~~bash
$ cd %HOMPATH%   # cd ~
$ mkdir hello
$ cd hello
$ touch helloworld.go

# 编写 helloworld.go

$ go mod init example/hello
$ go run .
Hello, World!
~~~


### 在线编译运行

官方提供的在线编译器 [playground](https://go.dev/play/)，注意页面底部有一些关于 playground 的注意事项，比如标准库、时间、运行限制等，具体见详情描述。

## 扩展内容

### 管理 golang 编译器

在一台机器上安装多个版本的 golang 编译器：[Managing Go installations](https://go.dev/doc/manage-install)
卸载 go 编译器：[uninstalling](https://go.dev/doc/manage-install#uninstalling)

### 如何在 windows 下部署环境

略

### go build/install/run 区别？

Helloworld 示例中的 `go run` 命令等同于下面的命令：

~~~go
$ go run helloworld.go
Hello, World!

// go run 等同于下列的执行

$ go build helloworld.go
$ ls
helloworld    helloworld.go

$ ./helloworld
Hello, World!
~~~

通过 `$ go` 查看每个命令的释义：

*   go build（`compile packages and dependencies`）：用于测试编译包，在项目目录下生成可执行文件（有main包）。
*   go install（`compile and install packages and dependencies`）：主要用来生成库和工具。一是编译包文件（无main包），将编译后的包文件放到 pkg 目录下（`$GOPATH/pkg）。二是编译生成可执行文件（有main包），将可执行文件放到 bin 目录（$`GOPATH/bin）。
*   go run（`compile and run Go program`）：编译并运行可执行文件。

[compile-install](https://golang.google.cn/doc/tutorial/compile-install)：build 命令将源码编译为可执行程序，但是如果想要正确运行的话，需要在可执行程序的目录下，或指定执行的路径才行。如果想要不指定路径就能正确执行的话，需要安装（install）这个可执行程序（指设定环境变量 `GOBIN` 的路径）。

使用 `$ go list -f '{{.Target}}'`，查看 install 命令对于当前包的安装路径，`C:\Users\zhuml1\go\bin\hello.exe` 表示二进制程序已经被安装到该路径下。使用下面命令，能不指定路径就执行应用：

```go
// ============ set GOBIN=/path/to/your/install/dir =============
$ export PATH=$PATH:/path/to/your/install/directory  // Linux or Mac
$ set PATH=%PATH%;C:\path\to\your\install\directory  // windows

// ============ if GOBIN=$HOME/bin，change it =============
$ go env -w GOBIN=/path/to/your/bin
$ go env -w GOBIN=C:\path\to\your\bin

// ================ Run and test ============
$ go install
$ hello
```


### 为 Go 开发配置 Visual Studio Code

[为 Go 开发配置 Visual Studio Code](https://learn.microsoft.com/zh-cn/azure/developer/go/configure-visual-studio-code)：

1. [安装 Go](https://go.dev/doc/install)，打开命令提示符，然后运行 `go version` 以确认已安装 Go；
2. [安装 Visual Studio Code](https://code.visualstudio.com/)；
3. 安装 Go 扩展，活动栏“扩展”视图或者快捷方式 (Ctrl+Shift+X) 。搜索 Go 扩展，然后选择“安装”；
4. 更新 Go 工具。命令面板“帮助>显示所有命令”，或快捷方式 (Ctrl+Shift+P)，`Go: Install/Update tools` 搜索 ，然后从托盘运行命令。出现提示时，选择所有可用的 Go 工具，然后单击“确定”。等待 Go 工具完成更新；
5. 编写示例 Go 程序。新建 `main.go`，打开终端 “终端 > 新建终端”，然后运行命令 `go mod init sample-app` 来初始化示例 Go 应用（初始化 `go.mod` 文件）；
6. 运行调试器。创建断点，活动栏“调试”视图，或快捷方式 (Ctrl+Shift+D) 。单击“ 运行并调试”或按 F5，鼠标悬停断点处的变量上以查看其值。 单击调试器栏上的“继续”或按 F5 退出调试器。


~~~go
package main

import "fmt"

func main() {
    name := "Go Developers"
    fmt.Println("Azure for", name)   // breakpoint
}
~~~

Go 工具中包含包 [gopls](https://pkg.go.dev/golang.org/x/tools/gopls) (pronounced "Go please") 是 Go 团队开发的官方 Go 语言服务器。它为任何 [LSP](https://microsoft.github.io/language-server-protocol/) 兼容编辑器提供 IDE 功能。语言服务器协议 (LSP，Language Server Protocol) 定义了编辑器或 IDE 与语言服务器之间使用的协议，该协议提供自动补全、转到定义、查找所有引用等语言功能。
