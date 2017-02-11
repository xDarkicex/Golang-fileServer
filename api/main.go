package api

import "github.com/xDarkicex/FileServer/server"

func Start(r *server.Router) {
	r.OPTIONS("/{rest:.*}", server.ContructRequest)
	handleFiles(r.Group("/files"))
	r.AllRoutes(server.DestructRequest)
}
