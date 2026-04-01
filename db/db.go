package db

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func seedRoles() {
	_, _ = DB.Exec(`INSERT OR IGNORE INTO roles (id, name) VALUES (1,'admin')`)
	_, _ = DB.Exec(`INSERT OR IGNORE INTO roles (id, name) VALUES (2,'user')`)
}
func seedUser() {
	var count int
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("1111"), bcrypt.DefaultCost)
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", "root@example.ru").Scan(&count)
	if err != nil {
		log.Println("Seed check error:", err)
		return
	}

	if count > 0 {
		return // уже есть
	}

	_, err = DB.Exec(`
		INSERT INTO users (name, login, email, age, country, password, role_id)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
		"Root",
		"root",
		"root@example.ru",
		30,
		"RU",
		string(hashedPassword),
		1,
	)

	if err != nil {
		log.Println("Seed insert error:", err)
		return
	}

	log.Println("Seed user created")
}

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "users.db")
	if err != nil {
		log.Fatal(err)
	}
	createTable()
	seedRoles()
	seedUser() // 👈 добавили
}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    login TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    age INTEGER,
    country TEXT,
    password TEXT NOT NULL,
    role_id INTEGER NOT NULL,
    FOREIGN KEY(role_id) REFERENCES roles(id)
);
CREATE TABLE IF NOT EXISTS roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);
`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
