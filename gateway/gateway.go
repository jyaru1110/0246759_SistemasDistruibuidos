package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	api "server/api/v1"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr_todo = flag.String("addr_todo", "todoservice:8081", "the address to connect to")
var addr_log = flag.String("addr_log", "logservice:8080", "the address to connect to")

func main() {
	todo_conn, err := grpc.NewClient(*addr_todo, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("could not connect to todo grpcserver: %v", err)
		return
	}

	log_conn, err := grpc.NewClient(*addr_log, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("could not connect to log grpcserver: %v", err)
		return
	}

	defer log_conn.Close()
	defer todo_conn.Close()

	log_grpc_client := api.NewLogClient(log_conn)
	todo_grpc_client := api.NewTodoServiceClient(todo_conn)

	ctx, cancel := context.WithTimeout(context.Background(), 100000*time.Millisecond)
	defer cancel()

	rest_server := mux.NewRouter()

	rest_server.HandleFunc("/todo", func(w http.ResponseWriter, r *http.Request) {
		var todo api.Todo
		json.NewDecoder(r.Body).Decode(&todo)
		_, error := todo_grpc_client.ProduceTodo(ctx, &api.ProduceTodoRequest{Todo: &todo})

		if error != nil {
			fmt.Println("Error: ", error)
			log_res, _ := log_grpc_client.Produce(ctx, &api.ProduceRequest{Record: &api.Record{Value: []byte(error.Error())}})
			w.WriteHeader(http.StatusInternalServerError)
			log_offset_str := string(log_res.Offset)
			w.Write([]byte(log_offset_str))
			return
		}

		w.WriteHeader(http.StatusOK)
	}).Methods("POST")

	rest_server.HandleFunc("/todo/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		log_response, error := log_grpc_client.Produce(ctx, &api.ProduceRequest{Record: &api.Record{Value: []byte("Get todo by id: " + id)}})

		fmt.Println("Log offset: ", log_response.Offset)

		if error != nil {
			fmt.Print("Failed to log: ", error)
		}

		todo, error := todo_grpc_client.Get(ctx, &api.GetRequest{Id: id})

		if error != nil {
			log_grpc_client.Produce(ctx, &api.ProduceRequest{Record: &api.Record{Value: []byte(error.Error())}})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(todo)
	}).Methods("GET")

	rest_server.HandleFunc("/log/{offset}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		offset := vars["offset"]

		uint_offset, err := strconv.ParseUint(offset, 10, 64)

		if err != nil {
			fmt.Println("Error: ", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log, error := log_grpc_client.Consume(ctx, &api.ConsumeRequest{Offset: uint_offset})

		if error != nil {
			fmt.Println("Error: ", error)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(log)
	}).Methods("GET")

	fmt.Println("Server started at :8082")
	http.ListenAndServe(":8082", rest_server)
}
