package test

import (
	"context"
	"net"
	"testing"
	"user/gen/user"
	"github.com/example/user/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024

func TestLogin(t *testing.T) {
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	user.RegisterUserServiceServer(s, &server.Server{})
	go s.Serve(lis)  // 内存网络层测试‌:ml-citation{ref="4,6" data="citationList"}

	conn, _ := grpc.DialContext(context.Background(), "bufnet",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	client := user.NewUserServiceClient(conn)
	tests := []struct {
		name     string
		req      *user.LoginRequest
		wantCode int32
		wantErr  bool
	}{
		{"Valid", &user.LoginRequest{Username: "admin", Password: "123456"}, 0, false},
		{"Invalid", &user.LoginRequest{Username: "guest", Password: "wrong"}, 401, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.Login(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Test %s failed: %v", tt.name, err)
			}
			if !tt.wantErr && resp.Code != tt.wantCode {
				t.Errorf("Code mismatch: got %d, want %d", resp.Code, tt.wantCode)
			}
		})
	}
}

