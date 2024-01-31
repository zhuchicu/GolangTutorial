# 04-ModuleReference

- 依赖管理
    - 依赖管理的历史
        - 依赖管理的主要内容
    - 创建 module 
        - 如何使用 module
        - go.work 文件
        - go.mod 文件
        - go.sum 文件
- 扩展内容
    - 模块感知模式
        - 设置 GOPROXY
    - go.work 文件解析
    - go.mod 文件解析
        - replace directive
        - retract directive
    - go.sum 文件解析
        - 使用 checksum database
    - Go 环境变量
    - Version
        - Pseudo-versions
        - major version suffix
    - FQA
        - Q：package 来自哪个 module
    - Reference

## 依赖管理

### 依赖管理的历史

Go 语言的依赖管理经历了三个主要阶段：GOPATH、Go Vendor 和 Go Module。

GOPATH 阶段：这是 Go语 言早期的一个依赖管理方式，它是一个环境变量，也是 Go 项目的工作区。在GOPATH下，项目源码、编译生成的库文件和项目编译的二进制文件都有特定的存放路径。然而，如果多个项目依赖同一个库，则每个项目只能使用该库的同一份代码，容易触发依赖冲突，无法实现库的多版本控制，GOPATH 管理模式就显得力不从心。

Go Vendor 阶段：在 Go 1.5 版本中推出了 vendor 机制，每个项目的根目录下有一个 vendor 目录，里面存放了该项目的依赖包。go build 命令在查找依赖包时会先查找 vendor 目录，再查找 GOPATH 目录。这解决了多个项目需要使用同一个包的依赖冲突问题。然而，如果不同工程想重用相同的依赖包，每个工程都需要复制一份在自己的vendor目录下，导致冗余度上升，无法实现库的多版本控制。

Go Module 阶段：从 Go 1.11 版本开始，官方推出 Go Module 作为包管理工具，并在 Go 1.16 版本中默认开启。在项目目录下有一个 go.mod 文件，且工程项目可以放在 GOPATH 路径之外。通过 go.mod 文件来描述模块的依赖关系，并使用 go get/go mod 指令工具来管理依赖包的下载和更新。Go Module 解决了之前存在的问题，实现了库的多版本控制，并且可以自动化管理依赖包的下载、编译和安装过程。

#### 依赖管理的主要内容

