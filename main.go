package main

import (
	"log"
	"runtime"

	"github.com/xDarkicex/FileServer/api"
	"github.com/xDarkicex/FileServer/server"
	"github.com/xDarkicex/fileServer/datastore"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() / 2)
	datastore.ConnectPostgre()
	log.Println("Listening on localhost:8080")
	server := server.NewRouter()
	api.Start(server.Group("/api"))
	server.Static("/", "./assets/")
	log.Fatal(server.ListenAndServe("127.0.0.1:8080"))

}
