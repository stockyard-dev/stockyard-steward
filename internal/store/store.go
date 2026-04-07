package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct{ db *sql.DB }

// Expense represents a single expense entry.
// Amount is stored as cents (integer) to avoid floating-point money math.
type Expense struct {
	ID            string `json:"id"`
	Description   string `json:"description"`
	Amount        int    `json:"amount"`
	Category      string `json:"category"`
	Vendor        string `json:"vendor"`
	Date          string `json:"date"`
	PaymentMethod string `json:"payment_method"`
	Receipt       string `json:"receipt"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
}

func Open(d string) (*DB, error) {
	if err := os.MkdirAll(d, 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", filepath.Join(d, "steward.db")+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS expenses(
		id TEXT PRIMARY KEY,
		description TEXT NOT NULL,
		amount INTEGER DEFAULT 0,
		category TEXT DEFAULT '',
		vendor TEXT DEFAULT '',
		date TEXT DEFAULT '',
		payment_method TEXT DEFAULT '',
		receipt TEXT DEFAULT '',
		status TEXT DEFAULT 'pending',
		created_at TEXT DEFAULT(datetime('now'))
	)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS extras(
		resource TEXT NOT NULL,
		record_id TEXT NOT NULL,
		data TEXT NOT NULL DEFAULT '{}',
		PRIMARY KEY(resource, record_id)
	)`)
	return &DB{db: db}, nil
}

func (d *DB) Close() error { return d.db.Close() }

func genID() string { return fmt.Sprintf("%d", time.Now().UnixNano()) }
func now() string   { return time.Now().UTC().Format(time.RFC3339) }

func (d *DB) Create(e *Expense) error {
	e.ID = genID()
	e.CreatedAt = now()
	_, err := d.db.Exec(
		`INSERT INTO expenses(id, description, amount, category, vendor, date, payment_method, receipt, status, created_at)
		 VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.Description, e.Amount, e.Category, e.Vendor, e.Date, e.PaymentMethod, e.Receipt, e.Status, e.CreatedAt,
	)
	return err
}

func (d *DB) Get(id string) *Expense {
	var e Expense
	err := d.db.QueryRow(
		`SELECT id, description, amount, category, vendor, date, payment_method, receipt, status, created_at
		 FROM expenses WHERE id=?`,
		id,
	).Scan(&e.ID, &e.Description, &e.Amount, &e.Category, &e.Vendor, &e.Date, &e.PaymentMethod, &e.Receipt, &e.Status, &e.CreatedAt)
	if err != nil {
		return nil
	}
	return &e
}

func (d *DB) List() []Expense {
	rows, _ := d.db.Query(
		`SELECT id, description, amount, category, vendor, date, payment_method, receipt, status, created_at
		 FROM expenses ORDER BY date DESC, created_at DESC`,
	)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Expense
	for rows.Next() {
		var e Expense
		rows.Scan(&e.ID, &e.Description, &e.Amount, &e.Category, &e.Vendor, &e.Date, &e.PaymentMethod, &e.Receipt, &e.Status, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}

func (d *DB) Update(e *Expense) error {
	_, err := d.db.Exec(
		`UPDATE expenses SET description=?, amount=?, category=?, vendor=?, date=?, payment_method=?, receipt=?, status=?
		 WHERE id=?`,
		e.Description, e.Amount, e.Category, e.Vendor, e.Date, e.PaymentMethod, e.Receipt, e.Status, e.ID,
	)
	return err
}

func (d *DB) Delete(id string) error {
	_, err := d.db.Exec(`DELETE FROM expenses WHERE id=?`, id)
	return err
}

func (d *DB) Count() int {
	var n int
	d.db.QueryRow(`SELECT COUNT(*) FROM expenses`).Scan(&n)
	return n
}

func (d *DB) Search(q string, filters map[string]string) []Expense {
	where := "1=1"
	args := []any{}
	if q != "" {
		where += " AND (description LIKE ? OR vendor LIKE ? OR category LIKE ?)"
		args = append(args, "%"+q+"%", "%"+q+"%", "%"+q+"%")
	}
	if v, ok := filters["category"]; ok && v != "" {
		where += " AND category=?"
		args = append(args, v)
	}
	if v, ok := filters["status"]; ok && v != "" {
		where += " AND status=?"
		args = append(args, v)
	}
	rows, _ := d.db.Query(
		`SELECT id, description, amount, category, vendor, date, payment_method, receipt, status, created_at
		 FROM expenses WHERE `+where+` ORDER BY date DESC, created_at DESC`,
		args...,
	)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Expense
	for rows.Next() {
		var e Expense
		rows.Scan(&e.ID, &e.Description, &e.Amount, &e.Category, &e.Vendor, &e.Date, &e.PaymentMethod, &e.Receipt, &e.Status, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}

// Stats returns total count, total amount (in cents), and breakdowns
// by category and status. Used by the dashboard's three stat cards.
func (d *DB) Stats() map[string]any {
	m := map[string]any{
		"total":        d.Count(),
		"total_amount": 0,
		"by_category":  map[string]int{},
		"by_status":    map[string]int{},
	}

	var totalAmount int
	d.db.QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM expenses`).Scan(&totalAmount)
	m["total_amount"] = totalAmount

	if rows, _ := d.db.Query(`SELECT category, SUM(amount) FROM expenses WHERE category != '' GROUP BY category`); rows != nil {
		defer rows.Close()
		by := map[string]int{}
		for rows.Next() {
			var cat string
			var amt int
			rows.Scan(&cat, &amt)
			by[cat] = amt
		}
		m["by_category"] = by
	}

	if rows, _ := d.db.Query(`SELECT status, COUNT(*) FROM expenses GROUP BY status`); rows != nil {
		defer rows.Close()
		by := map[string]int{}
		for rows.Next() {
			var s string
			var c int
			rows.Scan(&s, &c)
			by[s] = c
		}
		m["by_status"] = by
	}

	return m
}

// ─── Extras: generic key-value storage for personalization custom fields ───

func (d *DB) GetExtras(resource, recordID string) string {
	var data string
	err := d.db.QueryRow(
		`SELECT data FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	).Scan(&data)
	if err != nil || data == "" {
		return "{}"
	}
	return data
}

func (d *DB) SetExtras(resource, recordID, data string) error {
	if data == "" {
		data = "{}"
	}
	_, err := d.db.Exec(
		`INSERT INTO extras(resource, record_id, data) VALUES(?, ?, ?)
		 ON CONFLICT(resource, record_id) DO UPDATE SET data=excluded.data`,
		resource, recordID, data,
	)
	return err
}

func (d *DB) DeleteExtras(resource, recordID string) error {
	_, err := d.db.Exec(
		`DELETE FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	)
	return err
}

func (d *DB) AllExtras(resource string) map[string]string {
	out := make(map[string]string)
	rows, _ := d.db.Query(
		`SELECT record_id, data FROM extras WHERE resource=?`,
		resource,
	)
	if rows == nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id, data string
		rows.Scan(&id, &data)
		out[id] = data
	}
	return out
}
