package authgrpc

import (
	"context"

	auth1 "github.com/MrKrik/protos/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(
		ctx context.Context,
		login string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		login string,
		password string,
	) (user int64, err error)
}

type serverAPI struct {
	auth1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	auth1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

// brutforce
func (s *serverAPI) Login(
	ctx context.Context,
	req *auth1.LoginRequest,
) (*auth1.LoginResponse, error) {
	// TODO: validation v10
	if req.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "app_id is required")
	}

	token, err := s.auth.Login(ctx, req.GetLogin(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &auth1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *auth1.RegisterRequest,
) (*auth1.RegisterResponse, error) {
	// TODO: validation v10
	if req.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	uid, err := s.auth.RegisterNewUser(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		// TODO: ...
		// if errors.Is(err, storage.ErrUserExists) {
		// 	return nil, status.Error(codes.AlreadyExists, "user already exists")
		// }

		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &auth1.RegisterResponse{UserId: uid}, nil
}
