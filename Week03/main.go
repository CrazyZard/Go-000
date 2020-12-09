package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main()  {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	g, _ := errgroup.WithContext(ctx)

	//
	g.Go(func() error{
		if err := SignalServer(ctx);err != nil {
			cancel()
			return err
		}
		return nil
	})

	g.Go(func() error {
		if err := HttpServer(ctx,":8080",nil);err != nil {
			cancel()
			return err
		}
		return nil
	})

	if err:= g.Wait(); err!=nil{
		fmt.Printf("service is close , reasan is  %v \n",err)
	}
}

func HttpServer(ctx context.Context,addr string, handler http.HandlerFun) error {
	s := http.Server{
		Addr:    addr,
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(w, "Hello World, %s\n", srvName)
	})
	go func() {
		<- ctx.Done()
		_ = s.Shutdown(ctx)
	}()

	err := s.ListenAndServe()
	return err
}


func SignalServer(ctx context.Context) error {
	osChan := make(chan os.Signal, 1)
	println("---- start listen sign ----")
	signal.Notify(osChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		return nil
	case call := <-osChan:
		fmt.Printf("handle signal: %d \n", call)
		return nil
	}

}
