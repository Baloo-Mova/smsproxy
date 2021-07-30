package main

import (
	"flag"
	"log"
	"net/http"

	"gitlab.com/devskiller-tasks/messaging-app-golang/restapi"
)

func main() {
	var port = flag.Int("port", 8080, "port")
	flag.Parse()

	server := restapi.NewServer(*port)
	server.BindEndpoints()
	if err := server.Run(); err != http.ErrServerClosed {
		panic(err)
	}
	log.Println("shutdown: completed")
}
