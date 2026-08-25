package handler

import (
	"math"
	"strings"

	"count/backend/internal/model"

	"gorm.io/gorm"
)

func calculateOrderItemWage(item CreateOrderItemRequest) (model.WageRule, float64, float64, float64, bool) {
	rule, ok := findWageRuleCandidate(item)
	if !ok {
		return model.WageRule{}, 0, 0, 0, false
	}

	unitPrice := roundMoney(rule.BaseUnitPrice)
	multiplier := quantityWageMultiplier(item.Quantity)
	amount := roundMoney(unitPrice * multiplier * item.Quantity)

	rule.BaseUnitPrice = unitPrice
	return rule, unitPrice, multiplier, amount, true
}

func calculateOrderItemWageOnDB(db *gorm.DB, item CreateOrderItemRequest) (model.WageRule, float64, float64, float64, bool, error) {
	rule, unitPrice, multiplier, amount, ok := calculateOrderItemWage(item)
	if !ok {
		return model.WageRule{}, 0, 0, 0, false, nil
	}

	var dbRule model.WageRule
	if err := db.Where("code = ?", rule.Code).First(&dbRule).Error; err != nil {
		return model.WageRule{}, 0, 0, 0, false, err
	}

	unitPrice = roundMoney(dbRule.BaseUnitPrice)
	amount = roundMoney(unitPrice * multiplier * item.Quantity)
	dbRule.BaseUnitPrice = unitPrice
	return dbRule, unitPrice, multiplier, amount, true, nil
}

func findWageRuleCandidate(item CreateOrderItemRequest) (model.WageRule, bool) {
	componentType := strings.TrimSpace(item.ComponentType)
	diameter := dimensionNumber(item.Dimensions, "innerDiameter")
	length := dimensionNumber(item.Dimensions, "length")

	for _, rule := range wageRuleCatalog() {
		if rule.ComponentType != componentType {
			continue
		}
		if diameter < rule.MinDiameter || diameter > rule.MaxDiameter {
			continue
		}
		if length < rule.MinLength || length > rule.MaxLength {
			continue
		}
		return rule, true
	}
	return model.WageRule{}, false
}

func quantityWageMultiplier(quantity float64) float64 {
	switch {
	case quantity <= 8:
		return 2
	case quantity <= 20:
		return 1.5
	case quantity <= 40:
		return 1.3
	default:
		return 1
	}
}

func roundMoney(value float64) float64 {
	return math.Round(value*100) / 100
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
