package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Expense struct {
	ID string `json:"id"`
	Description string `json:"description"`
	Amount int `json:"amount"`
	Category string `json:"category"`
	Vendor string `json:"vendor"`
	Date string `json:"date"`
	PaymentMethod string `json:"payment_method"`
	Receipt string `json:"receipt"`
	Status string `json:"status"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"steward.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS expenses(id TEXT PRIMARY KEY,description TEXT NOT NULL,amount INTEGER DEFAULT 0,category TEXT DEFAULT '',vendor TEXT DEFAULT '',date TEXT DEFAULT '',payment_method TEXT DEFAULT '',receipt TEXT DEFAULT '',status TEXT DEFAULT 'pending',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Expense)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO expenses(id,description,amount,category,vendor,date,payment_method,receipt,status,created_at)VALUES(?,?,?,?,?,?,?,?,?,?)`,e.ID,e.Description,e.Amount,e.Category,e.Vendor,e.Date,e.PaymentMethod,e.Receipt,e.Status,e.CreatedAt);return err}
func(d *DB)Get(id string)*Expense{var e Expense;if d.db.QueryRow(`SELECT id,description,amount,category,vendor,date,payment_method,receipt,status,created_at FROM expenses WHERE id=?`,id).Scan(&e.ID,&e.Description,&e.Amount,&e.Category,&e.Vendor,&e.Date,&e.PaymentMethod,&e.Receipt,&e.Status,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Expense{rows,_:=d.db.Query(`SELECT id,description,amount,category,vendor,date,payment_method,receipt,status,created_at FROM expenses ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Expense;for rows.Next(){var e Expense;rows.Scan(&e.ID,&e.Description,&e.Amount,&e.Category,&e.Vendor,&e.Date,&e.PaymentMethod,&e.Receipt,&e.Status,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM expenses WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM expenses`).Scan(&n);return n}
