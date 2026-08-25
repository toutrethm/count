package handler

import (
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

func TestPreviewScanOrderReturnsWorkerPendingProcess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open("file:scan_preview_order?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.Process{}, &model.Order{}, &model.OrderItem{}, &model.OrderProcess{}, &model.ScanRecord{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	worker := model.User{
		Username:     "测试工人",
		Phone:        "13800000002",
		PasswordHash: string(passwordHash),
		Role:         "worker",
		StationRole:  "blanking_center",
		Status:       1,
	}
	if err := db.Create(&worker).Error; err != nil {
		t.Fatalf("create worker: %v", err)
	}

	process := model.Process{Code: "blanking_center", Name: "切料和打中心孔", StationRole: "blanking_center", Sort: 10, Status: 1}
	if err := db.Create(&process).Error; err != nil {
		t.Fatalf("create process: %v", err)
	}

	order := model.Order{OrderNo: "ORD-PREVIEW-1", QRToken: "ORD-PREVIEW-1", Quantity: 12, ScanLimit: 3, Status: "draft"}
	if err := db.Create(&order).Error; err != nil {
		t.Fatalf("create order: %v", err)
	}

	item := model.OrderItem{OrderID: order.ID, ItemNo: "1", ComponentType: "guide_pillar", PartName: "导柱", Spec: "30*220", Quantity: 12}
	if err := db.Create(&item).Error; err != nil {
		t.Fatalf("create order item: %v", err)
	}

	orderProcess := model.OrderProcess{
		OrderID:     order.ID,
		OrderItemID: item.ID,
		ProcessID:   process.ID,
		StationRole: "blanking_center",
		Sort:        10,
		Status:      1,
	}
	if err := db.Create(&orderProcess).Error; err != nil {
		t.Fatalf("create order process: %v", err)
	}

	h := New(db, "secret")
	engine := gin.New()
	engine.GET("/api/scans/preview/:qrToken", middleware.Auth([]byte("secret")), h.PreviewScanOrder)
	token, err := auth.GenerateToken([]byte("secret"), worker.ID, worker.Role, worker.StationRole, worker.Phone, worker.Username)
	if err != nil {
		t.Fatalf("token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/scans/preview/ORD-PREVIEW-1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		CanRecord      bool               `json:"can_record"`
		PendingProcess model.OrderProcess `json:"pending_process"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if !response.CanRecord {
		t.Fatal("expected can_record to be true")
	}
	if response.PendingProcess.ID != orderProcess.ID {
		t.Fatalf("unexpected pending process id: %d", response.PendingProcess.ID)
	}
}
