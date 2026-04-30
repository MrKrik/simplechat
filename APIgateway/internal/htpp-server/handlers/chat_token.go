package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	grpc "tcp-server/internal/grpc/auth"
)

func GetChatToken(client *grpc.Client, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request string
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Error("Failed decode login request")
		}
		token, errMSG := client.GetChatToken(request, log)
		if errMSG != "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(errMSG)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(token)
	}
}
