# ARQAP Backend

Backend de ARQAP escrito en **Go**, con entorno de desarrollo montado en **Docker Compose** y **Air** para soportar *hot reload* durante el desarrollo.

## Requisitos

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Tecnologías

- [Go](https://go.dev/) — Lenguaje principal
- [Air](https://github.com/air-verse/air) — Hot reload para Go
- [PostgreSQL](https://www.postgresql.org/) — Base de datos

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
- `api` → backend en Go con hot reload (Air).
- `db` → PostgreSQL 16.

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