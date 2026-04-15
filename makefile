# 1. Detección del Sistema Operativo
ifeq ($(OS),Windows_NT)
    # Configuración para Windows
    SHELL := powershell.exe
    # En PowerShell el separador es ';'
    SEP := ;
    # Comando para abrir terminal nueva en Windows
    EXT_TERM := start powershell -NoExit -Command
else
    # Configuración para Linux/macOS
    SHELL := /bin/sh
    # En Bash/Sh el separador es '&&'
    SEP := &&
    # Comando para abrir terminal nueva (esto varía según la distro, ej. gnome-terminal)
    EXT_TERM := x-terminal-emulator -e
endif

.PHONY: go-run npm-dev dev

# 2. Comandos usando las variables detectadas

# Ejecuta el Backend en la terminal actual
go-run:
	cd backend $(SEP) go run cmd/api/main.go

# Ejecuta el Frontend en la terminal actual
npm-dev:
	cd frontend $(SEP) npm run dev

# Ejecuta ambos en terminales separadas
# Nota: En Linux esto depende de tener un emulador de terminal instalado
dev:
ifeq ($(OS),Windows_NT)
	powershell -Command "Start-Process powershell -ArgumentList '-NoExit', '-Command', 'cd backend; go run cmd/api/main.go'"
	powershell -Command "Start-Process powershell -ArgumentList '-NoExit', '-Command', 'cd frontend; npm run dev'"
else
	@echo "Lanzando terminales en Linux..."
	x-terminal-emulator -e "sh -c 'cd backend && go run cmd/api/main.go; exec sh'" &
	x-terminal-emulator -e "sh -c 'cd frontend && npm run dev; exec sh'" &
endif