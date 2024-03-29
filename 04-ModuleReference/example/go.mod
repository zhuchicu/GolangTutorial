module github.com/zhuchicu/GolangTutorial/04-ModuleReference/example

go 1.18

require github.com/zhuchicu/GolangTutorial/04-ModuleReference/mypkg v1.0.0
replace github.com/zhuchicu/GolangTutorial/04-ModuleReference/mypkg => ../mypkg
