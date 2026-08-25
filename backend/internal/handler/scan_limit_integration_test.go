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

func TestRecordScanRejectsSecondSleeveScanWhenLimitReached(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.Process{}, &model.WorkerProcess{}, &model.Order{}, &model.OrderItem{}, &model.OrderProcess{}, &model.ScanRecord{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	worker := model.User{
		Username:     "测试工人",
		Phone:        "13800000000",
		PasswordHash: string(passwordHash),
		Role:         "worker",
		StationRole:  "turn_sleeve_auto",
		Status:       1,
	}
	if err := db.Create(&worker).Error; err != nil {
		t.Fatalf("create worker: %v", err)
	}

	process := model.Process{
		Code:        "turn_sleeve",
		Name:        "车套",
		StationRole: "turn_sleeve_auto",
		Sort:        10,
		Status:      1,
	}
	if err := db.Create(&process).Error; err != nil {
		t.Fatalf("create process: %v", err)
	}

	order := model.Order{
		OrderNo:      "ORD-TEST-1",
		QRToken:      "ORD-TEST-1",
		CustomerName: "客户A",
		ProductName:  "导套",
		Quantity:     1,
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
		ComponentType: "guide_bush",
		PartName:      "导套",
		Quantity:      1,
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

	call := func() *httptest.ResponseRecorder {
		body, _ := json.Marshal(RecordScanRequest{QRToken: "ORD-TEST-1"})
		req := httptest.NewRequest(http.MethodPost, "/api/scans/record", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		recorder := httptest.NewRecorder()
		engine.ServeHTTP(recorder, req)
		return recorder
	}

	first := call()
	if first.Code != http.StatusOK {
		t.Fatalf("expected first scan 200, got %d: %s", first.Code, first.Body.String())
	}

	second := call()
	if second.Code != http.StatusConflict {
		t.Fatalf("expected second scan 409, got %d: %s", second.Code, second.Body.String())
	}
}
