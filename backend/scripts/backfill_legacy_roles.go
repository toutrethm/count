package main

import (
	"log"

	"count/backend/config"
	"count/backend/internal/model"
	"count/backend/internal/seed"
	"count/backend/pkg/db"
	"gorm.io/gorm"
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

	log.Printf("legacy role rows remaining: users=%d order_processes=%d scan_records=%d",
		countLegacyUsers(mysqlDB),
		countLegacyOrderProcesses(mysqlDB),
		countLegacyScanRecords(mysqlDB),
	)
}

func countLegacyUsers(db *gorm.DB) int64 {
	var count int64
	_ = db.Model(&model.User{}).Where("station_role IN ?", []string{"A", "B", "C", "A,B", "A,C", "B,C", "A,B,C"}).Count(&count).Error
	return count
}

func countLegacyOrderProcesses(db *gorm.DB) int64 {
	var count int64
	_ = db.Model(&model.OrderProcess{}).Where("station_role IN ?", []string{"A", "B", "C", "A,B", "A,C", "B,C", "A,B,C"}).Count(&count).Error
	return count
}

func countLegacyScanRecords(db *gorm.DB) int64 {
	var count int64
	_ = db.Model(&model.ScanRecord{}).Where("station_role IN ?", []string{"A", "B", "C", "A,B", "A,C", "B,C", "A,B,C"}).Count(&count).Error
	return count
}
