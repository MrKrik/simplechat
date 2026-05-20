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
	api     validate.ValidateServiceClient
	timeout time.Duration
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
		api:     validate.NewValidateServiceClient(conn),
		timeout: timeout,
	}, nil
}

func (c *Client) ValidateToken(token string) (ok bool, errMSG string) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	res, err := c.api.ValidateToken(ctx, &validate.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		return false, err.Error()
	}
	st, ok := status.FromError(err)
	if ok {
		switch st.Code() {
		case codes.InvalidArgument:
			return false, st.Message()
		}
	}
	if res.Success {
		return true, ""
	} else {
		return false, ""
	}
}
