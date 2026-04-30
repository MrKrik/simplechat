package grpc

import (
	"context"
	"fmt"
	"time"

	validate "github.com/MrKrik/protos/gen/go/validate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func (c *Client) ValidateToken(token string) (bool, error) {
	res, err := c.api.ValidateToken(context.Background(), &validate.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		return false, err
	}
	if res.Success {
		return true, nil
	}
	return false, nil
}
