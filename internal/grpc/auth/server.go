package auth

import (
	"context"

	ssov1 "github.com/Bobahan/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Transport слой, который дергает Service слой, а именно Auth Service

const EMPTY_VALUE = 0

// интерфейс сервисного слоя `Auth`, в реализациях этих методов будет содержаться основная бизнес логика
type IAuth interface {
	Register(ctx context.Context, email string, password string) (userId int64, err error)
	Login(ctx context.Context, email string, password string, appId int) (token string, err error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type serverApi struct {
	ssov1.UnimplementedAuthServer       // дает возможность временно не реализовывать все методы
	auth                          IAuth // подключаем интерфейс Auth сервисного слоя
}

func PlugHandlers(gRPC *grpc.Server, auth IAuth) {
	ssov1.RegisterAuthServer(gRPC, &serverApi{auth: auth}) // передаем Auth
}

func (s *serverApi) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	userId, err := s.auth.Register(ctx, req.GetEmail(), req.GetPassword())

	if err != nil {
		return nil, status.Error(codes.Internal, "inernal error")
	}

	return &ssov1.RegisterResponse{
		UserId: userId,
	}, nil
}

func (s *serverApi) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == EMPTY_VALUE {
		return nil, status.Error(codes.InvalidArgument, "app_id is required")
	}

	// обращаемся к сервисному слою Auth.Login()
	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))

	if err != nil {
		return nil, status.Error(codes.Internal, "inernal error")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverApi) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if req.GetUserId() == EMPTY_VALUE {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	// обращаемся к сервисному слою Auth.IsAdmin()
	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())

	if err != nil {
		return nil, status.Error(codes.Internal, "inernal error")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}
