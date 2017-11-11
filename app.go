package main

import (
	"fmt"
	"strconv"
	"os"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"log"
)
type Request struct {
	ClientId string
	ObjectType string
	ObjectId string
}

type RemoveClientRequest struct {
	ClientId string
}

type App int

func (t *App) RegisterClient (r *http.Request, request *Request, result *bool) error {
	fmt.Println("RegisterClient")

	insert(request.ClientId, request.ObjectType, request.ObjectId)

	fmt.Println("Client registered")
	*result = true
	return nil
}

func (t *App) UnregisterClient (r *http.Request, request *Request, result *bool) error {
	fmt.Println("UnregisterClient")

	delete(request.ClientId, request.ObjectType, request.ObjectId)

	fmt.Println("Client unregistered")
	*result = true
	return nil
}


func (t *App) RemoveClient (r *http.Request, request *RemoveClientRequest, result *bool) error {
	fmt.Println("RemoveClient")

	deleteClient(request.ClientId)

	fmt.Println("Client removed")
	*result = true
	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}


func main() {
	fmt.Println("Start")

	port, _ := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	connection = connect(
		os.Getenv("MYSQL_HOST"),
		port,
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_DATABASE"),
	)

	go runWorker()

	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
	app := new(App)
	s.RegisterService(app, "")
	r := mux.NewRouter()
	r.Handle("/rpc", s)

	appPort := os.Getenv("APP_PORT")
	fmt.Println("Server has started on port " + appPort)
	http.ListenAndServe(":" + appPort, r)
}