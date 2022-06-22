package svc

import (
	"context"
	"fmt"
	"strings"

	"github.com/gusaul/grpcox/internal/proto/demov1pb"
)

func NewServer() demov1pb.TransportServer {
	return &server{}
}

type server struct {
}

func (s *server) Bidi(bidiServer demov1pb.Transport_BidiServer) error {
	for {
		var msg, _ = bidiServer.Recv()
		_ = bidiServer.Send(msg)
	}
}

func (s *server) ClientStream(streamServer demov1pb.Transport_ClientStreamServer) error {
	for i := 10; i > 0; i-- {
		fmt.Println(streamServer.Recv())
	}
	return nil
}

func (s *server) ServerStream(message *demov1pb.Message, streamServer demov1pb.Transport_ServerStreamServer) error {
	for i := 10; i > 0; i-- {
		message.Hello = fmt.Sprintf(strings.TrimSpace(message.Hello)+" %d-------\n", i)
		_ = streamServer.Send(message)
	}
	return nil
}

func (s *server) Unary(ctx context.Context, message *demov1pb.Message) (*demov1pb.Message, error) {
	return message, nil
}
