package nums_streamer

import (
	"context"
	"fmt"
	"grpc_app/pkg/api"
	"log"
	"time"
)

const (
	singleUserLogin    = "qwerty"
	singleUserPassword = "qwerty"
	accessToken        = "accessed"
)

type GRPCServer struct {
	api.UnimplementedNumsStreamerServer
	stop chan bool
}

func NewGRPCServer() *GRPCServer {
	return &GRPCServer{}
}

func (G *GRPCServer) Authenticate(ctx context.Context, request *api.AuthenticationRequest) (resp *api.AuthenticationResponse, err error) {
	println(request.GetLogin(), " tries to connect with password: ", request.GetPassword())
	if request.GetLogin() == singleUserLogin && request.GetPassword() == singleUserPassword {
		println("success")
		return &api.AuthenticationResponse{Authorization: accessToken}, nil
	} else {
		println("access denied")
		return nil, fmt.Errorf("incorrect login or password")
	}
}

func (G *GRPCServer) StartStream(request *api.StartRequest, server api.NumsStreamer_StartStreamServer) error {
	if !checkAuthorization(request.GetAuthorization()) {
		err := fmt.Errorf("have no rights")
		return err
	}
	func() {

		var i int64
		i = 1

		for {
			select {
			case <-G.stop:
				println("Stream stopped")
				return
			case <-server.Context().Done():
				println("Stream ended")
				return
			default:
				println(i)

				err := server.SendMsg(&api.IncreasingNumber{Number: i})
				if err != nil {
					log.Print(err)
					server.SendMsg(err)
					break
				}

				time.Sleep(time.Duration(request.GetInterval()) * time.Millisecond)

			}
			i++
		}

	}()

	return nil
}

func checkAuthorization(authorization string) bool {
	return authorization == accessToken
}

func (G *GRPCServer) StopStream(ctx context.Context, request *api.StopRequest) (*api.StopResponse, error) {
	if !checkAuthorization(request.GetAuthorization()) {
		return nil, fmt.Errorf("have no rights")
	}
	G.stop <- true
	return &api.StopResponse{}, nil
}

func (G *GRPCServer) mustEmbedUnimplementedNumsStreamerServer() {
}
