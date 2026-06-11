package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"golang.org/x/text/encoding/charmap"

	_ "modernc.org/sqlite"
)

func main() {
	dbPath := os.Args[1]

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("open: %v", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT rowid, codigo, descripcion FROM PRODUCTOS")
	if err != nil {
		log.Fatalf("query: %v", err)
	}
	defer rows.Close()

	enc := charmap.Windows1252.NewDecoder()
	fixed := 0
	for rows.Next() {
		var rowid int
		var codigo, descripcion string
		rows.Scan(&rowid, &codigo, &descripcion)

		raw := []byte(descripcion)
		hasHigh := false
		for _, b := range raw {
			if b > 127 {
				hasHigh = true
				break
			}
		}
		if !hasHigh {
			continue
		}

		utf8Bytes, err := enc.Bytes(raw)
		if err != nil {
			log.Printf("Skip %s: %v", codigo, err)
			continue
		}
		utf8Str := string(utf8Bytes)
		if utf8Str == descripcion {
			continue
		}

		db.Exec("UPDATE PRODUCTOS SET descripcion = ? WHERE rowid = ?", utf8Str, rowid)
		fmt.Printf("Fixed %s: %q -> %q\n", codigo, descripcion, utf8Str)
		fixed++
	}
	fmt.Printf("Total fixed: %d\n", fixed)
}
