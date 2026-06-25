package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"net/http"
	"time"
)

type contextKey string

const (
	sessionTokenBytes = 32
	ctxUserID         contextKey = "userID"
	ctxRole           contextKey = "role"
)

func generateSessionToken() (string, error) {
	b := make([]byte, sessionTokenBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func initSessionsTable(db *sql.DB) {
	db.Exec(`CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		user_id INTEGER NOT NULL,
		role TEXT NOT NULL,
		expires_at TEXT NOT NULL,
		created_at TEXT DEFAULT (datetime('now','localtime')),
		FOREIGN KEY (user_id) REFERENCES USUARIOS(id)
	)`)
	db.Exec("DELETE FROM sessions WHERE expires_at <= datetime('now','localtime')")
}

func createSession(db *sql.DB, userID int, role string) (token string, err error) {
	token, err = generateSessionToken()
	if err != nil {
		return "", err
	}
	expiresAt := time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05")
	_, err = db.Exec("INSERT INTO sessions (id, user_id, role, expires_at) VALUES (?, ?, ?, ?)",
		token, userID, role, expiresAt)
	if err != nil {
		return "", err
	}
	return token, nil
}

func deleteSession(db *sql.DB, token string) error {
	_, err := db.Exec("DELETE FROM sessions WHERE id=?", token)
	return err
}

var ErrNoSession = errors.New("no session")
var ErrSessionExpired = errors.New("session expired")

func validateSession(r *http.Request) (userID int, role string, err error) {
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		return 0, "", ErrNoSession
	}
	token := cookie.Value
	var uid int
	var rol string
	var expiresAt string
	err = db.QueryRow("SELECT user_id, role, expires_at FROM sessions WHERE id=?", token).Scan(&uid, &rol, &expiresAt)
	if err == sql.ErrNoRows {
		return 0, "", ErrNoSession
	}
	if err != nil {
		return 0, "", err
	}
	expires, err := time.Parse("2006-01-02 15:04:05", expiresAt)
	if err != nil || time.Now().After(expires) {
		deleteSession(db, token)
		return 0, "", ErrSessionExpired
	}
	return uid, rol, nil
}

func userIDFromContext(ctx context.Context) int {
	if v, ok := ctx.Value(ctxUserID).(int); ok {
		return v
	}
	return 0
}

func roleFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxRole).(string); ok {
		return v
	}
	return ""
}
