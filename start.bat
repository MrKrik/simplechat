@echo off
setlocal

echo Starting services...

:: Запускаем сервисы и присваиваем окнам уникальный префикс "SRV_"
start "SRV_API_Gateway" cmd /k "cd APIgateway && go run ./cmd/main.go --config=./config/local.yaml"
start "SRV_Auth_Service" cmd /k "cd auth && go run ./cmd/auth/main.go --config=./config/local.yaml"
start "SRV_WS_Gateway" cmd /k "cd wsgateway && go run ./cmd/main.go"

echo.
echo All services are running in separate windows.
echo To STOP all services and close their windows, press any key IN THIS WINDOW.
echo.

pause

echo Stopping all services...

:: Закрываем все окна, заголовок которых начинается на SRV_
taskkill /FI "WINDOWTITLE eq SRV_*" /F /T

echo All services have been stopped.
pause