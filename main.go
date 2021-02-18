package main

import (
	"context"
	"fmt"
	"github.com/zmisgod/goSpider/bloom_filter"
	"github.com/zmisgod/goSpider/music"
	"github.com/zmisgod/goSpider/utils"
)

func main() {
	res, err := music.NewFetchMusic("王心凌《大眠》 https://c.y.qq.com/base/fcgi-bin/u?__=d3IYVRj @QQ音乐")
	if err != nil {
		utils.CheckError(err)
	}else{
		resUrl, err := res.GetDownloadURL(context.Background())
		fmt.Println(resUrl, err)
	}
}

func test7() {
	filter := bloom_filter.NewBloomFilter()
	fmt.Println(filter.Func[1].Seed)
	str1 := "my name is zm"
	filter.Add(str1)
	str2 := "my name is ad"
	filter.Add(str2)
	str3 := "cool boys"
	filter.Add(str3)

	fmt.Println(filter.Set.Count())
	fmt.Println(filter.Contains(str1))
	fmt.Println(filter.Contains(str2))
	fmt.Println(filter.Contains(str3))
	fmt.Println(filter.Contains("cool boys"))
	fmt.Println(filter.Contains("cool girls"))
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
