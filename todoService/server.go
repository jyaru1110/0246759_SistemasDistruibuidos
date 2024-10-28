package todoService

import (
	"context"
	"fmt"

	api "server/api/v1"

	todo "server/todo"
)

var _ api.TodoServiceServer = (*GrpcServer)(nil)

type GrpcServer struct {
	api.UnimplementedTodoServiceServer
	Todo *todo.TodoController
}

func newgrpcServer(todoController *todo.TodoController) (srv *GrpcServer, err error) {
	srv = &GrpcServer{
		Todo: todoController,
	}
	return srv, nil
}

func (s *GrpcServer) ProduceTodo(ctx context.Context, req *api.ProduceTodoRequest) (*api.ProduceTodoResponse, error) {
	insertRes, err := s.Todo.CreateTodo(ctx, req.Todo)
	if err != nil {
		return nil, err
	}
	fmt.Println(insertRes.InsertedID)
	id_res := insertRes.InsertedID
	id_res_string := fmt.Sprintf("%v", id_res)
	return &api.ProduceTodoResponse{Id: id_res_string}, nil
}

func (s *GrpcServer) Get(ctx context.Context, req *api.GetRequest) (*api.GetResponse, error) {
	mongo_res, err := s.Todo.GetTodo(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	decoded := &api.Todo{}

	err = mongo_res.Decode(decoded)

	if err != nil {
		return nil, err
	}

	return &api.GetResponse{Todo: decoded}, nil
}
