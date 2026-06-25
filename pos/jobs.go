package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
)

type JobType string

const (
	JobExportVentas    JobType = "export_ventas"
	JobExportInventario JobType = "export_inventario"
	JobReporteCaja     JobType = "reporte_caja"
)

type Job struct {
	ID        string                 `json:"id"`
	Type      JobType                `json:"type"`
	Params    map[string]interface{} `json:"params"`
	Status    string                 `json:"status"`
	Result    string                 `json:"result,omitempty"`
	UserID    int                    `json:"user_id"`
	CreatedAt string                 `json:"created_at"`
}

var (
	jobQueue   = make(chan Job, 100)
	jobOnce    sync.Once
)

func initJobsDB(db *sql.DB) {
	jobOnce.Do(func() {
		db.Exec(`
			CREATE TABLE IF NOT EXISTS jobs (
				id TEXT PRIMARY KEY,
				type TEXT NOT NULL,
				params TEXT,
				status TEXT DEFAULT 'pending',
				result TEXT,
				user_id INTEGER,
				created_at DATETIME DEFAULT (datetime('now','localtime')),
				completed_at DATETIME
			)
		`)
		for i := 0; i < 3; i++ {
			go jobWorker(db)
		}
	})
}

func jobWorker(db *sql.DB) {
	for job := range jobQueue {
		processJob(db, job)
	}
}

func processJob(db *sql.DB, job Job) {
	db.Exec("UPDATE jobs SET status = 'processing' WHERE id = ?", job.ID)

	var result string
	var err error

	switch job.Type {
	case JobExportVentas:
		result, err = exportVentas(db, job.Params)
	case JobExportInventario:
		result, err = exportInventario(db, job.Params)
	case JobReporteCaja:
		result, err = generarReporteCaja(db, job.Params)
	}

	status := "completed"
	if err != nil {
		status = "failed"
		result = err.Error()
	}

	db.Exec(
		"UPDATE jobs SET status = ?, result = ?, completed_at = datetime('now','localtime') WHERE id = ?",
		status, result, job.ID,
	)
}

func generateJobID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func exportVentas(db *sql.DB, params map[string]interface{}) (string, error) {
	desde, _ := params["desde"].(string)
	hasta, _ := params["hasta"].(string)

	query := `SELECT t.id, t.folio, t.total, t.creado_en, COALESCE(c.nombre,'') as cliente
		FROM VENTATICKETS t LEFT JOIN CLIENTES c ON t.cliente_id=c.numero
		WHERE t.esta_cancelado='f'`

	var args []interface{}
	if desde != "" {
		query += " AND t.creado_en >= ?"
		args = append(args, desde)
	}
	if hasta != "" {
		query += " AND t.creado_en <= ?"
		args = append(args, hasta+" 23:59:59")
	}
	query += " ORDER BY t.creado_en"

	rows, err := db.Query(query, args...)
	if err != nil {
		return "", fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	type row struct {
		ID      int     `json:"id"`
		Folio   *int    `json:"folio"`
		Total   float64 `json:"total"`
		Fecha   string  `json:"fecha"`
		Cliente string  `json:"cliente"`
	}
	var data []row
	for rows.Next() {
		var r row
		if err := rows.Scan(&r.ID, &r.Folio, &r.Total, &r.Fecha, &r.Cliente); err != nil {
			continue
		}
		data = append(data, r)
	}

	b, _ := json.Marshal(data)
	return string(b), nil
}

func exportInventario(db *sql.DB, params map[string]interface{}) (string, error) {
	categoria, _ := params["categoria"].(string)

	query := `SELECT codigo, descripcion, COALESCE(dinventario,0), COALESCE(pcosto,0), COALESCE(pventa,0)
		FROM PRODUCTOS WHERE activo=1`
	var args []interface{}
	if categoria != "" {
		query += " AND categoria=?"
		args = append(args, categoria)
	}
	query += " ORDER BY descripcion"

	rows, err := db.Query(query, args...)
	if err != nil {
		return "", fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	type row struct {
		Codigo      string  `json:"codigo"`
		Descripcion string  `json:"descripcion"`
		Inventario  float64 `json:"inventario"`
		Costo       float64 `json:"costo"`
		Precio      float64 `json:"precio"`
	}
	var data []row
	for rows.Next() {
		var r row
		if err := rows.Scan(&r.Codigo, &r.Descripcion, &r.Inventario, &r.Costo, &r.Precio); err != nil {
			continue
		}
		data = append(data, r)
	}

	b, _ := json.Marshal(data)
	return string(b), nil
}

func generarReporteCaja(db *sql.DB, params map[string]interface{}) (string, error) {
	operacionID, _ := params["operacion_id"].(float64)

	var op Operacion
	err := db.QueryRow(`SELECT id, dinero_en_caja, tipo_de_cambio, inicio_usuario_id, inicio_en, cerro_en, caja_id, abierta FROM OPERACIONES WHERE id=?`, int(operacionID)).
		Scan(&op.ID, &op.DineroEnCaja, &op.TipoDeCambio, &op.InicioUsuarioID, &op.InicioEn, &op.CerroEn, &op.CajaID, &op.Abierta)
	if err != nil {
		return "", fmt.Errorf("operacion no encontrada: %w", err)
	}

	var totalVentas, totalGanancias float64
	db.QueryRow(`SELECT COALESCE(SUM(total),0), COALESCE(SUM(ganancia),0) FROM VENTATICKETS WHERE operacion_id=? AND esta_cancelado='f'`, int(operacionID)).
		Scan(&totalVentas, &totalGanancias)

	result := map[string]interface{}{
		"operacion_id":  op.ID,
		"inicio":       op.InicioEn,
		"cierre":       op.CerroEn,
		"ventas":       totalVentas,
		"ganancias":    totalGanancias,
	}

	b, _ := json.Marshal(result)
	return string(b), nil
}
