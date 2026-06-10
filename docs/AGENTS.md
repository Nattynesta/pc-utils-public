# Abarrotes PDV — Guía para Agentes

## Resumen del Proyecto
Sistema POS web en Go + SQLite para abarrotes. Migración desde Firebird (PDVDATA.FDB).

### Stack
- **Backend**: Go 1.22+ (net/http, html/template, sqlite3)
- **Frontend**: HTMX 2, Font Awesome 6, Tabler Icons
- **DB**: SQLite (WAL mode) en `~/.abarrotes-pdv/pdv.db`
- **Auth**: Cookies `session` (usuario) + `role` (admin/helper)

### Credenciales
| Usuario | Pass | Rol | Acceso |
|---------|------|-----|--------|
| admin | admin | admin | Todo |
| helper | helper | helper | Ventas, POS, Tickets |

### URLs
- Local: `http://localhost:8080`
- Tailscale: `http://100.92.186.120:8080`

---

## Arquitectura Clave

### Roles (Middleware `withAdmin`)
```go
// main.go:176
func withAdmin(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        roleCookie, _ := r.Cookie("role")
        if roleCookie.Value != "admin" {
            http.Redirect(w, r, "/ventas/pos", http.StatusSeeOther)
            return
        }
        next(w, r)
    }
}
```
- Rutas admin: `/productos*`, `/clientes*`, `/proveedores`, `/reportes`, `/usuarios`, `/cajas`
- Helper redirigido a `/ventas/pos` al login

### Renderizado de Templates
```go
// main.go:233 - firma cambió para aceptar *http.Request
func render(w http.ResponseWriter, r *http.Request, name string, data PageData) {
    // Auto-puebla data.User y data.Role desde cookies
}
```
- `PageData` incluye `Role string` para nav condicional en `base.html`

### Migración Firebird → SQLite
- **Encoding crítico**: Firebird usa `latin-1` (cp1252). Leer con `utf-8` + `errors="replace"` rompe ñ/acentos.
- Fix: `migrate_full.py` usa `encoding="latin-1"` al leer CSV/export.
- Tabla `USUARIOS` añade columna `rol TEXT DEFAULT 'helper'`.

### Iconos Productos (POS)
- **Tabler Icons** (MIT) vía CDN: `ti ti-droplet`, `ti ti-bottle`, etc.
- 13 categorías mapeadas por palabras clave en descripción
- **OpenFoodFacts**: 34/271 productos con foto real (`off_image_url`, `off_image_small`)
- Fallback: Tabler icon si no hay foto

---

## Comandos Útiles

```bash
# Build
cd pos && go build -o pos_server .

# Run (foreground)
./pos_server

# Run (background)
nohup ./pos_server > pos.log 2>&1 &

# Kill
pkill -f pos_server

# Test API
curl http://127.0.0.1:8080/api/productos | jq '.[] | select(.off_image_url)'

# Tailscale
sudo /snap/tailscale/115/bin/tailscale --socket=/var/snap/tailscale/common/socket/tailscaled.sock status
sudo /snap/tailscale/115/bin/tailscale --socket=/var/snap/tailscale/common/socket/tailscaled.sock up

# DB
sqlite3 ~/.abarrotes-pdv/pdv.db "SELECT codigo, descripcion, off_image_url FROM PRODUCTOS LIMIT 5;"
```

---

## Compresión de Tokens (para contexto limitado)

### Reglas de salida
- Truncar logs/lists a primeras 10 + últimas 5 líneas
- `git diff` → solo headers de hunk + 3 líneas
- `ls/glob` >30 entries → mostrar count ("47 files omitted")
- JSON largo → primeras 2 + últimas 2 líneas
- HTML → strip tags, solo texto relevante

### Herramientas preferidas
- `glob` > `find/ls` para búsqueda archivos
- `grep` (ripgrep) > `bash grep` para contenido
- `read` con offset/limit > `cat/head/tail`
- `edit` > `sed/awk` para cambios
- `Task` agent para exploración multi-paso

---

## Efectos Red Tailscale (observados)

| Nodo | IP | Estado |
|------|-----|--------|
| agentelinux-thinnote14 | 100.92.186.120 | **activo (este servidor)** |
| agentelinux-hp-probook-6475b | 100.124.160.13 | offline |
| desktop-32gvmvp | 100.108.145.29 | offline |
| desktop-ch99s8o | 100.119.11.79 | online |
| desktop-n2pmeqq | 100.65.114.99 | offline |
| redmi-note-12 | 100.97.245.113 | idle |

**Acceso POS**: `http://100.92.186.120:8080` desde cualquier nodo conectado.

**Firewall**: ufw inactivo. Puerto 8080 escucha en `*:8080` (todas interfaces).

---

## Próximos Pasos Pendientes

- [ ] Cache TTL para OpenFoodFacts (revalidar semanal)
- [ ] Selector manual de icono en formulario producto
- [ ] Tests unitarios (Go test)
- [ ] CI/CD GitHub Actions
- [ ] Backup automático DB