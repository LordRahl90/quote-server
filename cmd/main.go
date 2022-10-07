package main

import (
	"fmt"
	"net"

	"github.com/LordRahl90/quote-server/proto"
	"github.com/LordRahl90/quote-server/server"
	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:50056")
	if err != nil {
		panic(err)
	}

	// create a new server
	s := grpc.NewServer()

	// create a quote server object
	quoteServer := &server.QuoteServer{}

	proto.RegisterQuoteServer(s, quoteServer)

	fmt.Printf("Server started\n")

	if err := s.Serve(listener); err != nil {
		panic(err)
	}
}
