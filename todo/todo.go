package todo

import (
	"context"
	api "server/api/v1"

	"go.mongodb.org/mongo-driver/mongo"
)

type TodoController struct {
	db *mongo.Collection
}

func NewTodoController(db *mongo.Collection) *TodoController {
	return &TodoController{db: db}
}

func (t *TodoController) CreateTodo(ctx context.Context, newTodo *api.Todo) (*mongo.InsertOneResult, error) {
	result, err := t.db.InsertOne(ctx, newTodo)

	if err != nil {
		return nil, err
	}

	return result, err
}

func (t *TodoController) GetTodo(ctx context.Context, id string) (*mongo.SingleResult, error) {
	single_res := t.db.FindOne(ctx, id)

	if single_res.Err() != nil {
		return nil, single_res.Err()
	}

	return single_res, nil
}