跟依赖相关（[Go Modules Reference](https://go.dev/ref/mod)）的内容主要是以下这些：

- 创建自己的包：[`go mod` 命令](https://pkg.go.dev/cmd/go#hdr-Module_maintenance)。官方文档 [Create a Go module](https://go.dev/doc/tutorial/create-module)、[Call your code from another module](https://go.dev/doc/tutorial/call-module-code)
- 管理自己包的依赖：[`go.mod` 文件](https://go.dev/doc/modules/gomod-ref)和 `go.sum` 文件。官方文档 [Managing dependencies](https://go.dev/doc/modules/managing-dependencies)
- 将自己的包公开发布：版本管理。官方文档 [Publishing a module](https://go.dev/doc/modules/publishing)，[Module version numbering](https://go.dev/doc/modules/version-numbers)


### 创建 module 

模块（Modules）是一组已发布的、带有版本控制和分布式的包（package）。它可以直接从版本控制仓库，或模块代理服务器上下载。一个模块路径（module path）标识一个模块，该路径定义在 `go.mod` 文件中，其中包含了这个模块的一些依赖，该文件位于模块的根目录下。含有 go 命令调用的模块是主模块（main module）。


#### 如何使用 module

把 golang 升级到 1.11+，设置环境变量 `GO111MODULE=on` 的值。

#### go.work 文件

go 1.18 版本引入多模块工作区，`go.work` 文件是 Go 工作空间（[workspace](https://go.dev/ref/mod#workspaces)）的配置文件。它位于工作空间的根目录下，包含的信息：工作空间的根目录、工作空间中的模块列表和模块的依赖关系。官方文档 [Getting started with multi-module workspaces](https://go.dev/doc/tutorial/workspaces)。

例如。示例如下：

~~~go
go 1.18

use (                 // 指定工作区
    ./hello           // 子模块名称为 example.com/hello
    ./example/hello   // 子模块名称为 golang.org/x/example/hello
)
~~~

`go.work` 相关的 go 命令：[`$ go work init`](https://go.dev/ref/mod#go-work-init) 创建，[`$ go work use`](https://go.dev/ref/mod#go-work-use) 添加模块，[`$ go work edit`](https://go.dev/ref/mod#go-work-edit) 底层编辑。如果没有指定工作空间根目录，那么 `go work init` 命令将使用当前目录作为工作空间根目录。


一些注意点：

- 同时使用 replaces 指令时，`go.work` 文件优先级高于 `go.mod`；
- 通常情况下，建议不要提交 `go.work` 到 git 上，因为它主要用于本地代码开发；
- 使用 `go.work` 时，不同的项目文件必须在有同一个根目录， 推荐在 `$GOPATH` 路径下执行，生成 `go.work` 文件；
- 目前仅 `go build` 会对 `go.work` 做出判断，`go mod tidy` 不会影响工作区。
- 若想要禁用工作区模式，可以通过 -workfile=off 指令来指定。

##### Q：既然已经有了 `go.mod`，为什么还需要 `go.work`？

多模块工作区能够使开发者能够更容易地同时处理多个模块的工作， 如：方便进行依赖的代码调试(打断点、修改代码)、排查依赖代码 bug 。方便同时进行多个仓库/模块并行开发调试。<b>本质上还是为了解决本地开发的诉求。</b>

在日常开发中会遇到两个经典问题：依赖本地 replace module 和依赖本地未发布的 module。

依赖本地 replace module 的场景。例如，为了解决一些本地依赖，会在 `go.mod` 文件中使用 replace 指令做一些模块的替换，这样就可以实现本地开发联调时的准确性。由于每个人 replace 后的模块路径都不一样（本地路径），导致文件被修改而被提交到仓库，这样就避免不了每次提交时都要记得修改 `go.mod` 文件中被 replace 掉的模块路径。

依赖本地未发布的 module。例如，在本地同时开发多个库时，由于存在尚未发布的库，而无法通过 `github.com/NAME/Project-NAME` 这样的模块路径 import 到依赖的模块，导致本地运行失败。为了解决这个问题，在 Go1.18 以前，我们会通过 replace 指令做一些替换（将会遭遇“依赖本地 replace module”的问题），或者先发布到可被 Go tools 拉取到的公开库上（例如 github，这样讲会导致整个开发流程更繁琐）。

##### Go modules v2

每个完整的 Go 项目中都会有 `go.mod` 文件，但只有使用 Go modules v2 的项目才会同时有 `go.work` 文件。Go modules v2 是 `Go 1.17` 中引入的新的模块系统。使用 `$ go version` 查看正在使用的 Go modules 版本。可以使用 `$ go mod init <module-name>` 命令将项目迁移到 Go modules v2。

#### go.mod 文件

当你的代码导入其他模块中的软件包（package）时，你需要通过自己的模块来管理这些依赖关系。该模块由 `go.mod` 文件定义，该文件跟踪提供这些包的模块。使用 [`$ go mod init [module-path]`](https://go.dev/ref/mod#go-mod-init) 命令来创建 `go.mod` 文件，命令中的传参是代码所在模块的名称，该名称就是模块的模块路径。官方提供了 [golang.org/x/mod/modfile](https://pkg.go.dev/golang.org/x/mod/modfile) 包来自动化解析、操作和生成 `go.mod` 文件。 注意：Go 模块路径区分大小写、不能包含空格、不能以 `.` 或 `_` 开头或结尾。

module-path 模块名称（模块路径）命名建议：

- 创建公开库，使用 `github.com` 作为前缀。例如 `github.com/johndoe`，其中 johndoe 是 GitHub 的用户名；（确保能被 Go tools 下载到）
- 创建公司内部库，使用反向域名表示法，将公司与组织的名称作为前缀。例如 `com.example/[project-name]`，其中 example 就是公司与组织的名称；
- 创建个人项目，使用 GitHub 用户名或项目名称作为前缀。例如 GitHub 用户名是 `johndoe`，那么模块路径可以是 `github.com/johndoe/my-project`。


示例如下：

~~~bash
$ go mod init example/hello   # hello 
go: creating new go.mod: module example/hello
~~~

上述示例的模块名称为 `example/hello`，前缀 `example` 是模块的作者或组织，`hello` 是项目名称，是真实的文件夹。在其内部创建 `hello.go` 文件，最终在 `./hello` 路径下，使用 [`$ go run .`](https://pkg.go.dev/cmd/go#hdr-Compile_and_run_Go_program) 命令编译并运行。


#### go.sum 文件

在模块的根目录中<b>可能</b>有一个与 `go.mod` 文件同层级的文本文件，叫 [`go.sum`](https://go.dev/ref/mod#go-sum-files)。它包含模块<b>直接和间接依赖关系</b>的加密哈希值。当 go 命令将模块 `.mod` 或 `.zip` 文件下载到模块缓存（module cache）中时，它会计算一个哈希值，并检查该哈希值是否与主模块的 `go.sum` 文件中的相应哈希值相匹配。如果模块没有依赖关系，或者所有依赖关系都被替换为本地目录，则 `go.sum` 文件可能为空或不存在。

示例如下：

~~~go
package main
import "fmt"
import "rsc.io/quote"

func main() {
    fmt.Println(quote.Go())
}
~~~

上述示例中的导入包 [rsc.io/quote](https://pkg.go.dev/search?q=quote) 是在 pkg.go.dev 中搜索到的。在使用 [`$ go run .`](https://pkg.go.dev/cmd/go#hdr-Compile_and_run_Go_program) 命令编译并运行前，需要[添加新模块要求以及验证](https://go.dev/ref/mod#authenticating)。使用 [`$ go mod tidy`](https://go.dev/ref/mod#go-mod-tidy) 添加缺失的包，或移除多余的包。在此示例中，该命令将会定位并下载 `rsc.io/quote` 包（默认是最新版本 latest）。示例如下：

~~~bash
$ go mod tidy
go: finding module for package rsc.io/quote
go: found rsc.io/quote in rsc.io/quote v1.5.2

$ go run .
~~~

## 扩展内容


### 模块感知模式

对于缺失模块的下载方式（即如何管理依赖）。

大多数 go 命令可以在<b>模块感知模式</b>（Module-aware）或 `GOPATH` 模式下运行。在模块感知模式下，go 命令使用 `go.mod` 文件来查找版本化依赖项，它通常从 module cache 中加载包，如果模块丢失则下载模块；而 `GOPATH` 模式下，go 命令会忽略模块；它会在供应商（vendor）目录和 GOPATH 中查找依赖项。可以使用环境变量 `GO111MODULE` 来设置 Module-aware 模式：

- `GO111MODULE=off` 运行 GOPATH 模式，go 命令忽略 `go.mod` 文件；
- `GO111MODULE=on` 或缺省，运行 Module-aware 模式，即使可能不存在 `go.mod` 文件；
- `GO111MODULE=off` 若当前目录或任一父级目录存在 `go.mod` 文件，则运行 Module-aware 模式。Go 1.16 以下默认该设置，Go 1.16 及以上默认 Module-aware 模式。使用 `go mod <command>` 或 `go install <version>` 时将按照第二条规则执行。

注意：在 Module-aware 模式下，GOPATH 不再表示构建期间的导入，但它仍然存储下载的依赖项（在 `GOPATH/pkg/mod` 中）和安装的命令（在 `GOPATH/bin` 中，除非设置了环境变量 GOBIN） 。

在 go 命令的选项中添加 `-mod=mod` flag，将会指示 go 命令尝试去找到一个新的模块，同时去更新 `go.mod` 和 `go.sum` 文件。`go get` 和 `go mod tidy` 命令是默认自动去寻找。

##### module cache

module cache 是存储通过 go 命令下载模块的目录，它不同于的 build cache。其默认位置是 `$GOPATH/pkg/mod`，也可以通过 go 环境变量  `GOMODCACHE` 来修改（使用 `$ go env` 查看变量具体值）。存储在缓存中的文件或目录只有“可读”权限，且无法通过 `rm -rf` 来删除，只能通过 `$ go clean -modcache` 命令来移除。使用参数 `-modcacherw` 能够创建“可读写”的目录。

注：如果是未发布的模块，则不要放置在该目录下。因为它是通过下载自动存储的。

#### 设置 GOPROXY

go 命令通过 GOPROXY 寻找含有缺失包的模块，[go 环境变量](http://docscn.studygolang.com/ref/mod#environment-variables) GOPROXY 的值由 URLs 或多个关键字 direct/off 组成（逗号分隔）：

- proxy URL 表示 go 命令使用 GOPROXY  协议访问 module proxy；
- direct 表示通过 VCS 寻找；
- off 表示不尝试访问任何服务。

~~~bash
$ go env
set GOPROXY=https://proxy.golang.org,direct  # 默认值
set GONOPROXY=  # 直接从 VCS 中获取，而不是模块代理。未设置时，GOPRIVATE 是其默认值
set GOPRIVATE=  # 决定模块是 private 还是 GOVCS
~~~~


### go.work 文件解析


`go.work` 文件的指令有：

- go 指令表示 `go.work` 文件使用的 go toolchain version。每个文件中至少有一个 go 指令；
- <b>toolchain</b> directive 声明了工作区中建议使用的 [go toolchain](https://go.dev/doc/toolchain)。只有当默认工具链比建议的工具链更早旧时，该指令才会生效；
- <b>use</b> directive 会将磁盘上的一个模块添加到工作区的主模块集合中。它的参数是一个相对路径，指向包含模块 `go.mod` 文件的目录。use 指令不会添加包含在其参数目录子目录中的模块。这些模块可以由包含其 go.mod 文件的目录在单独的 use 指令中添加。
- 

~~~go
toolchain go1.21.0

go 1.18   // go directive 

// use directive
use (     // 支持 block 语法
    ./my/first/thing
    ./my/second/thing
)

replace example.com/bad/thing v1.4.5 => example.com/good/thing v1.4.5
~~~




### go.mod 文件解析

官方文档 [`go.mod` file reference](https://go.dev/doc/modules/gomod-ref)。

模块定义在其根目录下的 `go.mod` 文件中，文件的每行都有一个单独的指令，它由关键字和紧随其后的参数们组成。为首的关键字支持 block 格式，类似导入。`go.mod` 文件被设计成人类可读、机器可写。go 提供了一些二级命令（`go mod <command>`）来修改该文件，如 `go get` 能升级或降级指定的依赖；`go mode graph` 能够在 `go.mod` 文件需要时按需升级。包 `golang.org/x/mod/modfile` 也能以代码的形式完成文件的修改。

~~~go
// Deprecated: use example.com/my/thing/v2 instead.
module example.com/my/thing   // 模块名

go 1.12

require (  // 注释，不支持“/* */”格式
    example.com/new/thing/v2 v2.3.4
    example.com/old/thing v1.2.3
)
exclude example.com/old/thing v1.2.3
replace example.com/bad/thing v1.4.5 => example.com/good/thing v1.4.5
retract [v1.9.0, v1.9.5]
~~~

##### 哪些模块需要 mod 文件？

主模块，和指定了本地路径的替换模块，需要 mod 文件。但模块缺少显式 `go.mod` 文件时，依然能够被引用，或者作为一个被指定路径与版本的替换模块使用（详见 [Compatibility with non-module repositories](http://docscn.studygolang.com/ref/mod#non-module-compat)）。

##### 文件中的关键词：指令与废弃注释

文件中的每一行就是一个单独的指令，指令中包含的关键词：

- module 用于定义主模块路径，一个 `go.mod` 中有且仅有一个该指令；
- go 表明模块是基于哪个版本的 go 编写，主要是用于表明不同版本下对新特性的支持。必须是有效的 go 版本号，如 `go 1.14`。这个指令最多只有一个，若缺省时 go 命令将会添加当前 go 版本作为默认；
- require 声明一个给定模块依赖的最低要求版本。如果解析的包不是由主模块中的包导入的，则 go 命令会自动给新 requirement 添加 `//indirect` 注释。；
- exclude 可以防止一个模块的版本被 go 命令加载，且只应用于主模块，在其他模块中将被忽略；
- replace 模块替换。同样仅应用于主模块：
- retract 表示不应该依赖 `go.mod` 所定义的模块的某个版本或一系列版本。

一个模块可以被以字符串 Deprecated 开头的注释标记为“废弃”，其废弃消息在“:” 后，以段落为结尾。这段注释在 Module 指令前，或紧跟在其后（同一行）。<b>手动添加</b>了废弃注释之后：

- Go 1.17 版本以后，可以使用 `go list -m -u` 检查 build list 中所有废弃模块的信息（是否升级、是否 retract）；也可以使用 `go get <pkg_path>` 检查指定包是否需要升级； 
- 当 go 命令识别到了模块中的废弃信息时，它会从符合 [@latest](http://docscn.studygolang.com/ref/mod#version-queries) 版本的模块中加载 `go.mod`，而不考虑 retractions 或 exclusions 指令。
- 模块作者添加了 `//Deprecated:` 注释后，要标记一个新的发布版本，并在下一个更高的版本中移除或修改废弃注释信息；
- 废弃注释将会应用于模块的所有小版本（minor）；
- 单独的小版本和补丁版本不能设置为废弃，使用 retract 指令更合适；


~~~go
// Deprecated: use example.com/mod/v2 instead.
module example.com/mod  // 或紧跟其后
~~~

#### <span id="replace_directive">replace directive</span>

示例：

~~~go
replace golang.org/x/net v1.2.3 => example.com/fork/net v1.4.5  // 替换指定的版本

replace (  // 支持 block 语句
    golang.org/x/net v1.2.3 => example.com/fork/net v1.4.5
    golang.org/x/net => example.com/fork/net v1.4.5             // 全部版本都将被替换
    golang.org/x/net v1.2.3 => ./fork/net                       // 模块路径，不能省略要替换的版本号
    golang.org/x/net => ./fork/net                              // 本地相对路径，省略要替换的版本号
)
~~~

- `=>` 左侧有版本号，表示只替换指定版本，其他版本不受影响。若没有版本号表示全部替换；
- `=>` 右侧的路径是绝对或相对路径，它表示替换模块根目录的本地文件路径，该目录必须包含一个 `go.mod` 文件。在这种情况下，需要替换的版本号必须省略；
- `=>` 右侧不是本地路径，那它则必须是有效的模块路径，在这种情况下用于替换的模块版本号要指名；

如果替换模块有 go.mod 文件，那它的 module 指令必须与其被替换的模块路径匹配。替换指令是独立的，它不会将模块添加到 [Module graph](http://docscn.studygolang.com/ref/mod#glos-module-graph) 中。在主模块或依赖模块的 `go.mod` 文件中，也需要指向被替换模块版本的 require 指令。若替换指令左侧的模块版本没有被引用，replace 指令也不会生效。


#### <span id="retract_directive">retract directive</span>

retract directive 表示不应该依赖 `go.mod` 所定义的模块的某个版本或一系列版本。retract 版本应该在版本控制库和模块代理上保持可用，以确保依赖它们的 build 不会被破坏。retact 的概念来源自学术文献：一篇被撤回的研究论文仍然可用，但它有问题，不应成为未来工作的基础。

注：Go 1.1.6 版本才添加了 retract directive，若在低于该版本的主模块 `go.mod` 文件中添加 retract 指令的话，go 将会报错，同时也会忽略依赖项的撤回指令。

##### retract directive 示例

示例，模块 `example.com/m` 的作者意外的发布了版本 v1.0.0，为了阻止用户们将模块版本升级到 v1.0.0，作者在 `go.mod` 文件中添加了两条 retract 指令，然后 tag 了包含这些撤回的 v1.0.1 版本：

~~~go
retract (
    v1.0.0 // Published accidentally.
    v1.0.1 // Contains retractions only.
)
~~~

retract 指令的形式可以是指定单一版本的（如 v1.0.0），或使用“[]”包围的多个版本（如 `[v1.1.0, v1.2.0]`），v1.0.0 等价于 `[v1.0.0]`。retract 也可以使用 block 语句将多条指令成组，如上示例。建议每个 retract 指令都添加一条解释撤回原因的注释，它们可以在 `go list` 中输出用以提醒。注释书写的格式类似 Module 指令中“废弃注释”，位于 retract block 上的注释将会应用于块内所有的指令。


##### 如何设置 retract 版本？

模块作者通过在 `go.mod` 文件中添加 retract 指令实现版本的撤回，同时发布一个包含这个指令的新版本。新版本必须高于其他发布或预发布版本，即 @latest 版本查询应解析已被 retract 之后的新版本（如上述示例）。使用 go 命令 `go list -m -retracted $modpath@latest` 加载并将所显示的模块版本的设置为 retract 版本。

一个包含 retract 指令的版本可能会自行撤回。如果模块的最高版本或预发布版本自行撤回时，@latest 查询版本时会解析一个（在排除所有撤回版本后）最高版本。如上述示例，当用户运行 `go get example.com/m@latest` 时，go 命令会读取 v1.0.1（现在是最高版本）的撤回指令。 v1.0.0 和 v1.0.1 都已撤回，因此 go 命令将升级（或降级）到下一个最高版本，可能是 v0.9.5。


##### go 命令对含有 retract 指令的模块差异性

当一个模块的版本被 retract 时，用户将无法通过 `go get`、`go mod tidy` 或其他命令，自动升级到它。依赖于 retract 版本的构建应该继续工作，但是当用户用 `go list -m -u` 检查更新或用 `go get` 更新相关模块时，用户会被通知撤回。

使用 `go list -m -versions` 打印模块版本时不会显示 retract 版本，若要显示需要添加 `-retracted` flag。当版本查询的查询语句类似 `@>=v1.2.3` 或 `#latest` 时，retract 版本也会被排除。



### go.sum 文件解析

模块根目录下有一个与 `go.mod` 并排的文本文件 `go.sum`，其内容是模块直接和间接依赖的加密哈希值。当 go 命令下载一个模块的 `.mod` 或 `.zip` 文件到 [module cache](http://docscn.studygolang.com/ref/mod#module-cache) 时，它会计算一个哈希值，并检查该哈希值是否与主模块的 `go.sum` 中的相应哈希值相匹配。`go.sum` 内容可为空，或不存在该文件：当模块没有依赖关系，或所有依赖都被替换成本地目录（使用 replace directive）。使用 `go mod verify` 命令来验证 `go.sum` 文件的完整性。如果 go.sum 文件被篡改，则 `go mod verify` 命令将报告错误。

`go.sum` 文件的每一行都有三段字段域（由空格分隔），分别是 Module path、version 和 hash：

- Module path 是哈希所属模块的名称；
- version 是指模块所属版本。若版本号以 `/go.mod` 结尾，表示 哈希值仅代表模块 `go.mod` 文件；否则代表模块的 zip 文件；
- 哈希列由算法名称（如 h1）和 Base64 编码的加密哈希组成，并用冒号 (:) 分隔。当然，目前 SHA-256 (h1) 是唯一支持的哈希算法。如果将来发现 SHA-256 中的漏洞，将添加对另一种算法（如名为 h2 的等等）的支持。

示例如下：

~~~go
<module> <version>[/go.mod] <hash>
github.com/google/uuid v1.3.0 h1:TIyPBB2g7aqYIRf/OVQ6dmjwBe3QDzSwfl9Qjot/M04=
github.com/google/uuid v1.3.0/go.mod h1:TIyPBB2g7aqYIRf/OVQ6dmjwBe3QDzSwfl9Qjot/M04=
~~~


`go.sum` 文件可能包含模块多个版本的哈希值。 go 命令可能需要从依赖项的多个版本加载 go.mod 文件，以便执行最小版本选择。 `go.sum` 可能会包含不再需要的模块版本的哈希值（例如升级后），而 `go mod tidy` 命令将添加缺失的哈希值，并从 `go.sum` 中删除不必要的哈希值。

#### 使用 checksum database

如果 `go.sum` 文件不存在，或者它不包含下载文件的哈希值，则 go 命令可以使用 [checksum database](http://docscn.studygolang.com/ref/mod#checksum-database)（它是公共可用模块的哈希值的全局源）来验证。一旦验证了哈希值，go 命令会将其添加到 `go.sum` 并将下载的文件添加到 Module cache 中。如果模块是私有的（由环境变量 GOPRIVATE 或 GONOSUMDB 匹配），或者如果禁用校验和数据库（通过设置 `GOSUMDB=off`），则 go 命令接受哈希并将文件添加到模块缓存而不验证它。

### Go 环境变量

官方文档 [environment-variables](https://go.dev/ref/mod#environment-variables)

### <span id="Version">Version</span>

[版本的语义规范](https://semver.org/lang/zh-CN/)

版本格式：`vx.y.z`（如 `v0.0.0` ）由三个非负整数组成，分别表示主版本 major、小版本 minor、补丁版本 patch。补丁版本可以与 pre-release 组合，如 `v8.0.5-pre`；补丁版本或 pre-release 后也可跟元数据标识，如 `v2.0.9+meta`。

##### 版本号的三个整数变化规则是什么？

- 当模块的功能无法向后兼容时，major 版本必须递增，minor 和 patch 版本设置为 0，如某个包被移除；
- 当有可兼容的修改时，minor 递增，path 设为 0，如新增了某函数；
- 当有不影响公共接口的修改时，path 递增，如 bug 修复或优化；
- pre-release 后缀表示预发布版本，它将归类到相应的发布版中，如 `v1.2.3-pre` 来自 `v1.2.3`；
- 生成元数据（build metadata）后缀，主要用于忽略版本比较，但它会按照指定版本保存在 go.mod 文件中；`+incompatible` 后缀表示模块之前发布的版本已迁移到主版本2或更高。


主版本为 0 或含有 pre-release 后缀的版本，被认为是不稳定版本。不稳定的版本不受兼容性要求，如 `v0.2.0` 不兼容 `v0.1.0`。



#### <span id="Pseudo-versions">Pseudo-versions</span>


##### 什么是伪版本？

伪版本（Pseudo-versions）是版本控制系统（如 git，简称 VCS）所生成的预发布版本，其名称格式为修订标识符和时间戳格式，如 `$ go get -d golang.org/x/net@daa7c041` 将会生成伪版本 `v0.0.0-20191109021931-daa7c04131f5`，其中修订标识符 daa7c04131f5 是 git commit 哈希值。

伪版本的名称格式（`base-timestamp-revision`）：

- base version 前缀，如 `vx.0.0` 或 `vX.Y.Z-0`；
- timestamp 时间戳，格式：yyyymmddhhmmss。在 git 中是提交时间，而不是作者时间；
- 修订标识符，如 abcdefabcdef，是 commit hash 的前 12 个字符，或者由 0 填充的修订字符子版本。


##### 伪版本是如何产生的？

尽管 go 在 VCS 中可以通过 tag、branch、version 来访问那些不满足版本语义规范的模块。但有例外，当<b>主模块</b>的版本名称也不符合那些规范时，go 命令会自动修改，在这个过程中，还会移除生成元数据后缀（除了 `+incompatible` 外）。所以 go 命令的自主修改将会产生伪版本。


##### 伪版本不止一种名称格式？


由于 base version 的不同，伪版本可能有三种形式：<b>高于基础版本，低于下一个 tagged version</b>：

- 当没有已知的 base 版本时，伪版本格式：`vX.0.0-yyyymmddhhmmss-abcdefabcdef`，且主版本 X 必须匹配模块的主版本后缀；
- 当 base 版本是预发布版本时，如 `vX.Y.Z-pre`，伪版本格式：`vX.Y.Z-pre.0.yyyymmddhhmmss-abcdefabcdef`；（由 0 填充的修订字符子版本）
- 当 base 版本是正式版本，如 `vX.Y.Z`，伪版本格式：`vX.Y.(Z+1)-0.yyyymmddhhmmss-abcdefabcdef`。

##### 多种伪版本号格式的作用？

当生成伪版本后，标记了（tagged）较低版本时，会发生一个以上的伪版本可以通过使用不同的 base 版本来引用同一个提交。上述三种伪版本形式带来的好处是：具有已知 base version 的伪版本排序高于这些版本，但低于其他后续的预发布版本；有相同 base version 前缀的伪版本们将按照时间排序。

##### 为什么要设计伪版本号规范？如何检查是否符合规范？

go 命令执行多项检查，以确保模块作者可以控制如何将伪版本与其他版本进行比较，并且伪版本实际上引用的是模块提交历史记录一部分的修订版本：

- 如果指定了一个 base version，就必须有一个相应的语义版本标签，它是伪版本所描述的修订版的祖先。这可以防止开发人员使用比所有标记版本（tagged versions）更高的伪版本来绕过最小版本选择（[Minimal version selection](http://docscn.studygolang.com/ref/mod#minimal-version-selection)），例如  `v1.999.999-99999999999999-daa7c04131f5`；
- 时间戳必须与修订版本的时间戳匹配。这可以防止攻击者用无限数量的相同伪版本淹没模块代理。这也可以防止模块使用者更改版本的相对顺序；
- 该修订版本必须是模块库的某个分支或标签的祖先。这可以防止攻击者引用未经批准的修改或 pull 请求。


##### 命令自动生成伪版本号？

伪版本的生成格式不需要手动输入，很多命令都支持自动将 commit hash 或 branch name 转码为版本号：

~~~bash
$ go get -d example.com/mod@master
$ go list -m -json example.com/mod@abcd1234
~~~



#### <span id="major_version_suffix">major version suffix</span>

从主版本2开始，模块路径必须有一个主版本后缀（如`/v2`），与主版本相匹配。例如，如果一个模块 `v1.0.0` 版本的路径是 `example.com/mod`，而 `v2.0.0` 必须是 `example.com/mod/v2`。


##### 为什么要设计主版本后缀规则？

由于主版本遵循导入兼容性规则（[import compatibility rule](https://research.swtch.com/vgo-import)）：如果旧包和新包具有相同的导入路径，则新包必须向后兼容旧包，所以从 v2 版本开始，新的导入路径可以避免这个兼容规则。

主版本后缀规则，能让多个主版本能够共存在同一个 Build 中，它很好的解决了“菱形引用问题”（[diamond dependency problem](https://research.swtch.com/vgo-import#dependency_story)）。通常，如果一个模块的两个不同版本都被依赖，那么将使用较高的版本。若两个版本不兼容，那么都将无法满足需求。由于不兼容的版本必须有不同的主版本号，它们也必须有不同的模块路径（因为主版本后缀）。这就解决了冲突：具有不同后缀的模块被视为独立的模块，（即使这些包所在子文件的路径，相对于模块根目录是相同的）它们也是不同的。


##### 为什么 v0 和 v1 版本不需要添加模块路径后缀？

因为 v0 版本原本就认定为不稳定版本，它不需要考虑向后兼容问题。对于大多数模块来说，v1 版本是向后兼容上一个 v0 版本的；v1 版本是对兼容性的一种承诺，而不是表明与 v0 相比有不兼容的变化。


##### 是否有不遵循主版本后缀规则的例外情况？

当然这种规则也有特例：以 `gopkg.in/` 开头的模块路径无论哪个版本，都必须要带主版本后缀，且后缀是以 `.` 开始，而不是 `-`，如 `gopkg.in/yaml.v2`。

许多 Go 项目在迁移到模块之前，发布的正式版本都是 v2 或更高，但它们没有使用主版本后缀。这些版本都将会标注上 `+incompatible` 构建标签，如 `v2.0.0+incompatible`。


### FQA

#### Q：package 来自哪个 module

当 go command 使用 package path 加载包时，需要确定是哪个模块提供的。


##### 如何找到包的模块来源？

go 命令首先在 Build list 中列罗的模块路径搜寻，它们的路径前缀将是包的路径。如导入了包 `example.com/a/b`，而 build list 中也有 `example.com/a` 模块，那么 go 命令将会检查该模块路径下文件夹 b 是否包含包。若文件夹下有至少一个 `.go` 文件，则该文件夹将被当做一个包。

- 若 build list 下某模块含有目标包，那么该模块将被认为是真正使用者；
- 若一个都没有，或超过两个及以上的模块，那么 go 命令将会报错。


##### 找不到包模块来源，怎么办？

在 go 命令的选项中添加 `-mod=mod` flag，将会指示 go 命令尝试去找到一个新的模块，同时去更新 `go.mod` 和 `go.sum` 文件。`go get` 和 `go mod tidy` 命令是默认自动去寻找。

go 命令通过 GOPROXY 寻找含有缺失包的模块，[go 环境变量](http://docscn.studygolang.com/ref/mod#environment-variables) GOPROXY 的值由 URLs 或多个关键字 direct/off 组成（逗号分隔）：

- proxy URL 表示 go 命令使用 GOPROXY  协议访问 module proxy；
- direct 表示通过 VCS 寻找；
- off 表示不尝试访问任何服务。

~~~bash
$ go env
set GOPROXY=https://proxy.golang.org,direct  # 默认值
set GONOPROXY=  # 直接从 VCS 中获取，而不是模块代理。未设置时，GOPRIVATE 是其默认值
set GOPRIVATE=  # 决定模块是 private 还是 GOVCS
~~~~


##### 基于 GOPROXY 找到的候选模块如何处理？

GOPROXY list 将会罗列符合条件的实体，每个实体中将会罗列含缺失包的模块路径。对于拥有多个版本的模块，go 命令将优先选择其最新版本，具体筛选规则如下：

- 若有一个及以上的模块都符合，优先选路径最长的模块；
- 若有符合的模块，但不包含缺失包，go 命令将会报错；
- 若没有合格模块，go 命令尝试在下一个实体中寻找；
- 若所有实体中都没有符合的，go 命令将会报错。

当找到合适的模块后，go 命令会将新的 requirement 以及新模块的路径和版本添加到主模块的 `go.mod` 文件中。保证以后加载相同包时，将会适用相同版本的模块。如果解析的包不是由主模块中的包导入的，则 go 命令会自动给新 requirement 添加 `//indirect` 注释。

##### 定位包的模块来源示例

假设 go 命令正在寻找含包路径 `golang.org/x/net/html` 的模块，环境变量 `GOPROXY=https://corp.example.com,https://proxy.golang.org`，go 命令的具体请求过程：首先访问 `https://corp.example.com`（并行）：

- 请求最新版本的 `golang.org/x/net/html`；
- 请求最新版本的 `golang.org/x/net`；
- 请求最新版本的 `golang.org/x`；
- 请求最新版本的 `golang.org`；

若上述请求都返回 404 或 410，则访问 `https://proxy.golang.org`，并行请求上述模块的最新版本。




### Reference


1. [深入理解 Go Modules：高效管理你的 Golang 项目依赖](https://juejin.cn/post/7309692103055507491)
