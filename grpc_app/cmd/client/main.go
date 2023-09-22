package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc_app/pkg/api"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	fmt.Println("Client initialization")

	fmt.Print("Enter login:")
	var login string
	fmt.Fscan(os.Stdin, &login)

	fmt.Print("Enter password:")
	var password string
	fmt.Fscan(os.Stdin, &password)

	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	c := api.NewNumsStreamerClient(conn)
	resp, err := c.Authenticate(context.Background(), &api.AuthenticationRequest{Login: login, Password: password})
	if err != nil {
		println("access denied")
		log.Fatal(err)
	}
	println("Authenticated")
	authorization := resp.Authorization

	fmt.Print("Enter interval between responses(ms):")
	var interval int64
	fmt.Fscan(os.Stdin, &interval)

	fmt.Print("Enter streaming time(ms):")
	var streamingTime int64
	fmt.Fscan(os.Stdin, &streamingTime)

	fmt.Print("Enter time until printing output(ms):")
	var timeUntilOutput int64
	fmt.Fscan(os.Stdin, &timeUntilOutput)

	var command string
	for {
		fmt.Fscan(os.Stdin, &command)
		switch command {
		case "/start":

			stream, err := c.StartStream(context.Background(), &api.StartRequest{Authorization: authorization, Interval: interval})
			if err != nil {
				log.Fatal(err)
			}
			println("Stream started")
			respChan := make(chan map[time.Time]int64)
			go func(ch chan map[time.Time]int64) {
				stop := make(chan bool)
				go func(stop chan bool) {
					time.Sleep(time.Duration(streamingTime) * time.Millisecond)
					stop <- true
				}(stop)

				buf := make(map[time.Time]int64)
				for {
					select {
					default:
						res, err := stream.Recv()
						if err == io.EOF {
							break
						}
						if err != nil {
							log.Fatal(err)
						}
						buf[time.Now()] = res.GetNumber()
					case <-stop:
						err := stream.CloseSend()
						if err != nil {
							println("error")
							respChan <- buf
							return
						}
						println("error is nil")
						respChan <- buf
					}
				}
			}(respChan)

			go func(chan map[time.Time]int64) {
				time.Sleep(time.Duration(timeUntilOutput) * time.Millisecond)
				buf := <-respChan
				for t, num := range buf {
					println(t.String(), num)
				}
			}(respChan)

		}

	}
}
