# FileTransfer — Clon tipo WeTransfer

## 📌 Descripción

FileTransfer es una aplicación web tipo WeTransfer que permite subir archivos, generar enlaces públicos de descarga y gestionar la expiración automática de archivos.

Este proyecto implementa una arquitectura distribuida utilizando:

* Backend en Go
* Frontend en Vue.js
* Base de datos principal PostgreSQL
* Supabase como almacenamiento y respaldo

---

## 🏗️ Arquitectura

```bash
Frontend (Vue)
     │
     │ HTTP API
     ▼
Backend (Go)
     │
     ├── PostgreSQL (Metadata principal)
     │
     └── Supabase Storage (Archivos)
```

---

## 🛠️ Tecnologías a utilizar

## Backend

* Go
* Fiber
* PostgreSQL
* Supabase Storage

## Frontend

* Vue 3
* TailwindCSS
* Vite

## Base de datos

* PostgreSQL (principal)
* Supabase (secundaria)

---

## 🚀 Instalación

### 1. Clonar repositorio

```bash
git clone https://github.com/licandragon/FileTransfer.git

cd FileTransfer
```

---

## ⚙️ **Backend**

### Configuración

```bash
cd backend

cp .env.example .env

```

Editar variables de entorno:

```bash
PORT=3000

DATABASE_URL=postgres://user:password@localhost:5432/filetransfer

SUPABASE_URL=your_supabase_url
SUPABASE_SERVICE_KEY=your_service_key
SUPABASE_BUCKET=uploads
```

---

### Ejecutar Backend

```bash
go run ./cmd/server
```

Servidor disponible en:

```bash
http://localhost:3000
```

---

### 📡 **API Endpoints**

#### **Iniciar transferencia**

```bash
POST /api/transfer/ini
```

Se genera una nueva transferencia y retorna token valido para subida de archivo {upload_token}.

---

#### **Subir archivos**

```bash
POST /api/transfer/{upload_Token}/file
```

Sube archivo a transferencia con token valido.

---

```bash
PATCH /api/transfer/complete
```

Finaliza transferencia retornando el token de url publica (download_token).

---

#### **Descargar archivos**

```bash
GET /api/download/{download_token}
```

Descarga archivo mediante token público.

---

```bash
GET /api/download/{download_token}/files/
```

Genera el url firmado para la descarga un archivo en especifico, el enlace es valido durante 1hr.

---

#### **Información archivo**

```bash
GET /file/{token}
```

Obtiene metadata del archivo.

---

#### **Eliminar archivo (Opcional)**

```bash
DELETE /file/{token}
```

Elimina archivo manualmente.

---

## 🎨 **Frontend**

```bash
cd frontend

cp .env.example .env
```

Variables:

```env
VITE_API_URL=http://localhost:3000

VITE_SUPABASE_URL=your_supabase_url
VITE_SUPABASE_ANON_KEY=your_anon_key
```

---

### Ejecutar Frontend

```bash
npm install

npm run dev
```

Frontend disponible en:

```bash
http://localhost:5173
```

---

## 🔐 **Seguridad**

Seguridad a implementar:

* Validación MIME
* Límite de tamaño
* Tokens UUID
* Protección Path Traversal
* Validación de archivos ejecutables
* Manejo de errores

---

## ⏱️ Expiración Automática

Los archivos tienen fecha de expiración.

Un worker eliminara:

* Archivos expirados
* Metadata en PostgreSQL
* Archivos en Supabase Storage

---

## 👤 Usuarios

El sistema implementara:

* Usuarios invitados
* Links públicos
* Subida sin registro (limitado 1-3 dias)

---

## 🧪 Desarrollo

Ejecutar backend:

```bash
go run ./cmd/server
```

Ejecutar frontend:

```bash
npm run dev
```

## 🤖 **Uso de IA**

Durante el  desarrollo se hizo uso de la IA para lo siguiente:

* Consultar documentacion
* Correcciones de bugs y mejorar la experincia del usuario UI/UX
* Corregir flujo de trajo para que fuera mas estructurado y escalable
