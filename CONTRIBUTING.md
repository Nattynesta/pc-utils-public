# Contribuyendo a Abarrotes PDV

## Branches

- `main` — estable, protegido. Solo merge via PR.
- `develop` — integración. branch base para features.
- `feat/<nombre>` — nuevas features. Se mergean a `develop`.
- `fix/<nombre>` — bugfixes. Se mergean a `develop`.
- `docs/<nombre>` — documentación.

## Flujo de trabajo

1. Crear un Issue describiendo el cambio
2. Crear branch desde `develop`
3. Implementar con commits atómicos
4. Push y abrir Pull Request a `develop`
5. Esperar revisión + CI verde
6. Merge

## Commits

Usamos [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: agregar selector de icono en formulario producto
fix: encoding de ñ en exportación CSV
docs: actualizar setup en README
refactor: extraer lógica de búsqueda a helper
test: agregar tests unitarios para handlers
```

## Estándares de código

- Go: `gofmt` + `go vet` pasando siempre
- HTML: indentación 2 espacios
- Sin comentarios en código (el código se explica solo)
- Nombres en español para dominio del negocio (productos, ventas, tickets)
- Nombres en inglés para infraestructura (handlers, middleware, db)

## Pull Requests

- Título descriptivo con prefix (feat/fix/docs/refactor/test)
- Descripción: qué cambia y por qué
- Link al Issue correspondiente
- CI debe pasar (build + test)
- Sin merge conflicts con `develop`

## API

Nuevos endpoints siguen el patrón REST existente:

```go
mux.HandleFunc("GET /api/recurso", handleRecursoList)
mux.HandleFunc("POST /api/recurso", handleRecursoCreate)
```

Respuestas siempre JSON con `jsonResp` o `jsonErr`.
