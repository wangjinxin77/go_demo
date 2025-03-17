package main

import (
	"context"
	"log"
	"time"
	"github.com/example/user/gen/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, _ := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	client := user.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := client.Login(ctx, &user.LoginRequest{Username: "admin", Password: "123456"})
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	log.Printf("Token: %s", resp.Token)
}

