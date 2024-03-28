package main

import (
	"flag"
	"log"

	"github.com/spmadness/otus-go-project/internal/client"
)

var (
	M    int64
	N    int64
	Type int64
	host string
	port int64
)

func init() {
	flag.Int64Var(&N, "N", 1, "amount of seconds for stats fetch interval")
	flag.Int64Var(&M, "M", 5, "amount of seconds to use for average values")
	flag.Int64Var(&Type, "type", 0, "type of stats to get")
	flag.StringVar(&host, "host", "localhost", "server host")
	flag.Int64Var(&port, "port", 50051, "server port")
}

func main() {
	flag.Parse()
	log.Println("starting client...")
	c := client.NewClient(host, port, client.OutputTable)

	err := c.Connect()
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer c.Disconnect()

	err = c.Start(M, N, Type)
	if err != nil {
		log.Println(err.Error())
	}
}
