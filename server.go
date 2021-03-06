package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

const (
	DB_USER     = "docker"
	DB_PASSWORD = "docker"
	DB_NAME     = "docker"
)

type HTTPConn struct {
	in  io.Reader
	out io.Writer
}

func (c *HTTPConn) Read(p []byte) (n int, err error)  { return c.in.Read(p) }
func (c *HTTPConn) Write(d []byte) (n int, err error) { return c.out.Write(d) }
func (c *HTTPConn) Close() error                      { return nil }

type Handler struct {
	rpcServer *rpc.Server
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serverCodec := jsonrpc.NewServerCodec(&HTTPConn{
		in:  r.Body,
		out: w,
	})
	w.Header().Set("Content-type", "application/json")
	err := h.rpcServer.ServeRequest(serverCodec)
	if err != nil {
		log.Println(err.Error())
	}
}

func main() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	userManager := NewUserManager(db)

	server := rpc.NewServer()
	server.Register(userManager)

	userHandler := &Handler{
		rpcServer: server,
	}
	http.Handle("/rpc", userHandler)

	fmt.Println("starting server at :8081")
	panic(http.ListenAndServe(":8081", nil))
}
