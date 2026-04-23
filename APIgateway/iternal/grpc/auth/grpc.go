package grpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	auth1 "github.com/MrKrik/protos/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api auth1.AuthClient
}

var ErrUserExists = errors.New("user already exists")

func New(
	addr string,
	timeout time.Duration,
) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
	}
	return &Client{
		api: auth1.NewAuthClient(conn),
	}, nil
}

func (c *Client) Login(login string, password string) (token string, err error) {
	res, err := c.api.Login(context.Background(), &auth1.LoginRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		return "", nil
	}
	return res.GetToken(), nil
}

func (c *Client) Register(login string, password string) (err error) {
	_, err = c.api.Register(context.Background(), &auth1.RegisterRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		return nil
	}
	return nil
}
