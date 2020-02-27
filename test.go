package main

import "fmt"

func main() {
	test1()
}

func Foo() (err1 error) {
	if err := Bar(); err != nil {
		return
	}
	return
}

func Bar() error {
	return nil
}

//可变参数是空接口类型
func test1() {
	var a = []interface{}{1, 2, 3}

	fmt.Println(a)
	fmt.Println(a...)
}

//数组是值传递
func test2() {
	x := []int{1, 2, 3}

	func(arr []int) {
		arr[0] = 7
		fmt.Println(arr)
	}(x)

	fmt.Println(x)
}