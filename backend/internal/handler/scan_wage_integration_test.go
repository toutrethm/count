package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"count/backend/internal/middleware"
	"count/backend/internal/model"
	"count/backend/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestRecordScanStoresWageSnapshot(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.Process{}, &model.WorkerProcess{}, &model.Order{}, &model.OrderItem{}, &model.OrderProcess{}, &model.WageRule{}, &model.ScanRecord{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	worker := model.User{
		Username:     "测试工人",
		Phone:        "13800000001",
		PasswordHash: string(passwordHash),
		Role:         "worker",
		StationRole:  "blanking_center",
		Status:       1,
	}
	if err := db.Create(&worker).Error; err != nil {
		t.Fatalf("create worker: %v", err)
	}

	process := model.Process{
		Code:        "blanking_center",
		Name:        "切料和打中心孔",
		StationRole: "blanking_center",
		Sort:        10,
		Status:      1,
	}
	if err := db.Create(&process).Error; err != nil {
		t.Fatalf("create process: %v", err)
	}

	if err := db.Create(&model.WageRule{
		Code:          "guide_pillar_d30_l151_250",
		ComponentType: "guide_pillar",
		MinDiameter:   30,
		MaxDiameter:   30,
		MinLength:     151,
		MaxLength:     250,
		BaseUnitPrice: 0.35,
		Status:        1,
	}).Error; err != nil {
		t.Fatalf("create wage rule: %v", err)
	}

	order := model.Order{
		OrderNo:      "ORD-WAGE-1",
		QRToken:      "ORD-WAGE-1",
		CustomerName: "客户A",
		ProductName:  "产品A",
		Quantity:     12,
		ScanLimit:    1,
		ScanCount:    0,
		Status:       "draft",
	}
	if err := db.Create(&order).Error; err != nil {
		t.Fatalf("create order: %v", err)
	}

	item := model.OrderItem{
		OrderID:       order.ID,
		ItemNo:        "1",
		ComponentType: "guide_pillar",
		PartName:      "导柱",
		Quantity:      12,
		Dimensions:    `{"innerDiameter":"30","length":"220"}`,
	}
	if err := db.Create(&item).Error; err != nil {
		t.Fatalf("create order item: %v", err)
	}

	orderProcess := model.OrderProcess{
		OrderID:     order.ID,
		OrderItemID: item.ID,
		ProcessID:   process.ID,
		StationRole: process.StationRole,
		Sort:        10,
		Status:      1,
	}
	if err := db.Create(&orderProcess).Error; err != nil {
		t.Fatalf("create order process: %v", err)
	}

	h := New(db, "secret")
	engine := gin.New()
	engine.POST("/api/scans/record", middleware.Auth([]byte("secret")), h.RecordScan)
	token, err := auth.GenerateToken([]byte("secret"), worker.ID, worker.Role, worker.StationRole, worker.Phone, worker.Username)
	if err != nil {
		t.Fatalf("token: %v", err)
	}

	body, _ := json.Marshal(RecordScanRequest{QRToken: "ORD-WAGE-1"})
	req := httptest.NewRequest(http.MethodPost, "/api/scans/record", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Record model.ScanRecord `json:"record"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if response.Record.WageRuleCode != "guide_pillar_d30_l151_250" {
		t.Fatalf("unexpected wage rule code: %s", response.Record.WageRuleCode)
	}
	if response.Record.WageUnitPrice != 0.35 {
		t.Fatalf("unexpected unit price: %.2f", response.Record.WageUnitPrice)
	}
	if response.Record.WageAmount != 6.3 {
		t.Fatalf("unexpected wage amount: %.2f", response.Record.WageAmount)
	}
}
