package main   // 表示一个可独立执行的程序，每个 Go 程序都包含一个 main 包
import "fmt"   // fmt 包（即 format）实现了格式化 IO（输入/输出）的函数

func main() {  // 如果有 init() 函数则会先执行该函数，'{' 不能独自成行
   fmt.Println("Hello, World!")
}
