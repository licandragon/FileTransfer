# 1. Detección del Sistema Operativo
ifeq ($(OS),Windows_NT)
	# Configuración para Windows
	SHELL := powershell.exe
	SEP := ;
    CURL_CMD := curl.exe
else
	# Configuración para Linux/macOS
	SHELL := /bin/sh
	SEP := &&
    CURL_CMD := curl

endif

# Declarar todos los targets como PHONY para evitar conflictos con archivos locales
.PHONY: go-run npm-dev dev docker-up docker-down migrate-up migrate-down test-upload

# Variable de conexión (Centralizada)
DB_URL := "postgres://postgres:postgres@localhost:5432/filetransfer?sslmode=disable"

FILE_PATH := ./test/image (4).png

# 2. Comandos

go-run:
	cd backend $(SEP) go run cmd/api/main.go

npm-dev:
	cd frontend $(SEP) npm run dev

docker-up:
	docker compose up -d

docker-down:
	docker compose down -v

migrate-up:
	cd backend $(SEP) migrate -path ./migrations -database $(DB_URL) up

migrate-down:
	cd backend $(SEP) migrate -path ./migrations -database $(DB_URL) down

test-upload:
	$(CURL_CMD) "http://127.0.0.1:3000/upload" -X POST -F "sender_email=tu_correo@ejemplo.com" -F "subject_email=Envio de prueba" -F "message_email=Hola desde PowerShell" -F "recipients=amigo1@ejemplo.com" -F "recipients=amigo2@ejemplo.com" -F "files=@$(FILE_PATH)"
# Ejecuta ambos en terminales separadas
dev:
ifeq ($(OS),Windows_NT)
	powershell -Command "Start-Process powershell -ArgumentList '-NoExit', '-Command', 'cd backend; go run cmd/api/main.go'"
	powershell -Command "Start-Process powershell -ArgumentList '-NoExit', '-Command', 'cd frontend; npm run dev'"
else
	@echo "Lanzando terminales en Linux..."
	x-terminal-emulator -e "sh -c 'cd backend && go run cmd/api/main.go; exec sh'" &
	x-terminal-emulator -e "sh -c 'cd frontend && npm run dev; exec sh'" &
endif