@echo off

echo Starting API Gateway...
start "API Gateway" cmd /k "cd APIgateway && go run ./cmd/main.go --config=./config/local.yaml"

echo Starting Auth Service...
start "Auth Service" cmd /k "cd auth && go run ./cmd/auth/main.go --config=./config/local.yaml"

echo Starting WS Gateway...
start "WS Gateway" cmd /k "cd wsgateway && go run ./cmd/main.go"

echo All services are starting in separate windows.
pause