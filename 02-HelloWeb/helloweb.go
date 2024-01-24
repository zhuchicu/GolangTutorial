package main

import (        // 嵌套结构导入
    "fmt"
    "net/http"  // 官方 package
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {  // 匿名函数，可以直接调用，或赋值给一个变量
        fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
    })

    http.ListenAndServe(":1234", nil)
}
