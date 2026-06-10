# Abarrotes PDV — Estado Actual

## Versión
- Commit: `7374702` (HEAD)
- Build: `pos_server` (Go 1.22+)

---

## Funcionalidades Implementadas

### 1. Sistema de Roles (admin/helper)
```
admin/admin  → acceso total (productos, clientes, proveedores, reportes, usuarios, cajas)
helper/helper → solo Ventas, POS, Tickets (redirige a /ventas/pos al login)
```

### 2. Iconografía
| Ubicación | Librería | Detalle |
|-----------|----------|---------|
| Nav, botones, forms | Font Awesome 6 | `fa-store`, `fa-cash-register`, `fa-box`, `fa-users`, etc. |
| Productos POS | Tabler Icons (MIT) | 13 categorías: `ti ti-droplet`, `ti ti-bottle`, `ti ti-bread`, etc. |
| Fallback productos | Tabler | Por palabras clave en descripción |

### 3. OpenFoodFacts Integration
- **34/271 productos** con foto real desde OpenFoodFacts
- Tabla `PRODUCTOS_OFF` (codigo PK, image_url, image_small, name, nutriscore, nova)
- API `/api/productos` incluye `off_image_url`, `off_image_small`, `off_name`
- POS: muestra `<img>` circular si existe, fallback a icono Tabler
- Ejemplos con foto: Coca-Cola (10 variantes), Bonafont, Oreo, Doritos, Gatorade, etc.

### 4. Fix Encoding ñ/acentos
- **Problema**: Firebird usa `latin-1` (cp1252). Lectura `utf-8` + `errors="replace"` = �
- **Fix**: `migrate_full.py` usa `encoding="latin-1"` al leer CSV
- Re-migración exitosa: ñ, á, é, í, ó, ú, ü se muestran correctamente

### 5. Base de Datos
- **Ubicación**: `~/.abarrotes-pdv/pdv.db` (SQLite WAL mode)
- **Tablas clave**: USUARIOS (rol), PRODUCTOS, PRODUCTOS_OFF, VENTATICKETS, OPERACIONES, CLIENTES
- **Migración**: `ALTER TABLE USUARIOS ADD COLUMN rol TEXT DEFAULT 'helper'` (idempotente via pragma)

---

## API Endpoints

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/api/productos` | Lista productos (+ off_image_*) |
| GET | `/api/productos?q=` | Buscar por código/descripción |
| GET | `/api/clientes` | Lista clientes |
| GET | `/api/usuarios` | Lista usuarios (+ rol) |
| POST | `/api/usuarios` | Crear usuario (incluye rol) |
| PUT | `/api/usuarios/:id` | Actualizar usuario (incluye rol) |
| GET | `/api/tickets` | Lista tickets |
| POST | `/api/tickets` | Crear ticket (venta) |
| GET | `/api/departamentos` | Lista departamentos |

---

## Flujo Helper (Cajero)

1. Login `helper` / `helper`
2. Redirige automático a `/ventas/pos`
3. Nav muestra solo: **Ventas**, **POS**, **Tickets**
4. Bloqueado: Productos, Clientes, Proveedores, Reportes, Usuarios, Cajas (303 → POS)
5. Solo POS para vender, Tickets para ver historial

---

## Flujo Admin

1. Login `admin` / `admin`
2. Redirige a `/` (Dashboard)
3. Nav completa con todos los módulos
4. Puede crear usuarios con rol `admin` o `helper`

---

## Pendientes / Mejoras

| Área | Tarea |
|------|-------|
| OFF Cache | TTL + revalidación semanal (cron/job) |
| Iconos | Selector manual en formulario producto |
| Tests | `go test ./...` (unit + integration) |
| CI/CD | GitHub Actions (build, test, lint) |
| Backup | DB dump automático diario |
| UX | Teclado numérico en POS móvil |
| Sync | Offline-first con IndexedDB (futuro) |

---

## Estructura Archivos Clave

```
pos/
├── main.go              # Server, middleware, render, migrate
├── handlers.go          # API endpoints (productos, usuarios, tickets...)
├── handlers_pages.go    # Page handlers (login, POS, dashboard...)
├── db.go                # Structs + queries (Producto + OffImageUrl)
├── schema.sql           # DDL (USUARIOS.rol, PRODUCTOS_OFF)
├── migrate_full.py      # Firebird→SQLite (latin-1 encoding)
├── templates/
│   ├── base.html        # Nav condicional por .Role, Font Awesome + Tabler CSS
│   ├── login.html       # FA icon
│   └── ventas/pos.html  # Product cards con off_image / Tabler icons
└── static/              # (vacío, CDN only)
```

---

## Ejecución

```bash
# Desarrollo
cd pos && go run .

# Producción
cd pos && go build -o pos_server .
nohup ./pos_server > pos.log 2>&1 &

# Verificar
curl http://100.92.186.120:8080/login
curl http://100.92.186.120:8080/api/productos | jq '.[] | select(.off_image_url)'
```