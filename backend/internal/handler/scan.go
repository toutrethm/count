package handler

import (
	"errors"
	"net/http"
	"time"

	"count/backend/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var errOrderScanLimitReached = errors.New("order scan limit reached")

func (h *Handler) PreviewScanOrder(c *gin.Context) {
	claims, ok := currentClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "missing token"})
		return
	}

	qrToken := c.Param("qrToken")
	order, err := h.findOrderByCode(qrToken)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "order not found"})
		return
	}

	orderWithGraph, err := h.fetchOrder(order.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	userRoles, err := h.userStationRolesOnDB(h.DB, claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	matchingProcesses := matchingOrderProcessesForRoles(orderWithGraph.Processes, userRoles)
	if claims.Role != "admin" && len(matchingProcesses) == 0 {
		c.JSON(http.StatusForbidden, gin.H{"message": "permission denied"})
		return
	}

	pendingStep, err := h.workerCanScanProcessOnDB(h.DB, claims.UserID, order.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var pendingProcess *model.OrderProcess
	if pendingStep != nil {
		pendingProcess = findOrderProcessByID(orderWithGraph.Processes, pendingStep.ID)
	}

	c.JSON(http.StatusOK, gin.H{
		"order":              orderWithGraph,
		"matching_processes": matchingProcesses,
		"pending_process":    pendingProcess,
		"can_record":         pendingProcess != nil,
	})
}

func (h *Handler) RecordScan(c *gin.Context) {
	claims, ok := currentClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "missing token"})
		return
	}

	var req RecordScanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var order model.Order
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		lockedOrder, err := h.findOrderByCodeOnDB(tx.Clauses(clause.Locking{Strength: "UPDATE"}), req.QRToken)
		if err != nil {
			return err
		}

		if lockedOrder.ScanCount >= orderScanLimitForOrder(lockedOrder) {
			return errOrderScanLimitReached
		}

		order = lockedOrder
		return nil
	}); err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "order not found"})
		case errors.Is(err, errOrderScanLimitReached):
			c.JSON(http.StatusConflict, gin.H{"message": "order scan limit reached"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}

	step, err := h.workerCanScanProcessOnDB(h.DB, claims.UserID, order.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if step == nil && claims.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"message": "permission denied"})
		return
	}
	if step == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "no pending process"})
		return
	}

	userRoles, err := h.userStationRolesOnDB(h.DB, claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	scanRole := matchStationRole(step.StationRole, userRoles)
	if scanRole == "" {
		scanRole = step.StationRole
	}

	record := model.ScanRecord{
		OrderID:        order.ID,
		OrderItemID:    step.OrderItemID,
		OrderProcessID: step.ID,
		ProcessID:      step.ProcessID,
		UserID:         claims.UserID,
		StationRole:    scanRole,
		ScannedAt:      time.Now().Truncate(time.Second),
		Source:         "scan",
	}
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		var currentOrder model.Order
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Preload("Items").First(&currentOrder, order.ID).Error; err != nil {
			return err
		}

		if currentOrder.ScanCount >= orderScanLimitForOrder(currentOrder) {
			return errOrderScanLimitReached
		}

		orderItem := findOrderItemByID(currentOrder.Items, step.OrderItemID)
		if orderItem == nil {
			return gorm.ErrRecordNotFound
		}
		wageRule, wageUnitPrice, _, wageAmount, ok, err := calculateOrderItemWageOnDB(tx, CreateOrderItemRequest{
			ComponentType: orderItem.ComponentType,
			Quantity:      orderItem.Quantity,
			Dimensions:    decodeDimensions(orderItem.Dimensions),
		})
		if err != nil {
			return err
		}
		if ok {
			record.WageRuleID = &wageRule.ID
			record.WageRuleCode = wageRule.Code
			record.WageUnitPrice = wageUnitPrice
			record.WageAmount = wageAmount
		}

		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		return tx.Model(&model.Order{}).Where("id = ?", currentOrder.ID).UpdateColumn("scan_count", gorm.Expr("scan_count + ?", 1)).Error
	}); err != nil {
		switch {
		case errors.Is(err, errOrderScanLimitReached):
			c.JSON(http.StatusConflict, gin.H{"message": "order scan limit reached"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}

	orderWithGraph, err := h.fetchOrder(order.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"record": record,
		"order":  orderWithGraph,
	})
}

func (h *Handler) ListMyScanRecords(c *gin.Context) {
	claims, ok := currentClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "missing token"})
		return
	}

	query := h.DB.Preload("Order").Preload("OrderItem").Preload("OrderProcess").Preload("Process").Preload("User").Model(&model.ScanRecord{}).Where("user_id = ?", claims.UserID)
	if claims.Role == "admin" {
		query = h.DB.Preload("Order").Preload("OrderItem").Preload("OrderProcess").Preload("Process").Preload("User").Model(&model.ScanRecord{})
	}

	records := make([]model.ScanRecord, 0)
	if err := query.Order("scanned_at desc").Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": records})
}

func (h *Handler) ListAllScanRecords(c *gin.Context) {
	records := make([]model.ScanRecord, 0)
	if err := h.DB.Preload("Order").Preload("OrderItem").Preload("OrderProcess").Preload("Process").Preload("User").Order("scanned_at desc").Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": records})
}

func findOrderItemByID(items []model.OrderItem, id uint) *model.OrderItem {
	for i := range items {
		if items[i].ID == id {
			return &items[i]
		}
	}
	return nil
}

func findOrderProcessByID(processes []model.OrderProcess, id uint) *model.OrderProcess {
	for i := range processes {
		if processes[i].ID == id {
			return &processes[i]
		}
	}
	return nil
}

func matchingOrderProcessesForRoles(processes []model.OrderProcess, roles []string) []model.OrderProcess {
	result := make([]model.OrderProcess, 0)
	for _, process := range processes {
		if matchStationRole(process.StationRole, roles) != "" {
			result = append(result, process)
		}
	}
	return result
}
