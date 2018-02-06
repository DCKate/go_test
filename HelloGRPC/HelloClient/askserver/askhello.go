package askserver

import (
	"context"
	"fmt"
	"log"

	pb "../../hello"
	"google.golang.org/grpc"
)

type HelloClient struct {
	Addr   string
	Port   int
	Client pb.GreeterClient
	Conn   *grpc.ClientConn
}

func (cc *HelloClient) GetConnectAddr() string {
	where := fmt.Sprintf("%v:%v", cc.Addr, cc.Port)
	return where
}
func (cc *HelloClient) SetConnection(conn *grpc.ClientConn) {
	cc.Conn = conn
	cc.Client = pb.NewGreeterClient(conn)
}
func (cc *HelloClient) SayHello(name string) string {
	r, err := cc.Client.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
	return r.Message
}
func (cc *HelloClient) CloseConnection() {
	defer cc.Conn.Close()
}
