@echo off

start "Discovery Server" cmd /c "cd discovery-server && discovery-server.exe > discovery-server.log 2>&1"
start "API Gateway" cmd /c "cd api-gateway && api-gateway.exe > api-gateway.log 2>&1"
start "Product Service" cmd /c "cd product-service && product-service.exe > product-service.log 2>&1"
start "Order Service" cmd /c "cd order-service && order-service.exe > order-service.log 2>&1"
start "Inventory Service" cmd /c "cd inventory-service && inventory-service.exe > inventory-service.log 2>&1"

rem Start multiple instances of notification-service
start "Notification Service 1" cmd /c "cd notification-service && notification-service.exe > notification-service-1.log 2>&1"
start "Notification Service 2" cmd /c "cd notification-service && notification-service.exe > notification-service-2.log 2>&1"
start "Notification Service 3" cmd /c "cd notification-service && notification-service.exe > notification-service-3.log 2>&1"

echo All services started.
pause
