package main

import (
	"context"
	"fmt"
	"log"
	"net"

	proto "Grpc/proto"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"google.golang.org/grpc"
)

type Todo struct {
	ID    int `gorm:"primaryKey;autoIncrement"`
	Title string
}

type server struct {
	proto.UnimplementedExampleServer
	db *gorm.DB
}

func main() {
	db, err := gorm.Open(sqlite.Open("todos.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&Todo{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}

	srv := grpc.NewServer()
	proto.RegisterExampleServer(srv, &server{db: db})

	log.Println("Server is running on port 9000")
	if err := srv.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *server) Addtodo(ctx context.Context, req *proto.AddTodoRequest) (*proto.AddTodoResponse, error) {
	todo := Todo{
		Title: req.Todo.Title,
	}

	if err := s.db.Create(&todo).Error; err != nil {
		return nil, fmt.Errorf("failed to create todo: %v", err)
	}

	return &proto.AddTodoResponse{Todo: &proto.Todo{
		Id:    int32(todo.ID),
		Title: todo.Title,
	}}, nil
}

func (s *server) Gettodo(ctx context.Context, req *proto.GetTodoRequest) (*proto.GetTodoResponse, error) {
	var todos []Todo

	if err := s.db.Find(&todos).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve todos: %v", err)
	}

	var protoTodos []*proto.Todo
	for _, todo := range todos {
		protoTodos = append(protoTodos, &proto.Todo{
			Id:    int32(todo.ID),
			Title: todo.Title,
		})
	}

	return &proto.GetTodoResponse{Todo: protoTodos}, nil
}

func (s *server) ServerReply(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	fmt.Println("Received from client:", req.Somestring)
	return &proto.HelloResponse{Reply: "Reply from server: " + req.Somestring}, nil
}
