package main

import (
	"bufio"
	"context"
	"log"
	"net"
	"runtime"
)
//发送数据
func sendMessage(ctx context.Context,conn net.Conn, ch <-chan string) {

	wr := bufio.NewWriter(conn)
	for {
		select {
		case <- ctx.Done():
			log.Printf("writer ctx err %+v", ctx.Err())
			log.Printf("writer done")
			log.Printf("Number of active goroutines %d", runtime.NumGoroutine())
			return
		case line := <- ch:
			log.Printf("发送数据：%s",line)
			wr.Write([]byte(line))
			wr.Flush()
		}
	}
}

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatalf("listen error: %v\n", err)
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("accept error: %v\n", err)
			continue
		}
		// 开始goroutine监听连接
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {

	defer conn.Close()

	channel := make(chan string, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go sendMessage(ctx,conn,channel)

	//启动接受 conn 的协程
	input := bufio.NewScanner(conn)
	for input.Scan() {
		channel <- input.Text()
	}

	if err := input.Err(); err != nil {
		log.Println("读取错误：", err)
	}

}



