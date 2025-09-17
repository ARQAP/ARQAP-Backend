# ARQAP Backend

Backend de ARQAP escrito en **Go**, con entorno de desarrollo montado en **Docker Compose** y **Air** para soportar _hot reload_ durante el desarrollo.

## Requisitos

-   [Docker](https://www.docker.com/get-started)
-   [Docker Compose](https://docs.docker.com/compose/install/)

## Tecnologías

-   [Go](https://go.dev/) — Lenguaje principal
-   [Air](https://github.com/air-verse/air) — Hot reload para Go
-   [PostgreSQL](https://www.postgresql.org/) — Base de datos

## Estructura del proyecto

```
ARQAP-Backend/
├── .air.toml           # Configuración de Air
├── .gitignore          # Configuración de Gitignore
├── docker-compose.yml  # Orquestación de servicios
├── Dockerfile          # Imagen del servicio API
├── go.mod              # Dependencias de Go
└── main.go             # Código fuente en Go

```

## Levantar el entorno

Compilar y levantar servicios:

```bash
docker compose up --build
```

Esto arranca:

-   `api` → backend en Go con hot reload (Air).
-   `db` → PostgreSQL 16.

El backend queda escuchando en:

```
http://localhost:8080
```

## Hot Reload

Gracias a **Air**, cada vez que se edite un archivo `.go` dentro de `app/`, el servicio se recompila automáticamente sin necesidad de reiniciar manualmente el contenedor.

Logs esperados en la consola de `api`:

```
watching .
building...
running...
Servidor escuchando en :8080
```

## Detener servicios

```bash
docker compose down -v
```

(`-v` también elimina volúmenes, útil si querés borrar la base de datos y empezar de cero).

---

## PGAdmin

Dashboard de la BD disponible en:

[http://localhost:5050](http://localhost:5050)

Credenciales por defecto:

-   Usuario: `admin`
-   Contraseña: `admin`
-   Contraseña maestra: `pass`

## API

### Autenticación

-   `POST /register` → Registrar nuevo usuario (recibe JSON con `username` y `password`).
-   `POST /login` → Iniciar sesión (recibe JSON con `username` y `password`, devuelve JWT).

### Credenciales de usuario creadas en la inicialización

-   Usuario: `arqap`
-   Contraseña: `arqap`

### Como realizar una peticion a una ruta protegida

Para realizar una petición protegida, es necesario incluir el token JWT que responde la API al iniciar sesión en el HEADER llamado `Authorization` dentro de la solicitud HTTP.

El formato del encabezado debe ser exactamente el siguiente:

`Authorization: Bearer <JWT_TOKEN>`

## Entidades y endpoints disponibles:

### 👤 Usuarios

- **CRUD**

---

### 🏺 Piezas Arqueológicas

#### 📋 Información básica de pieza arqueológica

-   **CRUD** 

#### 📷 Imágen adjunta de pieza arqueológica

**Subir imagen:**

-   **Método:** `POST`
-   **Formato:** Multipart/Form-Data
-   **URL:** `{host}/artefacts/:id/picture/`
-   **Key:** `picture`
-   **Content-Type:** `file` / `Auto`
-   **Value:** Seleccionar archivo

**Servir imagen:**

-   **Método:** `GET`
-   **URL:** `{host}/artefacts/:id/picture/`

#### 📄 Imágen adjunta de ficha histórica de pieza arqueológica

**Subir documento:**

-   **Método:** `POST`
-   **Formato:** Multipart/Form-Data
-   **URL:** `{host}/artefacts/:id/historical-record/`
-   **Key:** `document`
-   **Content-Type:** `file` / `Auto`
-   **Value:** Seleccionar archivo

**Servir documento:**

-   **Método:** `GET`
-   **URL:** `{host}/artefacts/:id/historical-record/`

---

### 👨‍🔬 Arqueólogos

-   **CRUD** 

---

### 📚 Colecciones

-   **CRUD** 

---

### 💬 Menciones

-   **CRUD** 

---

### 🌍 Ubicaciones Geográficas

#### 🏳️ Países

-   **CRUD** 

#### 🗺️ Regiones

-   **CRUD** 

#### 🏛️ Sitios Arqueológicos

-   **CRUD** 

---

### 🏷️ Clasificadores de piezas arqueológicas

#### 📊 Clasificadores INPL

-   **CRUD** 

#### 🔖 Clasificadores Interno

-   **CRUD** 

---

### 🤝 Préstamos y solicitantes

#### 📋 Préstamos

-   **CRUD** 

#### 👤 Solicitante de Préstamo

-   **CRUD** 

---

### 📍 Ubicaciones Físicas

#### 📚 Estanterías

-   **CRUD** 

#### 🏢 Ubicación Física

-   **CRUD** 
