package handler

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"count/backend/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *Handler) CreateOrder(c *gin.Context) {
	claims, ok := currentClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "missing token"})
		return
	}

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "missing items"})
		return
	}

	orderNo := buildOrderNo()

	var deliveryDate *time.Time
	if strings.TrimSpace(req.DeliveryDate) != "" {
		parsed, err := time.Parse("2006-01-02", req.DeliveryDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid delivery_date"})
			return
		}
		deliveryDate = &parsed
	}

	order := model.Order{
		OrderNo:      orderNo,
		QRToken:      orderNo,
		CustomerName: strings.TrimSpace(req.CustomerName),
		ProductName:  strings.TrimSpace(req.ProductName),
		Spec:         strings.TrimSpace(req.Spec),
		Quantity:     totalItemQuantity(req.Items, req.Quantity),
		ScanLimit:    deriveOrderScanLimit(req.Items),
		ScanCount:    0,
		DeliveryDate: deliveryDate,
		Status:       defaultOrderStatus(req.Status),
		CreatedBy:    &claims.UserID,
		Remark:       strings.TrimSpace(req.Remark),
	}

	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		items := make([]model.OrderItem, 0, len(req.Items))
		for _, item := range req.Items {
			if len(buildOrderWorkflowPlans(item)) == 0 {
				return fmt.Errorf("unsupported component_type: %s", strings.TrimSpace(item.ComponentType))
			}
			dimensionsJSON, err := marshalDimensions(item.Dimensions)
			if err != nil {
				return err
			}
			items = append(items, model.OrderItem{
				OrderID:       order.ID,
				ItemNo:        strings.TrimSpace(item.ItemNo),
				ComponentType: strings.TrimSpace(item.ComponentType),
				PartName:      strings.TrimSpace(item.PartName),
				Model:         normalizeOrderItemModel(item.Model, item.PartName),
				Spec:          strings.TrimSpace(item.Spec),
				Dimensions:    dimensionsJSON,
				Material:      strings.TrimSpace(item.Material),
				Quantity:      item.Quantity,
				Unit:          strings.TrimSpace(item.Unit),
				Remark:        strings.TrimSpace(item.Remark),
			})
		}
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}

		processByCode, err := h.processesByCode(tx)
		if err != nil {
			return err
		}
		processes := make([]model.OrderProcess, 0)
		for itemIndex, item := range req.Items {
			plans := buildOrderWorkflowPlans(item)
			for _, plan := range plans {
				process, ok := processByCode[plan.ProcessCode]
				if !ok {
					return fmt.Errorf("process code not seeded: %s", plan.ProcessCode)
				}
				processes = append(processes, model.OrderProcess{
					OrderID:     order.ID,
					OrderItemID: items[itemIndex].ID,
					ProcessID:   process.ID,
					StationRole: plan.StationRole,
					Sort:        itemIndex*100 + plan.Sort,
					Status:      1,
				})
			}
		}
		if len(processes) > 0 {
			if err := tx.Create(&processes).Error; err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if err := h.loadOrderGraph(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": order})
}

func (h *Handler) ListOrders(c *gin.Context) {
	orders := make([]model.Order, 0)
	if err := h.DB.Preload("Items").Preload("Processes.Process").Order("id desc").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": orders})
}

func (h *Handler) GetOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	claims, ok := currentClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "missing token"})
		return
	}

	order, err := h.fetchOrder(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "order not found"})
		return
	}

	if claims.Role != "admin" {
		allowed, err := h.workerCanAccessOrder(claims.UserID, order.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"message": "permission denied"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"item": order})
}

func (h *Handler) GetOrderByNo(c *gin.Context) {
	orderNo := strings.TrimSpace(c.Param("orderNo"))
	if orderNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid order_no"})
		return
	}

	claims, ok := currentClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "missing token"})
		return
	}

	order, err := h.findOrderByCode(orderNo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "order not found"})
		return
	}

	if claims.Role != "admin" {
		allowed, err := h.workerCanAccessOrder(claims.UserID, order.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"message": "permission denied"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"item": order})
}

