package db

import (
	"database/sql"
	"fmt"

	"github.com/098765432m/config"
	"github.com/098765432m/logger"
	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(cfg *config.DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=%t",
	cfg.User,
	cfg.Password,
	cfg.Host,
	cfg.Port,
	cfg.Name,
	cfg.ParseTime,
	)

	// logger.NewLogger().Info.Println(dsn)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	
	logger.NewLogger().Info.Println("Connect to database successfully!")

	database := &Database{db:db}

	database.InitDb()

	return database, nil
}

func (d *Database) InitDb() {
	query := `CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		email VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);`

	_, err := d.db.Exec(query)
	if err != nil {
		logger.NewLogger().Error.Fatal("Failed to create users table:", err)
	}
}

func (d *Database) CloseDb() {
	if d.db != nil {
		d.db.Close()
		logger.NewLogger().Info.Println("Database connection closed.")
	}
}