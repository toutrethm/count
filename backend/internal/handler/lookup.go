package handler

import (
	"encoding/json"
	"strings"

	"count/backend/internal/model"

	"gorm.io/gorm"
)

func (h *Handler) assignedProcesses(userID uint) ([]model.Process, error) {
	userRoles, err := h.userStationRoles(userID)
	if err != nil {
		return nil, err
	}

	processes := make([]model.Process, 0)
	if err := h.DB.Where("status = 1").Order("sort ASC, id ASC").Find(&processes).Error; err != nil {
		return nil, err
	}

	filtered := make([]model.Process, 0, len(processes))
	for _, process := range processes {
		matched := matchStationRole(process.StationRole, userRoles)
		if matched != "" {
			next := process
			next.StationRole = matched
			filtered = append(filtered, next)
		}
	}
	return filtered, nil
}

func (h *Handler) allProcesses() ([]model.Process, error) {
	var processes []model.Process
	err := h.DB.Where("status = 1").Order("sort ASC, id ASC").Find(&processes).Error
	return processes, err
}

func (h *Handler) fetchOrder(orderID uint) (model.Order, error) {
	var order model.Order
	err := h.DB.Preload("Items").Preload("Processes.Process").Preload("Processes.OrderItem").First(&order, orderID).Error
	return order, err
}

func (h *Handler) loadOrderGraph(order *model.Order) error {
	return h.DB.Preload("Items").Preload("Processes.Process").Preload("Processes.OrderItem").First(order, order.ID).Error
}

func (h *Handler) findOrderByCode(code string) (model.Order, error) {
	return h.findOrderByCodeOnDB(h.DB, code)
}

func (h *Handler) findOrderByCodeOnDB(db *gorm.DB, code string) (model.Order, error) {
	var order model.Order
	err := db.Preload("Items").Where("order_no = ? OR qr_token = ?", code, code).First(&order).Error
	return order, err
}

func (h *Handler) orderHasProcess(orderID, processID uint) (bool, error) {
	var count int64
	err := h.DB.Model(&model.OrderProcess{}).
		Where("order_id = ? AND process_id = ?", orderID, processID).
		Count(&count).Error
	return count > 0, err
}

func (h *Handler) workerCanAccessOrder(userID, orderID uint) (bool, error) {
	userRoles, err := h.userStationRoles(userID)
	if err != nil {
		return false, err
	}

	var steps []model.OrderProcess
	if err := h.DB.Where("order_id = ?", orderID).Order("sort ASC, id ASC").Find(&steps).Error; err != nil {
		return false, err
	}
	for _, step := range steps {
		if matchStationRole(step.StationRole, userRoles) != "" {
			return true, nil
		}
	}
	return false, nil
}

func (h *Handler) workerCanScanProcess(userID, orderID uint) (*model.OrderProcess, error) {
	return h.workerCanScanProcessOnDB(h.DB, userID, orderID)
}

func (h *Handler) workerCanScanProcessOnDB(db *gorm.DB, userID, orderID uint) (*model.OrderProcess, error) {
	userRoles, err := h.userStationRoles(userID)
	if err != nil {
		return nil, err
	}

	var steps []model.OrderProcess
	if err := db.Where("order_id = ?", orderID).Order("sort ASC, id ASC").Find(&steps).Error; err != nil {
		return nil, err
	}

	scanned, err := h.orderProcessScannedIDsOnDB(db, orderID)
	if err != nil {
		return nil, err
	}

	step := nextPendingOrderProcess(steps, scanned, userRoles)
	return step, nil
}

func (h *Handler) processesByCode(tx *gorm.DB) (map[string]model.Process, error) {
	processes := make([]model.Process, 0)
	if err := tx.Where("status = 1").Find(&processes).Error; err != nil {
		return nil, err
	}

	result := make(map[string]model.Process, len(processes))
	for _, process := range processes {
		result[process.Code] = process
	}
	return result, nil
}

func (h *Handler) orderProcessScannedIDs(orderID uint) (map[uint]struct{}, error) {
	return h.orderProcessScannedIDsOnDB(h.DB, orderID)
}

func (h *Handler) orderProcessScannedIDsOnDB(db *gorm.DB, orderID uint) (map[uint]struct{}, error) {
	rows := make([]struct {
		OrderProcessID uint
	}, 0)
	if err := db.Model(&model.ScanRecord{}).Select("order_process_id").Where("order_id = ?", orderID).Scan(&rows).Error; err != nil {
		return nil, err
	}

	result := make(map[uint]struct{}, len(rows))
	for _, row := range rows {
		result[row.OrderProcessID] = struct{}{}
	}
	return result, nil
}

func decodeDimensions(raw string) map[string]string {
	if strings.TrimSpace(raw) == "" {
		return map[string]string{}
	}
	var dimensions map[string]string
	if err := json.Unmarshal([]byte(raw), &dimensions); err != nil {
		return map[string]string{}
	}
	return dimensions
}

func (h *Handler) userStationRoles(userID uint) ([]string, error) {
	return h.userStationRolesOnDB(h.DB, userID)
}

func (h *Handler) userStationRolesOnDB(db *gorm.DB, userID uint) ([]string, error) {
	var user model.User
	if err := db.Select("station_role").First(&user, userID).Error; err != nil {
		return nil, err
	}
	return stationRoleParts(user.StationRole), nil
}

func matchStationRole(required string, allowed []string) string {
	requiredRoles := stationRoleParts(required)
	if len(requiredRoles) == 0 {
		return ""
	}
	requiredSet := make(map[string]struct{}, len(requiredRoles))
	for _, role := range requiredRoles {
		requiredSet[role] = struct{}{}
	}
	for _, role := range allowed {
		normalized := normalizeStationRole(role)
		if normalized == "" {
			continue
		}
		if _, ok := requiredSet[normalized]; ok {
			return normalized
		}
	}
	return ""
}
