package config

import (
	"strings"
	"testing"
)

func TestLoadBuildsDBDSNFromDatabaseEnvFields(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("HTTP_ADDR", ":18080")
	t.Setenv("DB_HOST", "127.0.0.1")
	t.Setenv("DB_PORT", "3306")
	t.Setenv("DB_USER", "factory_user")
	t.Setenv("DB_PASSWORD", "factory_pass")
	t.Setenv("DB_NAME", "factory_count")

	cfg := Load()

	if cfg.DBHost != "127.0.0.1" {
		t.Fatalf("expected DBHost to be loaded")
	}
	if !strings.Contains(cfg.DBDSN, "factory_user:factory_pass@tcp(127.0.0.1:3306)/factory_count") {
		t.Fatalf("unexpected DBDSN: %s", cfg.DBDSN)
	}
}
