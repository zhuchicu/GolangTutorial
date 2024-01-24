package main
import "fmt"

// 封装，Animal 表示 class 
type Animal struct {
	Name string  // 首字母大写是 exported name，属性为 public
	bark string "how it call" // unexported name，属性为 private
}

// 表示 class 对应的方法，接收器为 *Animal
func (a *Animal) Say(msg string) {
	fmt.Printf("Animal[%v] [%v]: \"%v\"\n", a.Name, a.bark, msg)  // %v the value in a default format
}

// 继承
type Unicorn struct {
	Animal  // 嵌入类型 embedding types
	wings int
}

// 多态
type Fly interface {  // 接口
	GetWingNum() int  // 抽象 public 方法 
}

func (u *Unicorn) GetWingNum() int {  // 实现了接口抽象方法 GetWingNum
	return u.wings
}

func main() {
	u := &Unicorn{Animal{"Pony", "whinny"}, 2}  // 接收器为 *T，所以要取地址
	u.Say("Hi!")                                // 字段提升，直接引用
	fmt.Printf("[%v] has [%v] wings.\n", u.Animal.Name, u.GetWingNum())
}