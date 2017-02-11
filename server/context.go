package server

import (
	"net/http"

	"encoding/json"

	"log"

	"github.com/gorilla/mux"
)

type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{w, r}
}

func (c *Context) Header(name string, value string) {
	c.ResponseWriter.Header().Set(name, value)
}

func (c *Context) Param(name string) string {
	Params := mux.Vars(c.Request)
	return Params[name]
}

func (c *Context) RenderError(status int, err error) {
	http.Error(c.ResponseWriter, err.Error(), status)
}

func (c *Context) setStatus(status int) {
	c.ResponseWriter.WriteHeader(status)
}

func (c *Context) RenderJSON(status int, JSON interface{}) {
	c.Header("Content-Type", "application/json")
	c.setStatus(status)
	data, err := json.Marshal(JSON)
	if err != nil {
		c.setStatus(http.StatusInternalServerError)
		log.Printf("can not Marshal raw data: %s\nError messaage: %s", data, err.Error())
		return
	}
	c.ResponseWriter.Write(data)
}
