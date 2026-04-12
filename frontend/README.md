# Vue 3 + Supabase 🚀 (Login + Registro)

Aplicación web construida con Vue 3 que implementa autenticación de usuarios utilizando Supabase. Incluye registro, inicio de sesión, protección de rutas y un dashboard básico.

---

## 🧰 Tecnologías

* Vue 3 (Composition API)
* Vue Router
* Supabase (Auth + Database)
* Tailwind CSS

---

## ⚙️ Instalación

Clona el repositorio:

```bash
git clone https://github.com/licandragon/login_supabase.git
cd login_supabase
```

Instala dependencias:

```bash
npm install
```

---

---

## 🔐 Configuración de varibales de entorno

1. Crea un proyecto en Supabase
2. Copia tus credenciales (URL y anon key)
3. Copia el archivo de ejemplo:

```bash
cp .env.example .env
```

4. Edita `.env` con tus credenciales de **[Supabase](https://supabase.com)**:

```env
VITE_SUPABASE_URL=tu_url_real
VITE_SUPABASE_ANON_KEY=tu_anon_key_real
```

---

## ▶️ Ejecutar el proyecto

```bash
npm run dev
```

---

## 📁 Estructura básica

```
src/
├── composables/
│   └── useAuth.js
 views/
│   ├── LoginView.vue
│   ├── RegisterView.vue
│   └── DashboardView.vue
├── router/
│   └── index.js
├── lib/
│   └── supabaseClient.js
```

---

## 🔑 Funcionalidades

* Registro de usuario
* Inicio de sesión
* Cierre de sesión
* Protección de rutas
* Persistencia de sesión
* Dashboard básico

---

## 🔒 Seguridad

* Uso de autenticación de Supabase
* Protección de rutas con Vue Router
* Manejo de sesión reactivo

---

## 📌 Notas

* Asegúrate de desactivar la confirmación por email en Supabase si estás en desarrollo

---

## 📌 Notas
* Asegúrate de desactivar la confirmación por email en Supabase si estás en desarrollo
---

## ⚡ LiveDemo (Vercel)

[login-supabase-brown.vercel.app](login-supabase-brown.vercel.app)