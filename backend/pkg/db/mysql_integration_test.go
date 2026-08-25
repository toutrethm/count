package db

import (
	"database/sql"
	"os"
	"testing"

	"count/backend/config"

	"github.com/joho/godotenv"
)

func TestOpenMySQLFromEnvCanPingDatabase(t *testing.T) {
	if os.Getenv("RUN_DB_INTEGRATION") != "1" {
		t.Skip("set RUN_DB_INTEGRATION=1 to run the live database connection test")
	}

	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatalf("load env: %v", err)
	}

	cfg := config.Load()
	if cfg.DBDSN == "" {
		t.Fatal("DBDSN is empty; check DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME")
	}

	gormDB, err := OpenMySQL(cfg.DBDSN)
	if err != nil {
		t.Fatalf("open mysql: %v", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		t.Fatalf("get sql db: %v", err)
	}
	defer func(conn *sql.DB) {
		_ = conn.Close()
	}(sqlDB)

	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("ping mysql: %v", err)
	}
}
