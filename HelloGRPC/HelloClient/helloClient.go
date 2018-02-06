package main

import (
	"../HelloClient/askserver"
)

func main() {
	var hc = askserver.HelloClient{
		Addr: "127.0.0.1",
		Port: 50051,
	}
	askserver.ConnectServer(&hc)
	hc.SayHello("GRPC")
	hc.SayHello("TryTry")
	askserver.CloseServer(&hc)
}
