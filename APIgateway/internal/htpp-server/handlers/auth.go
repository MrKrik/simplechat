package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	grpc "tcp-server/internal/grpc/auth"
)

func Login(client *grpc.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request string
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			fmt.Println("f")
		}
		login, password, _ := strings.Cut(request, " ")
		tk, _ := client.Login(login, password)
		fmt.Println(tk)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tk)
	}
}

func Register(client *grpc.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request string
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			fmt.Println("f")
		}
		login, password, _ := strings.Cut(request, " ")
		err = client.Register(login, password)
		if err != nil {
			return
		}
		w.WriteHeader(200)
	}
}
