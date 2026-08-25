package seed

import (
	"fmt"

	"count/backend/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const adminPasswordHash = "$2a$10$7aWlx1WSuW.bL/E5FyKnsehUSjfHC9bHtmcpSrPXPLjZipt32w8R6"

var legacyUserRoleByPhone = map[string]string{
	"17742523836": "blanking_center,center_hole,drill_tap_batch",
	"13959578057": "turn_outer,turn_sleeve_manual",
	"18250505219": "turn_outer,turn_sleeve_auto",
	"15106017652": "turn_head,center_hole,drill_tap_small,drill_tap_batch,turn_sleeve_auto",
	"15501148077": "turn_sleeve_manual",
	"18039005017": "turn_sleeve_manual",
	"13950776176": "turn_sleeve_manual",
}

func Defaults(db *gorm.DB) error {
	if err := seedWageRules(db); err != nil {
		return err
	}

	processes := []model.Process{
		{Code: "blanking_center", Name: "切料和打中心孔", StationRole: "blanking_center", Sort: 10, Status: 1},
		{Code: "turn_outer", Name: "车外圆", StationRole: "turn_outer", Sort: 20, Status: 1},
		{Code: "turn_head", Name: "车大头", StationRole: "turn_head", Sort: 30, Status: 1},
		{Code: "center_hole", Name: "打中心孔", StationRole: "center_hole", Sort: 40, Status: 1},
		{Code: "turn_head_center", Name: "车大头和打中心孔", StationRole: "turn_head_center", Sort: 50, Status: 1},
		{Code: "drill_tap_small", Name: "钻孔攻牙", StationRole: "drill_tap_small", Sort: 60, Status: 1},
		{Code: "drill_tap_batch", Name: "批量钻孔攻牙", StationRole: "drill_tap_batch", Sort: 70, Status: 1},
		{Code: "turn_sleeve", Name: "车套", StationRole: "turn_sleeve_auto,turn_sleeve_manual", Sort: 80, Status: 1},
	}

	for _, process := range processes {
		if err := db.Where("code = ?", process.Code).Assign(process).FirstOrCreate(&process).Error; err != nil {
			return err
		}
	}

	admin := model.User{
		Username:     "管理员",
		Phone:        "13800000000",
		PasswordHash: adminPasswordHash,
		Role:         "admin",
		StationRole:  "",
		Status:       1,
	}
	if err := db.Where("phone = ?", admin.Phone).Assign(admin).FirstOrCreate(&admin).Error; err != nil {
		return err
	}

	workers := []struct {
		Username    string
		Phone       string
		StationRole string
	}{
		{Username: "赵吴洋", Phone: "17742523836", StationRole: "blanking_center,center_hole,drill_tap_batch"},
		{Username: "刘小印", Phone: "13959578057", StationRole: "turn_outer,turn_sleeve_manual"},
		{Username: "王海龙", Phone: "18250505219", StationRole: "turn_outer,turn_sleeve_auto"},
		{Username: "陈家兴", Phone: "15106017652", StationRole: "turn_head,center_hole,drill_tap_small,drill_tap_batch,turn_sleeve_auto"},
		{Username: "郑志权", Phone: "15501148077", StationRole: "turn_sleeve_manual"},
		{Username: "郑志勇", Phone: "18039005017", StationRole: "turn_sleeve_manual"},
		{Username: "吴丽先", Phone: "13950776176", StationRole: "turn_sleeve_manual"},
	}

	for _, item := range workers {
		var user model.User
		password := fmt.Sprintf("%s8888", last4(item.Phone))
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user = model.User{
			Username:     item.Username,
			Phone:        item.Phone,
			PasswordHash: string(passwordHash),
			Role:         "worker",
			StationRole:  item.StationRole,
			Status:       1,
		}
		if err := db.Where("phone = ?", item.Phone).Assign(user).FirstOrCreate(&user).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedWageRules(db *gorm.DB) error {
	for _, rule := range wageRuleCatalog() {
		if err := db.Where("code = ?", rule.Code).Assign(rule).FirstOrCreate(&rule).Error; err != nil {
			return err
		}
	}
	return nil
}

func BackfillLegacyRoles(db *gorm.DB) error {
	if err := backfillLegacyUsers(db); err != nil {
		return err
	}
	if err := backfillLegacyProcesses(db); err != nil {
		return err
	}
	if err := backfillLegacyOrderProcesses(db); err != nil {
		return err
	}
	if err := backfillLegacyScanRecords(db); err != nil {
		return err
	}
	return nil
}

func backfillLegacyUsers(db *gorm.DB) error {
	for phone, capability := range legacyUserRoleByPhone {
		if err := db.Model(&model.User{}).Where("phone = ?", phone).Update("station_role", capability).Error; err != nil {
			return err
		}
	}

	type legacyUserUpdate struct {
		ID          uint
		StationRole string
	}
	rows := make([]legacyUserUpdate, 0)
	if err := db.Model(&model.User{}).
		Select("id, station_role").
		Where("station_role IN ?", []string{"A", "B", "C", "A,B", "A,C", "B,C", "A,B,C"}).
		Scan(&rows).Error; err != nil {
		return err
	}

	for _, row := range rows {
		if err := db.Model(&model.User{}).Where("id = ?", row.ID).Update("station_role", normalizeLegacyCapabilities(row.StationRole)).Error; err != nil {
			return err
		}
	}
	return nil
}

func backfillLegacyProcesses(db *gorm.DB) error {
	return db.Exec(`
		UPDATE processes
		SET station_role = CASE code
			WHEN 'blanking_center' THEN 'blanking_center'
			WHEN 'turn_outer' THEN 'turn_outer'
			WHEN 'turn_head' THEN 'turn_head'
			WHEN 'center_hole' THEN 'center_hole'
			WHEN 'turn_head_center' THEN 'turn_head_center'
			WHEN 'turn_sleeve' THEN 'turn_sleeve_auto,turn_sleeve_manual'
			WHEN 'drill_tap' THEN 'drill_tap_small,drill_tap_batch'
			ELSE station_role
		END
		WHERE code IN ('blanking_center', 'turn_outer', 'turn_head', 'center_hole', 'turn_head_center', 'turn_sleeve', 'drill_tap')
	`).Error
}

func backfillLegacyOrderProcesses(db *gorm.DB) error {
	return db.Exec(`
		UPDATE order_processes op
		JOIN processes p ON p.id = op.process_id
		LEFT JOIN order_items oi ON oi.id = op.order_item_id
		SET op.station_role = CASE
			WHEN p.code = 'blanking_center' THEN 'blanking_center'
			WHEN p.code = 'turn_outer' THEN 'turn_outer'
			WHEN p.code = 'turn_head' THEN 'turn_head'
			WHEN p.code = 'center_hole' THEN 'center_hole'
			WHEN p.code = 'turn_head_center' THEN 'turn_head_center'
			WHEN p.code = 'turn_sleeve' THEN 'turn_sleeve_auto,turn_sleeve_manual'
			WHEN p.code = 'drill_tap' AND COALESCE(oi.quantity, 0) < 20 THEN 'drill_tap_small'
			WHEN p.code = 'drill_tap' THEN 'drill_tap_batch'
			ELSE op.station_role
		END
		WHERE op.station_role IN ('A', 'B', 'C', 'A,B', 'A,C', 'B,C', 'A,B,C')
	`).Error
}

func backfillLegacyScanRecords(db *gorm.DB) error {
	return db.Exec(`
		UPDATE scan_records sr
		JOIN processes p ON p.id = sr.process_id
		LEFT JOIN order_items oi ON oi.id = sr.order_item_id
		SET sr.station_role = CASE
			WHEN p.code = 'blanking_center' THEN 'blanking_center'
			WHEN p.code = 'turn_outer' THEN 'turn_outer'
			WHEN p.code = 'turn_head' THEN 'turn_head'
			WHEN p.code = 'center_hole' THEN 'center_hole'
			WHEN p.code = 'turn_head_center' THEN 'turn_head_center'
			WHEN p.code = 'turn_sleeve' THEN 'turn_sleeve_auto,turn_sleeve_manual'
			WHEN p.code = 'drill_tap' AND COALESCE(oi.quantity, 0) < 20 THEN 'drill_tap_small'
			WHEN p.code = 'drill_tap' THEN 'drill_tap_batch'
			ELSE sr.station_role
		END
		WHERE sr.station_role IN ('A', 'B', 'C', 'A,B', 'A,C', 'B,C', 'A,B,C')
	`).Error
}

func normalizeLegacyCapabilities(role string) string {
	switch role {
	case "A":
		return "blanking_center,center_hole,drill_tap_batch"
	case "B":
		return "turn_outer,turn_sleeve_manual"
	case "C":
		return "turn_head,center_hole,drill_tap_small,turn_sleeve_auto"
	case "A,B":
		return "blanking_center,center_hole,drill_tap_batch,turn_outer,turn_sleeve_manual"
	case "A,C":
		return "blanking_center,center_hole,drill_tap_batch,turn_head,drill_tap_small,turn_sleeve_auto"
	case "B,C":
		return "turn_outer,turn_sleeve_manual,turn_head,center_hole,drill_tap_small,turn_sleeve_auto"
	case "A,B,C":
		return "blanking_center,center_hole,drill_tap_batch,turn_outer,turn_sleeve_manual,turn_head,drill_tap_small,turn_sleeve_auto"
	default:
		return role
	}
}

func last4(phone string) string {
	if len(phone) <= 4 {
		return phone
	}
	return phone[len(phone)-4:]
}

func wageRuleCatalog() []model.WageRule {
	return []model.WageRule{
		{Code: "guide_pillar_d0_25_l0_150", ComponentType: "guide_pillar", MinDiameter: 0, MaxDiameter: 25, MinLength: 0, MaxLength: 150, BaseUnitPrice: 0.2, Status: 1},
		{Code: "guide_pillar_d0_25_l151_250", ComponentType: "guide_pillar", MinDiameter: 0, MaxDiameter: 25, MinLength: 151, MaxLength: 250, BaseUnitPrice: 0.28, Status: 1},
		{Code: "guide_pillar_d0_25_l251_350", ComponentType: "guide_pillar", MinDiameter: 0, MaxDiameter: 25, MinLength: 251, MaxLength: 350, BaseUnitPrice: 0.35, Status: 1},
		{Code: "guide_pillar_d30_l0_150", ComponentType: "guide_pillar", MinDiameter: 30, MaxDiameter: 30, MinLength: 0, MaxLength: 150, BaseUnitPrice: 0.24, Status: 1},
		{Code: "guide_pillar_d30_l151_250", ComponentType: "guide_pillar", MinDiameter: 30, MaxDiameter: 30, MinLength: 151, MaxLength: 250, BaseUnitPrice: 0.35, Status: 1},
		{Code: "guide_pillar_d30_l251_350", ComponentType: "guide_pillar", MinDiameter: 30, MaxDiameter: 30, MinLength: 251, MaxLength: 350, BaseUnitPrice: 0.4, Status: 1},
		{Code: "guide_pillar_d35_l0_150", ComponentType: "guide_pillar", MinDiameter: 35, MaxDiameter: 35, MinLength: 0, MaxLength: 150, BaseUnitPrice: 0.3, Status: 1},
		{Code: "guide_pillar_d35_l151_250", ComponentType: "guide_pillar", MinDiameter: 35, MaxDiameter: 35, MinLength: 151, MaxLength: 250, BaseUnitPrice: 0.38, Status: 1},
		{Code: "guide_pillar_d35_l251_350", ComponentType: "guide_pillar", MinDiameter: 35, MaxDiameter: 35, MinLength: 251, MaxLength: 350, BaseUnitPrice: 0.45, Status: 1},
		{Code: "guide_pillar_d40_l0_150", ComponentType: "guide_pillar", MinDiameter: 40, MaxDiameter: 40, MinLength: 0, MaxLength: 150, BaseUnitPrice: 0.36, Status: 1},
		{Code: "guide_pillar_d40_l151_250", ComponentType: "guide_pillar", MinDiameter: 40, MaxDiameter: 40, MinLength: 151, MaxLength: 250, BaseUnitPrice: 0.42, Status: 1},
		{Code: "guide_pillar_d40_l251_350", ComponentType: "guide_pillar", MinDiameter: 40, MaxDiameter: 40, MinLength: 251, MaxLength: 350, BaseUnitPrice: 0.53, Status: 1},
		{Code: "top_pin_d0_30_l0_150", ComponentType: "top_pin", MinDiameter: 0, MaxDiameter: 30, MinLength: 0, MaxLength: 150, BaseUnitPrice: 0.2, Status: 1},
		{Code: "top_pin_d0_30_l151_250", ComponentType: "top_pin", MinDiameter: 0, MaxDiameter: 30, MinLength: 151, MaxLength: 250, BaseUnitPrice: 0.26, Status: 1},
		{Code: "top_pin_d0_30_l251_350", ComponentType: "top_pin", MinDiameter: 0, MaxDiameter: 30, MinLength: 251, MaxLength: 350, BaseUnitPrice: 0.32, Status: 1},
		{Code: "top_pin_d35_l0_150", ComponentType: "top_pin", MinDiameter: 35, MaxDiameter: 35, MinLength: 0, MaxLength: 150, BaseUnitPrice: 0.26, Status: 1},
		{Code: "top_pin_d35_l151_250", ComponentType: "top_pin", MinDiameter: 35, MaxDiameter: 35, MinLength: 151, MaxLength: 250, BaseUnitPrice: 0.32, Status: 1},
		{Code: "top_pin_d35_l251_350", ComponentType: "top_pin", MinDiameter: 35, MaxDiameter: 35, MinLength: 251, MaxLength: 350, BaseUnitPrice: 0.38, Status: 1},
	}
}
