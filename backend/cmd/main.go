package main

import (
	"log"

	"count/backend/config"
	"count/backend/internal/handler"
	"count/backend/internal/model"
	"count/backend/internal/router"
	"count/backend/internal/seed"
	"count/backend/pkg/db"
)

func main() {
	cfg := config.Load()

	mysqlDB, err := db.OpenMySQL(cfg.DBDSN)
	if err != nil {
		log.Fatalf("open db failed: %v", err)
	}

	if err := mysqlDB.AutoMigrate(
		&model.User{},
		&model.Process{},
		&model.WorkerProcess{},
		&model.Order{},
		&model.OrderItem{},
		&model.OrderProcess{},
		&model.WageRule{},
		&model.ScanRecord{},
	); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}
	if err := mysqlDB.Migrator().AlterColumn(&model.User{}, "StationRole"); err != nil {
		log.Fatalf("alter users.station_role failed: %v", err)
	}
	if err := mysqlDB.Migrator().AlterColumn(&model.Process{}, "StationRole"); err != nil {
		log.Fatalf("alter processes.station_role failed: %v", err)
	}
	if err := mysqlDB.Migrator().AlterColumn(&model.OrderProcess{}, "StationRole"); err != nil {
		log.Fatalf("alter order_processes.station_role failed: %v", err)
	}
	if err := mysqlDB.Migrator().AlterColumn(&model.ScanRecord{}, "StationRole"); err != nil {
		log.Fatalf("alter scan_records.station_role failed: %v", err)
	}
	if err := seed.BackfillLegacyRoles(mysqlDB); err != nil {
		log.Fatalf("backfill legacy roles failed: %v", err)
	}
	if err := seed.Defaults(mysqlDB); err != nil {
		log.Fatalf("seed defaults failed: %v", err)
	}

	h := handler.New(mysqlDB, cfg.JWTSecret)
	r := router.New(h)
	if err := r.Run(cfg.HTTPAddr); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
