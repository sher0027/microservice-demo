@echo off

echo Stopping all services...

taskkill /F /IM discovery-server.exe
taskkill /F /IM api-gateway.exe
taskkill /F /IM product-service.exe
taskkill /F /IM order-service.exe
taskkill /F /IM inventory-service.exe
taskkill /F /IM notification-service.exe

echo All services stopped.
pause
