# Gator

Gator es un **agregador de feeds RSS** de línea de comandos escrito en Go. Permite
registrar usuarios, dar de alta feeds, seguirlos, y recolectar/alistar posts de
forma periódica, todo desde la terminal.

## Características

- Registro y login de usuarios (persistidos en Postgres).
- Alta y listado de feeds RSS/Atom.
- Seguir / dejar de seguir feeds por usuario.
- Recolección periódica de posts (`agg`) con rotación por `last_fetched_at`.
- Comando `browse` para listar los posts de tus feeds seguidos.
- `help` integrado con la lista de comandos.

## Requisitos

- Go 1.21+ (el `go.mod` usa 1.26.5).
- PostgreSQL en ejecución.
- [`sqlc`](https://sqlc.dev) para regenerar el código de base de datos (opcional en dev).
- [`goose`](https://github.com/pressly/goose) para aplicar migraciones.

## Instalación

```bash
git clone https://github.com/NicolasFerreras/Gator.git
cd Gator
go build -o bin/gator ./cmd/gator
# o ejecutar directamente
go run ./cmd/gator <comando>
```

## Configuración

Gator lee `~/.gatorconfig.json`. Este archivo **está ignorado en git**; cada
usuario debe crear el suyo. Formato:

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

- `db_url`: connection string de PostgreSQL.
- `current_user_name`: se completa solo tras `login` / `register`.

## Base de datos

1. Levantar PostgreSQL y crear la base (ej. `gator`).
2. Aplicar las migraciones (formato goose en `internal/database/schema`):

```bash
goose -dir internal/database/schema postgres "$DB_URL" up
```

Tablas: `users`, `feeds`, `posts`, `feed_follow`.

> Para regenerar el código Go tras editar queries: `sqlc generate`

## Modelo de datos (DER)

```mermaid
erDiagram
    USERS {
        uuid id PK
        timestamp created_at
        timestamp updated_at
        varchar username UK
    }
    FEEDS {
        uuid id PK
        timestamp created_at
        timestamp updated_at
        varchar name UK
        varchar url UK
        uuid user_id FK
        timestamp last_fetched_at
    }
    POSTS {
        uuid id PK
        timestamp created_at
        timestamp updated_at
        varchar title
        varchar url UK
        text description
        timestamp published_at
        uuid feed_id FK
    }
    FEED_FOLLOW {
        uuid id PK
        timestamp created_at
        timestamp updated_at
        uuid user_id FK
        uuid feed_id FK
    }

    USERS ||--o{ FEEDS : "crea (user_id)"
    USERS ||--o{ FEED_FOLLOW : "sigue (user_id)"
    FEEDS ||--o{ FEED_FOLLOW : "es seguido (feed_id)"
    FEEDS ||--o{ POSTS : "publica (feed_id)"
```

## Arquitectura y Patrones de Diseño

**Tipo de arquitectura:** Aplicación **CLI monolítica de arquitectura por capas
(layered)**, de estilo pragmático. No utiliza un framework de CLI (p.ej. Cobra);
el enrutado de comandos es propio y liviano.

**Capas:**
- **Entrada / CLI (`internal/cli`)** — parsea argumentos, registra y despacha
  comandos, aplica middleware de autenticación y contiene la lógica de handlers.
- **Acceso a datos (`internal/database`)** — código generado por `sqlc` desde
  `schema/` y `queries/`; actúa como repositorio.
- **Integración externa (`internal/rss`)** — fetch y parseo RSS/Atom vía HTTP +
  `encoding/xml`.
- **Configuración (`internal/config`)** — lectura/escritura de `.gatorconfig.json`.
- **Modelos / DTOs (`internal/models`)** — estructuras de transporte (`RSSFeed`).

**Patrones de diseño:**
- **Command Pattern** — `Commands` mantiene un mapa `cmdMap` de
  `nombre → handler`; `register` da de alta y `run` despacha por nombre.
- **Middleware / Decorator (higher-order function)** — `middlewareLoggedIn`
  envuelve un handler, resuelve el usuario desde la BD y lo inyecta, añadiendo
  el requisito de sesión sin modificar la lógica del handler.
- **Repository (implícito)** — `database.Queries` concentra toda la persistencia
  (sqlc), aislando el SQL de los handlers.
- **Dependency Injection vía `State`** — el struct `State{Config, Db}` se pasa a
  todos los handlers como contenedor de dependencias (DI manual).
- **Centralized Error Handling** — `errors_handling` con factories que envuelven
  errores con contexto.
- **Code Generation** — `sqlc` genera el data layer desde SQL; el schema usa
  marcas `goose` para migraciones.

**Limitaciones:**
- Los handlers están **acoplados a la implementación concreta**
  `*database.Queries` (no a una interfaz/puerto), por lo que **no** es
  arquitectura hexagonal ni Clean Architecture estricta.
- No hay capa de servicios / use-cases separada: la lógica de negocio vive en los
  handlers.
- La suite de tests actual **mockea la BD** (no hay tests de integración contra
  Postgres real).

**Cómo mejorar / refactorizar estas limitaciones:**

1. **Desacoplar la capa de datos (hacia arquitectura hexagonal / Clean):**
   - Definir una interfaz `Querier` (o puertos por dominio, p.ej.
     `UserRepository`, `FeedRepository`) en `internal/cli` o en un paquete
     `ports`, con los métodos que usan los handlers.
   - Cambiar `State.Db` de `*database.Queries` a esa interfaz. Como
     `*database.Queries` ya implementa esos métodos, no hace falta modificar el
     código generado por sqlc.
   - Beneficio: los handlers pasan a ser testeables con fakes en memoria y el
     proyecto evoluciona hacia puertos y adaptadores.

2. **Introducir una capa de servicios / use-cases:**
   - Extraer la lógica de negocio de los handlers hacia servicios
     (p.ej. `UserService.Register`, `FeedService.Follow`) que reciban los puertos
     y devuelvan errores de dominio.
   - Los handlers quedarían como controladores delgados: parsean argumentos,
     invocan al servicio y formatean la salida. Esto centraliza reglas de negocio
     y facilita reutilización y testing.

3. **Agregar tests de integración contra Postgres real:**
   - Levantar una base efímera (Testcontainers o `docker run` de Postgres) en un
     `TestMain` o helper, aplicar migraciones con goose y ejecutar los queries
     sqlc contra ella (sin `sqlmock`).
   - Mantener los tests unitarios con `sqlmock` para velocidad, y sumar los de
     integración para validar el SQL y el esquema reales.
   - Separar con sufijo `_integration_test.go` o build tags (p.ej. `//go:build
     integration`) para no ralentizar la suite diaria.

## Comandos

| Comando | Auth | Descripción |
|---------|------|-------------|
| `login <usuario>` | – | Inicia sesión como usuario existente. |
| `register <usuario>` | – | Crea un usuario y inicia sesión. |
| `reset` | – | Elimina todos los usuarios. |
| `users` | – | Lista usuarios (marca el actual con `*`). |
| `agg <duración>` | – | Recolecta feeds cada `<duración>` (ej. `1h`, `30m`). Loop continuo. |
| `addfeed <nombre> <url>` | ✓ | Agrega un feed y lo sigue. |
| `feeds` | – | Lista todos los feeds con su dueño. |
| `follow <url>` | ✓ | Sigue un feed existente. |
| `following` | ✓ | Lista los feeds que sigues. |
| `unfollow <url>` | ✓ | Deja de seguir un feed. |
| `browse [límite]` | – | Lista tus posts (límite por defecto: 2). |
| `help` | – | Muestra la ayuda. |

> ✓ = requiere haber hecho `login`/`register` primero.

```bash
gator register nico
gator addfeed "TechCrunch" https://techcrunch.com/feed/
gator follow https://www.wagslane.dev/index.xml
gator agg 30m
gator browse 5
```

## Desarrollo y Tests

Suite de tests unitarios **aislados** (sin Postgres ni red externa):
- `internal/rss/fetch_test.go` — `FetchFeed` con `httptest`.
- `internal/config/config_test.go` — round-trip de config aislado con `t.Setenv`.
- `internal/database/queries_test.go` — queries sqlc con `go-sqlmock`.
- `internal/cli/handlers_test.go` — handlers y `middlewareLoggedIn` con `go-sqlmock`.

```bash
go test ./... -v -cover
```

Dependencias de test: `github.com/DATA-DOG/go-sqlmock`.

## Licencia

Proyecto de uso académico / personal.
