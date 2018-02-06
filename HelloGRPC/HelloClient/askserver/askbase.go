package askserver

import (
	"log"

	"google.golang.org/grpc"
)

type InfaceClient interface {
	GetConnectAddr() string
	SetConnection(*grpc.ClientConn)
	CloseConnection()
}

func ConnectServer(cc InfaceClient) {
	addr := cc.GetConnectAddr()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Printf("fail to dial: %v\n", err)
	}
	// defer conn.Close()
	cc.SetConnection(conn)
	log.Printf("Connect to %v\n", addr)
}

func CloseServer(cc InfaceClient) {
	addr := cc.GetConnectAddr()
	cc.CloseConnection()
	log.Printf("Close Connection %v\n", addr)
}
