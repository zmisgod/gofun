package worker

import (
	"context"
	"errors"
	"log"
	"math"
	"sync"
)

type w struct {
	QueueNumber int      `json:"queue_number"`
	List        []string `json:"list"`
	_func       func(ctx context.Context, item string) error
}

func NewWorker(num int) *w {
	return &w{
		QueueNumber: num,
	}
}

func (a *w) Enqueue(item string) {
	a.List = append(a.List, item)
}

func (a *w) SetFunc(_func func(ctx context.Context, item string) error) {
	a._func = _func
}

func (a *w) Handle(ctx context.Context) error {
	if a._func == nil {
		return errors.New("参数错误")
	}
	_list := ChunkStrArray(a.List, a.QueueNumber)
	wg := sync.WaitGroup{}
	for _, v := range _list {
		_oneList := v
		wg.Add(1)
		go func(_list []string) {
			defer wg.Done()
			for _, j := range _list {
				if err := a._func(ctx, j); err != nil {
					log.Println(err, j)
				}
			}
		}(_oneList)
	}
	wg.Wait()
	return nil
}

func ChunkStrArray(list []string, chunkSize int) [][]string {
	length := len(list)
	chunks := int(math.Ceil(float64(length) / float64(chunkSize))) //each count
	count := 0
	result := make([][]string, 0)
	for i := 0; i < chunkSize; i++ {
		rows := make([]string, 0)
		for g := 0; g < chunks; g++ {
			if count < length {
				rows = append(rows, list[count])
				count++
			}
		}
		if len(rows) > 0 {
			result = append(result, rows)
		}
	}
	return result
}

//checkTaskIsRunning  检查任务是否还在运行
func (a *w) checkTaskIsRunning(taskId int) bool {
	return false
}

// getATask 获取一个任务
func (a *w) getATask() (string, error) {
	return "", nil
}

//checkTaskIsEmpty 检查当前任务队列是否未空
func (a *w) checkTaskIsEmpty() bool {
	return false
}

// 1,2,3,4,5
// 任务1：1，3，5
// 任务2：2，4