func (h *Handler) UpdateWorkerProcesses(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	var req SetWorkerProcessesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&model.WorkerProcess{}).Error; err != nil {
			return err
		}

		if len(req.ProcessIDs) == 0 {
			return tx.Model(&model.User{}).Where("id = ?", userID).Update("station_role", "").Error
		}

		rows := make([]model.WorkerProcess, 0, len(req.ProcessIDs))
		codes := make([]string, 0, len(req.ProcessIDs))
		seen := make(map[string]struct{}, len(req.ProcessIDs))
		for _, processID := range req.ProcessIDs {
			rows = append(rows, model.WorkerProcess{
				UserID:    uint(userID),
				ProcessID: processID,
			})

			var process model.Process
			if err := tx.First(&process, processID).Error; err != nil {
				return err
			}
			for _, code := range stationRoleParts(process.StationRole) {
				if _, ok := seen[code]; ok {
					continue
				}
				seen[code] = struct{}{}
				codes = append(codes, code)
			}
		}
		if err := tx.Create(&rows).Error; err != nil {
			return err
		}
		if len(codes) == 0 {
			return nil
		}
		return tx.Model(&model.User{}).Where("id = ?", userID).Update("station_role", strings.Join(codes, ",")).Error
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	processes, _ := h.assignedProcesses(uint(userID))
	c.JSON(http.StatusOK, gin.H{"items": processes})
}

func (h *Handler) ListWorkers(c *gin.Context) {
	workers := make([]model.User, 0)
	if err := h.DB.Where("role = ?", "worker").Order("id desc").Find(&workers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": workers})
}

func (h *Handler) ListMyOrders(c *gin.Context) {
	claims, ok := currentClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "missing token"})
		return
	}

	orders := make([]model.Order, 0)
	if err := h.DB.Preload("Items").Preload("Processes.Process").Order("id desc").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	filtered := make([]model.Order, 0, len(orders))
	for _, order := range orders {
		if claims.Role == "admin" {
			filtered = append(filtered, order)
			continue
		}
		allowed, err := h.workerCanAccessOrder(claims.UserID, order.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		if allowed {
			filtered = append(filtered, order)
		}
	}

	c.JSON(http.StatusOK, gin.H{"items": filtered})
}

func buildOrderNo() string {
	return buildOrderNoAt(time.Now(), randomOrderSuffix())
}

func buildOrderNoAt(now time.Time, suffix string) string {
	return fmt.Sprintf("%s%03d-%s", now.Format("20060102150405"), now.Nanosecond()/1e6, suffix)
}

func defaultOrderStatus(status string) string {
	status = strings.TrimSpace(status)
	if status == "" {
		return "draft"
	}
	return status
}

func totalItemQuantity(items []CreateOrderItemRequest, fallback float64) float64 {
	if len(items) == 0 {
		return fallback
	}

	total := 0.0
	for _, item := range items {
		total += item.Quantity
	}
	return total
}

func randomOrderSuffix() string {
	n, err := rand.Int(rand.Reader, big.NewInt(1000))
	if err != nil {
		fallback := time.Now().UnixNano() % 1000
		if fallback < 0 {
			fallback = -fallback
		}
		return fmt.Sprintf("%03d", fallback)
	}
	return fmt.Sprintf("%03d", n.Int64())
}

func normalizeOrderItemModel(modelValue, partName string) string {
	modelValue = strings.TrimSpace(modelValue)
	if modelValue != "" {
		return modelValue
	}
	return strings.TrimSpace(partName)
}

func marshalDimensions(dimensions map[string]string) (string, error) {
	if dimensions == nil {
		dimensions = map[string]string{}
	}
	bytes, err := json.Marshal(dimensions)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
