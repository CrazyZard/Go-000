package main

import (
	"context"
	"fmt"
	"time"
)
/**
	思路要点
	1. 在调用者 main 里面去启动一个goroutine
	2. 启动 tracker Run() 方法里去循环 ch 里的值
	3. 当 context 结束时 ，调用 tracker 的 shutdown 方法 关掉 ch
	4. 关掉 ch 的时候 Run() 里循环就退出了 stop
 */

func main()  {
	tr := NewTracker()
	go tr.Run()
	_ = tr.Event(context.Background(),"test")
	_ = tr.Event(context.Background(),"test")
	_ = tr.Event(context.Background(),"test")

	ctx ,cancel := context.WithDeadline(context.Background(),time.Now().Add(10 * time.Second))
	defer cancel()
	tr.stop <- struct{}{}
	tr.Shutdown(ctx)
}

func NewTracker() *Tracker {
	return &Tracker{
		ch : make(chan string,10),
		stop : make(chan struct{}),
	}
}

type Tracker struct {
	ch chan string
	stop chan struct{}
}

func (t *Tracker) Event (ctx context.Context,data string) error  {
	select {
	case t.ch <- data:
		return nil
	case <- ctx.Done():
		return ctx.Err()
	}
}

func (t *Tracker) Run()  {
	for data := range t.ch{
		time.Sleep(1 * time.Second)
		fmt.Println(data)
	}
	t.stop <- struct{}{}
}

func (t *Tracker) Shutdown(ctx context.Context)  {
	//close(t.ch)
	select {
	case <- t.stop:
		fmt.Println("stop")
	case <- ctx.Done():
		fmt.Println("ctx")
	}
}