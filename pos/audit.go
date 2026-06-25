package main

import (
	"database/sql"
	"net/http"
)

func initAuditTable(db *sql.DB) {
	db.Exec(`CREATE TABLE IF NOT EXISTS audit_log (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		action TEXT NOT NULL,
		resource_type TEXT NOT NULL,
		resource_id INTEGER,
		details TEXT,
		ip_address TEXT,
		created_at TEXT DEFAULT (datetime('now','localtime')),
		FOREIGN KEY (user_id) REFERENCES USUARIOS(id)
	)`)
}

func logAudit(db *sql.DB, userID int, action, resourceType string, resourceID int, details, ip string) {
	db.Exec(
		"INSERT INTO audit_log (user_id, action, resource_type, resource_id, details, ip_address) VALUES (?, ?, ?, ?, ?, ?)",
		userID, action, resourceType, resourceID, details, ip,
	)
}

func getUserIDForAudit(r *http.Request) int {
	if uid := userIDFromContext(r.Context()); uid > 0 {
		return uid
	}
	return 0
}
