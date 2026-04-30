package grpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	auth1 "github.com/MrKrik/protos/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
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

func (c *Client) Login(login string, password string, log *slog.Logger) (token string, errMessage string) {
	res, err := c.api.Login(context.Background(), &auth1.LoginRequest{
		Login:    login,
		Password: password,
		AppId:    1,
	})
	st, ok := status.FromError(err)
	if ok {
		switch st.Code() {
		case codes.InvalidArgument:
			log.Error(err.Error())
			return "", st.Message()
		case codes.Internal:
			log.Error(err.Error())
			return "", st.Message()
		}
	}
	return res.GetToken(), ""
}

func (c *Client) Register(login string, password string, log *slog.Logger) (errMessage string) {
	_, err := c.api.Register(context.Background(), &auth1.RegisterRequest{
		Login:    login,
		Password: password,
	})
	st, ok := status.FromError(err)
	if ok {
		switch st.Code() {

		case codes.InvalidArgument:
			log.Error(err.Error())
			return st.Message()

		case codes.Internal:
			log.Error(err.Error())
			return st.Message()

		case codes.AlreadyExists:
			log.Error(err.Error())
			return st.Message()
		}
	}
	return ""
}

func (c *Client) GetChatToken(authToken string, log *slog.Logger) (errMessage string, token string) {
	res, err := c.api.GetChatToken(context.Background(), &auth1.GetChatTokenRequest{
		AuthToken: authToken,
	})
	st, ok := status.FromError(err)
	if ok {
		switch st.Code() {
		case codes.InvalidArgument:
			log.Error(err.Error())
			return st.Message(), ""
		}
	}
	return res.ChatToken, ""
}
