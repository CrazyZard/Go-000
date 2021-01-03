package main

import (
	"log"
	"sync"
	"time"
)

type SlidingWindow struct {
	Windows map[int64]*window
	Mutex *sync.RWMutex
}

type window struct {
	Value int64
}

func newSlidingWindow() *SlidingWindow  {
	return  &SlidingWindow{
		Windows: make(map[int64]*window),
		Mutex: &sync.RWMutex{},
	}
}

func (sw *SlidingWindow)  getCurrentWindow() *window {
	now := time.Now().Unix()
	var w *window
	var ok bool
	if w ,ok = sw.Windows[now];!ok {
		w = &window{}
		sw.Windows[now] = w
	}
	return w
}

func (sw *SlidingWindow) removeOldWindows()  {
	now := time.Now().Unix() - 5
	for timestamp := range sw.Windows {
		if timestamp <= now {
			delete(sw.Windows,timestamp)
		}
	}
}

func (sw *SlidingWindow) Increment (i int64)  {
	if i == 0 {
		return
	}

	sw.Mutex.Lock()
	defer sw.Mutex.Unlock()

	b := sw.getCurrentWindow()
	b.Value += i
	sw.removeOldWindows()
}

func (sw *SlidingWindow) Sum (now time.Time) int64  {
	sum := int64(0)

	sw.Mutex.Lock()
	defer sw.Mutex.Unlock()

	for timestamp , window := range sw.Windows {
		if timestamp > now.Unix() - 5 {
			sum += window.Value
		}
	}

	return sum
}

func (sw *SlidingWindow) Avg (now time.Time) int64  {
	return sw.Sum(now) / 5
}

func main()  {
	window := newSlidingWindow()
	for _, request := range []int64{1,2,3,4,5,6,7,8,9} {
		window.Increment(request)
		time.Sleep(1 * time.Second)
	}
	log.Printf("-------- %d",window.Avg(time.Now()))
}