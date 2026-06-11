# pc-utils — Abarrotes PDV

Sistema POS web para abarrotes. Go + SQLite + HTMX.

## Stack

| Capa | Tecnología |
|------|-----------|
| Backend | Go 1.22+ (net/http, html/template) |
| Frontend | HTMX 2, Font Awesome 6, Tabler Icons |
| Base de datos | SQLite (WAL mode) |
| Auth | Cookies de sesión + roles (admin/helper) |

## Arquitectura

```
pos/
├── main.go              # Server, middleware, rutas, migración
├── handlers.go          # API endpoints (productos, tickets, usuarios...)
├── handlers_pages.go    # Page handlers (login, POS, dashboard...)
├── db.go                # Structs + queries SQL
├── schema.sql           # DDL
├── templates/
│   ├── base.html        # Layout con nav condicional por rol
│   ├── login.html       # Pantalla de login
│   └── ventas/pos.html  # POS con búsqueda y carrito
└── static/
    └── style.css
```

## Requisitos

- Go 1.22+
- Git

## Desarrollo

```bash
# Clonar
git clone https://github.com/Nattynesta/pc-utils.git
cd pc-utils/pos

# Compilar
go build -o pos_server .

# Ejecutar (localhost:8080)
./pos_server

# O en modo desarrollo con recarga manual
go run .
```

La base de datos se crea automáticamente en `~/.abarrotes-pdv/pdv.db`.

## Credenciales por defecto

| Usuario | Pass   | Rol    | Acceso                        |
|---------|--------|--------|-------------------------------|
| admin   | admin  | admin  | Todo el sistema               |
| helper  | helper | helper | Solo ventas, POS, tickets     |

## API

```
GET    /api/productos[?q=]    Listar / buscar productos
POST   /api/productos         Crear producto
PUT    /api/productos/{cod}   Actualizar producto
DELETE /api/productos/{cod}   Eliminar producto

GET    /api/tickets            Listar tickets
POST   /api/tickets            Crear ticket (venta)
POST   /api/tickets/{id}/cobrar   Cobrar ticket
POST   /api/tickets/{id}/cancelar Cancelar ticket

GET    /api/clientes[?q=]     Listar / buscar clientes
GET    /api/usuarios          Listar usuarios
...
```

## Contribuir

Ver [CONTRIBUTING.md](./CONTRIBUTING.md).

## Licencia

MIT — ver [LICENSE](./LICENSE).
