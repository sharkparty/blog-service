package main

import (
	blogProto "blog-service/rpc/blog"
	config "blog-service/config"

	"blog-service/server"
	"fmt"
	"net/http"
)



func main() {
	fmt.Println("Connecting to DB")
	config.ConnectToDB()
	fmt.Println("Starting Blog Service")
	server := &server.Server{}
	handler := blogProto.NewBlogServiceServer(server)
	fmt.Printf("Service listening on port: %v\n", config.Port)
	listener := fmt.Sprintf(":%v", config.Port)
	http.ListenAndServe(listener, handler)
}