package main

import (
	"flag"
	"log"
	"server"
)

func main() {
	var listenAddr string
	flag.StringVar(&listenAddr, "listen-addr", ":12177",
		"address to listen")
	flag.Parse()
	log.Printf("listening on %q", listenAddr)

	server.RunServer(listenAddr)
}
