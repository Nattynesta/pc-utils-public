# Token Optimization — Reglas de Compresión

## Principio
Presupuesto limitado (8K-128K ctx, ~500K tokens/día). Cada token cuenta.

---

## Reglas de Salida (Output)

### General
- Eliminar trailing whitespace y líneas vacías
- Deduplicar líneas consecutivas idénticas
- Colapsar logs `npm install` / `cargo build` → solo errores y última línea
- Truncar `ls -la`, `glob` >30 entries → `"47 files omitted"`
- Rutas: `/home/agentelinux` → `~`

### Git
- `git status` → solo listas modified/untracked
- `git diff` → headers de hunk + 3 líneas c/u
- `git log` → máx 10 commits, one-liner
- `git branch` → solo marcador `*`
- `git stash list` → últimas 5

### npm/pnpm/bun
- Eliminar "up to date", "audited X packages", "found X vulnerabilities"
- Solo errores/warnings reales
- Árbol deps → profundidad 1

### ls / tree / glob
- Dirs >30 entries → count
- `tree` → max depth 3, count por dir
- `ls -la` → primeras 10 + últimas 5

### Docker
- `docker ps` → names + status
- `docker images` → repo:tag + size
- `docker logs` → últimas 20 líneas, grep ERROR/WARN

### Cargo
- `cargo build` → solo errores `error[E...]`
- `cargo test` → nombres tests + pass/fail (sin backtrace)
- `cargo check` → solo errores reales

### WebFetch / File reads
- HTML → strip tags, solo texto
- JSON logs → primeras 2 + últimas 20 líneas
- Archivos >200 líneas → leer en chunks, resumir secciones

### Grep
- >20 matches → count + primeros 10
- Usar `-l` (solo archivos) cuando posible, luego leer matches individuales

---

## Presupuesto de Tokens
- Respuestas < 100 líneas salvo que se pida detalle
- Resumir antes de enviar texto grande al modelo
- Cerca de límites → `/compact` o pedir resumir
- Preferir `read` con offset/limit vs leer archivo completo

---

## Herramientas Preferidas (orden)

| Tarea | Herramienta | Evitar |
|-------|-------------|--------|
| Buscar archivos | `glob` | `find`, `ls -R` |
| Buscar contenido | `grep` (rg) | `bash grep`, `grep -r` |
| Leer archivo | `read` (offset/limit) | `cat`, `head`, `tail` |
| Editar archivo | `edit` | `sed`, `awk`, `perl` |
| Escribir archivo | `write` | `echo >`, `cat <<EOF` |
| Explorar codebase | `Task` (explore agent) | Múltiples `bash`/`grep` |
| Comunicar | Output directo | `echo` en bash, comentarios en código |

---

## Ejemplos de Compresión

### Antes (200 tokens)
```bash
$ ls -la /home/agentelinux/repos/pc-utils/pos/
total 17764
drwxr-xr-x  4 agentelinux agentelinux     4096 Jun 10 14:49 .
drwxr-xr-x  5 agentelinux agentelinux     4096 Jun 10 12:58 ..
-rw-r--r--  1 agentelinux agentelinux     7322 Jun 10 14:40 db.go
-rw-r--r--  1 agentelinux agentelinux      707 Jun 10 14:13 go.mod
...
```

### Después (30 tokens)
```
pos/ : 18 files (db.go, handlers.go, main.go, templates/, static/, pos_server)
```

### Antes (JSON 500 tokens)
```json
{"products": [{"codigo": "7501055300952", "descripcion": "Coca 1.25L", ...}, {...}...]}
```

### Después (50 tokens)
```
API productos: 271 items. 34 con off_image_url (Coca, Bonafont, Oreo, etc.)
```

---

## Checklist Pre-Respuesta
- [ ] ¿Puedo responder en < 5 líneas?
- [ ] ¿He truncado listas/logs?
- [ ] ¿He eliminado ruido (up-to-date, audited, etc.)?
- [ ] ¿Rutas abreviadas (`~`)?
- [ ] ¿JSON/HTML largo → resumen + muestra?