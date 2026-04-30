package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	grpc "tcp-server/internal/grpc/auth"
)

func Login(client *grpc.Client, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request string
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Error("Failed decode login request")
		}
		login, password, _ := strings.Cut(request, " ")
		tk, errMSG := client.Login(login, password, log)
		if errMSG != "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(errMSG)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tk)
	}
}

func Register(client *grpc.Client, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request string
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			fmt.Println("Error decode reigster json")
		}
		login, password, _ := strings.Cut(request, " ")
		errMSG := client.Register(login, password, log)
		if errMSG != "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(errMSG)
			return
		}
		w.WriteHeader(200)
	}
}
