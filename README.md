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

---

## 🔗 Integración con Google Drive API

El sistema puede descargar archivos desde Google Drive automáticamente durante la importación de artefactos desde Excel. Para habilitar esta funcionalidad, es necesario configurar las credenciales de Google Drive API.

### Configuración de Google Drive API

#### 1. Crear un Service Account en Google Cloud

1. Ve a [Google Cloud Console](https://console.cloud.google.com/)
2. Crea un nuevo proyecto o selecciona uno existente
3. Habilita la **Google Drive API** para el proyecto
4. Ve a **IAM & Admin** > **Service Accounts**
5. Crea un nuevo Service Account
6. Descarga el archivo JSON de credenciales

#### 2. Compartir archivos/carpetas con el Service Account

Para que el Service Account pueda acceder a los archivos de Google Drive:

1. Abre el archivo JSON de credenciales
2. Copia el **email del Service Account** (campo `client_email`)
3. En Google Drive, comparte los archivos/carpetas con ese email (dar permisos de "Lector")

#### 3. Configurar las credenciales en el backend

Tienes dos opciones:

**Opción A: Usar archivo de credenciales (recomendado para desarrollo local)**

```bash
export GOOGLE_DRIVE_CREDENTIALS_PATH="/ruta/al/archivo/credentials.json"
```

**Opción B: Usar JSON como variable de entorno (recomendado para producción/Docker)**

```bash
export GOOGLE_DRIVE_CREDENTIALS_JSON='{"type":"service_account","project_id":"...","private_key_id":"...","private_key":"...","client_email":"...","client_id":"...","auth_uri":"...","token_uri":"...","auth_provider_x509_cert_url":"...","client_x509_cert_url":"..."}'
```

#### 4. Agregar al docker-compose.yml (si usas Docker)

```yaml
services:
  app:
    environment:
      - GOOGLE_DRIVE_CREDENTIALS_PATH=/app/credentials.json
    volumes:
      - ./credentials.json:/app/credentials.json:ro
```

### Notas importantes

- El Service Account necesita permisos de **lectura** en los archivos/carpetas de Google Drive
- Los archivos deben estar compartidos con el email del Service Account
- Si no se configuran las credenciales, el sistema intentará usar descarga HTTP directa (puede fallar para archivos grandes)
- El sistema detecta automáticamente URLs de Google Drive y usa la API cuando está disponible 
