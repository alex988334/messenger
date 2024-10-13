package main

import (
	"flag"
	"fmt"

	wss_server "github.com/alex988334/messenger/pkg/messenger"
)

//	"net/http"

var (
	host     = flag.String("host", "localhost", "http service address") //	флаг консольной команды
	port     = flag.String("port", "25550", "an int")                   //	флаг консольной команды
	localUrl = flag.String("local-url", "", "a string")                 //	флаг консольной команды
// head head_unit.Head
)

func main() {

	defer fmt.Println("Exit programm")

	flag.Parse() //	расшифровываем наш флаг

	fmt.Println("host:", *host, "port:", *port)
	fmt.Println("local url:", *localUrl)
	fmt.Println("\n")

	for {
		head := wss_server.NewWssServer()
		head.Run()
	}
}
