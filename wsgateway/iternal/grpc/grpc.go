package grpc

import (
	"context"
	"fmt"
	"time"

	validate "github.com/MrKrik/protos/gen/go/validate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Client struct {
	api validate.ValidateServiceClient
}

func New(
	addr string,
	timeout time.Duration,
) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
	}
	return &Client{
		api: validate.NewValidateServiceClient(conn),
	}, nil
}

func (c *Client) ValidateToken(token string) (ok bool, errMSG string) {
	res, err := c.api.ValidateToken(context.Background(), &validate.ValidateTokenRequest{
		Token: token,
	})
	st, ok := status.FromError(err)
	if ok {
		switch st.Code() {
		case codes.InvalidArgument:
			return false, st.Message()
		}
	}
	if res.Success {
		return true, ""
	}
	return false, ""
}
