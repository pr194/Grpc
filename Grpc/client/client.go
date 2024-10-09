package main

import (
	"context"
	"fmt"
	"log"

	proto "Grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := proto.NewExampleClient(conn)

	// Call ServerReply
	callServerReply(client)

	// Call Addtodo
	callAddTodo(client)

	// Call Gettodo
	callGetTodo(client)
}

func callServerReply(client proto.ExampleClient) {
	req := &proto.HelloRequest{Somestring: "hi there it's me prince"}
	response, err := client.ServerReply(context.TODO(), req)
	if err != nil {
		log.Fatalf("Error calling ServerReply: %v", err)
	}

	fmt.Printf("Server response: %s\n", response.Reply)
}

func callAddTodo(client proto.ExampleClient) {
	req := &proto.AddTodoRequest{
		Todo: &proto.Todo{
			Title: "New Todo Item",
		},
	}
	response, err := client.Addtodo(context.TODO(), req)
	if err != nil {
		log.Fatalf("Error calling Addtodo: %v", err)
	}

	fmt.Printf("Added Todo: ID = %d, Title = %s\n", response.Todo.Id, response.Todo.Title)
}

func callGetTodo(client proto.ExampleClient) {
	req := &proto.GetTodoRequest{}
	response, err := client.Gettodo(context.TODO(), req)
	if err != nil {
		log.Fatalf("Error calling Gettodo: %v", err)
	}

	fmt.Println("Todos received from server:")
	for _, todo := range response.Todo {
		fmt.Printf("ID = %d, Title = %s\n", todo.Id, todo.Title)
	}
}
