package main

import (
	"fmt"
	"github.com/zmisgod/goSpider/bloom_filter"
)

func main() {
	test7()
}

func test7() {
	filter := bloom_filter.NewBloomFilter()
	fmt.Println(filter.Func[1].Seed)
	str1 := "hello,bloom filter!"
	filter.Add(str1)
	str2 := "A happy day"
	filter.Add(str2)
	str3 := "Greate wall"
	filter.Add(str3)

	fmt.Println(filter.Set.Count())
	fmt.Println(filter.Contains(str1))
	fmt.Println(filter.Contains(str2))
	fmt.Println(filter.Contains(str3))
	fmt.Println(filter.Contains("blockchain technology"))
}

func test5() {
	//ch := make(chan int)
	//go func(ch chan int) {
	//	ch <- 1
	//}(ch)
	//fmt.Println("get block channel:", <-ch)

	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	fmt.Println("get cache channel:", <-ch) //get cache channel:1
	fmt.Println("get cache channel:", <-ch) //get cache channel:2
	fmt.Println("get cache channel:", <-ch) //get cache channel:3
	close(ch)
}

func test6() {
	ch := make(chan int)
	go func() {
		ch <- 1
	}()
	select {
	case o, ok := <-ch:
		if ok {
			fmt.Println("ch = ", o)
		} else {
			fmt.Println("channel is closed")
		}
	}
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
