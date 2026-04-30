package validategrpc

import (
	"context"

	val1 "github.com/MrKrik/protos/gen/go/validate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Validate interface {
	ValidateToken(ctx context.Context, token string) (ok bool, err error)
}

type serverAPI struct {
	val1.UnimplementedValidateServiceServer
	validate Validate
}

func Register(gRPC *grpc.Server, validate Validate) {
	val1.RegisterValidateServiceServer(gRPC, &serverAPI{validate: validate})
}

func (s *serverAPI) ValidateToken(ctx context.Context, req *val1.ValidateTokenRequest) (*val1.ValidateTokenResponse, error) {
	ok, err := s.validate.ValidateToken(ctx, req.Token)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &val1.ValidateTokenResponse{Success: true}, nil
}
