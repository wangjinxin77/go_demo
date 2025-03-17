package main

import (
    "context"
    "log"
    "net"

    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    pb "github.com/example/user/gen/user"  // 需替换为实际模块路径‌:ml-citation{ref="1,3" data="citationList"}
)

// 实现UserService接口
type userServer struct {
    pb.UnimplementedUserServiceServer
}

// Login方法实现（基于最新proto定义）
func (s *userServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
    // 验证用户名和密码‌:ml-citation{ref="1,3" data="citationList"}
    if req.GetUsername() != "admin" || req.GetPassword() != "123456" {
        st := status.New(codes.Unauthenticated, "invalid credentials")
        return nil, st.Err()  // 返回gRPC标准错误‌:ml-citation{ref="2,5" data="citationList"}
    }

    // 成功时返回code=200和token‌:ml-citation{ref="1,3" data="citationList"}
    return &pb.LoginResponse{
        Code:  200,
        Token: "generated_jwt_token_here",
    }, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("监听失败: %v", err)
    }

    s := grpc.NewServer()
    pb.RegisterUserServiceServer(s, &userServer{})

    log.Println("Server running on :50051")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("服务启动失败: %v", err)
    }
}

