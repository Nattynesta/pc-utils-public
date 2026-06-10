# Red Tailscale — Estado y Comandos

## Nodos Registrados (Tailnet: lab238nesta@)

| IP (100.x.y.z) | Hostname | OS | Estado | Última vez |
|----------------|----------|-----|--------|------------|
| **100.92.186.120** | agentelinux-thinnote14 | Linux | **activo** (este servidor POS) | now |
| 100.124.160.13 | agentelinux-hp-probook-6475b | Linux | offline | — |
| 100.108.145.29 | desktop-32gvmvp | Windows | offline | — |
| 100.119.11.79 | desktop-ch99s8o | Windows | online | — |
| 100.65.114.99 | desktop-n2pmeqq | Windows | offline | — |
| 100.97.245.113 | redmi-note-12 | Android | idle | tx 9.5KB / rx 10.9KB |

---

## Servicios Expuestos en 100.92.186.120

| Puerto | Servicio | Acceso |
|--------|----------|--------|
| 8080 | Abarrotes POS (Go) | `http://100.92.186.120:8080` |
| — | — | — |

**Acceso POS desde cualquier nodo**: `http://100.92.186.120:8080`

---

## Comandos Tailscale (Snap)

```bash
# Status
/snap/tailscale/115/bin/tailscale --socket=/var/snap/tailscale/common/socket/tailscaled.sock status

# Up (requiere sudo)
sudo /snap/tailscale/115/bin/tailscale --socket=/var/snap/tailscale/common/socket/tailscaled.sock up

# IP propia
/snap/tailscale/115/bin/tailscale --socket=/var/snap/tailscale/common/socket/tailscaled.sock ip -4

# Ping a nodo
/snap/tailscale/115/bin/tailscale --socket=/var/snap/tailscale/common/socket/tailscaled.sock ping 100.119.11.79

# File sharing (opcional)
/snap/tailscale/115/bin/tailscale file cp archivo.txt 100.119.11.79:
```

---

## Configuración Actual

- **Daemon**: `/snap/tailscale/115/bin/tailscaled` (systemd/snap)
- **Socket**: `/var/snap/tailscale/common/socket/tailscaled.sock`
- **IP asignada**: 100.92.186.120/32 (CGNAT)
- **IPv6 ULA**: fd7a:115c:a1e0::b34:ba79/128
- **Firewall host**: ufw inactivo
- **Puerto POS**: 8080 en `*:8080` (todas interfaces)

---

## Troubleshooting

### "Connection refused" a 100.92.186.120:8080
1. Verificar servidor corriendo: `ss -tlnp | grep 8080`
2. Verificar binding: `*:8080` (no `127.0.0.1:8080`)
3. Firewall host: `ufw status` → inactivo OK
4. Tailscale permite tráfico entre nodos por defecto (ACLs default allow)

### Daemon no responde
```bash
# Verificar proceso
ps aux | grep tailscaled

# Reiniciar daemon (snap)
sudo snap restart tailscale

# Ver logs
journalctl -u snap.tailscale.tailscaled -f
```

### Auth URL no aparece
```bash
# Forzar auth interactivo
sudo /snap/tailscale/115/bin/tailscale --socket=/var/snap/tailscale/common/socket/tailscaled.sock up --print-auth-url
```

---

## Acceso desde Otros Dispositivos

### Navegador
```
http://100.92.186.120:8080
```

### curl / scripts
```bash
curl http://100.92.186.120:8080/api/productos
curl -X POST -d "usuario=admin&clave=admin" http://100.92.186.120:8080/login
```

### SSH (si habilitado)
```bash
# En cliente: tailscale ssh user@100.92.186.120
# Requiere: tailscale set --ssh=enabled
```

---

## Próximos Pasos Red

- [ ] Habilitar SSH via Tailscale (`tailscale set --ssh=enabled`)
- [ ] Configurar ACLs en admin console (restringir puertos)
- [ ] MagicDNS: acceder como `http://thinnote14:8080` o `http://thinnote14.lab238nesta.ts.net:8080`
- [ ] Subnet router para LAN local (si aplica)
- [ ] Exit node para tráfico saliente (opcional)