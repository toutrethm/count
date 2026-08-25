package model

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"size:64;not null" json:"username"`
	Phone        string    `gorm:"size:32;not null;uniqueIndex" json:"phone"`
	PasswordHash string    `gorm:"size:100;not null" json:"-"`
	Role         string    `gorm:"size:16;not null;default:worker;index" json:"role"`
	StationRole  string    `gorm:"size:255;not null;default:'';index" json:"station_role"`
	Status       int8      `gorm:"not null;default:1" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Process struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Code        string    `gorm:"size:32;not null;uniqueIndex" json:"code"`
	Name        string    `gorm:"size:64;not null" json:"name"`
	StationRole string    `gorm:"size:255;not null;default:'';index" json:"station_role"`
	Sort        int       `gorm:"not null;default:0" json:"sort"`
	Status      int8      `gorm:"not null;default:1" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type WorkerProcess struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex:uk_worker_processes" json:"user_id"`
	ProcessID uint      `gorm:"not null;uniqueIndex:uk_worker_processes" json:"process_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Order struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	OrderNo      string         `gorm:"size:64;not null;uniqueIndex" json:"order_no"`
	QRToken      string         `gorm:"size:64;not null;uniqueIndex" json:"qr_token"`
	CustomerName string         `gorm:"size:64;not null" json:"customer_name"`
	ProductName  string         `gorm:"size:128;not null" json:"product_name"`
	Spec         string         `gorm:"size:128" json:"spec"`
	Quantity     float64        `gorm:"type:decimal(18,3);not null;default:0" json:"quantity"`
	ScanLimit    int            `gorm:"not null;default:1" json:"scan_limit"`
	ScanCount    int            `gorm:"not null;default:0" json:"scan_count"`
	DeliveryDate *time.Time     `gorm:"type:date" json:"delivery_date"`
	Status       string         `gorm:"size:32;not null;default:draft;index" json:"status"`
	CreatedBy    *uint          `gorm:"index" json:"created_by"`
	Remark       string         `gorm:"size:255" json:"remark"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	Items        []OrderItem    `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	Processes    []OrderProcess `gorm:"foreignKey:OrderID" json:"processes,omitempty"`
}

type OrderItem struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	OrderID       uint      `gorm:"not null;index" json:"order_id"`
	ItemNo        string    `gorm:"size:64;not null" json:"item_no"`
	ComponentType string    `gorm:"size:64;not null;index" json:"component_type"`
	PartName      string    `gorm:"size:128;not null" json:"part_name"`
	Model         string    `gorm:"size:128" json:"model"`
	Spec          string    `gorm:"size:128" json:"spec"`
	Dimensions    string    `gorm:"type:json" json:"dimensions"`
	Material      string    `gorm:"size:128" json:"material"`
	Quantity      float64   `gorm:"type:decimal(18,3);not null;default:0" json:"quantity"`
	Unit          string    `gorm:"size:16" json:"unit"`
	Remark        string    `gorm:"size:255" json:"remark"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type OrderProcess struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	OrderID     uint      `gorm:"not null;uniqueIndex:uk_order_processes" json:"order_id"`
	OrderItemID uint      `gorm:"not null;uniqueIndex:uk_order_processes" json:"order_item_id"`
	ProcessID   uint      `gorm:"not null;uniqueIndex:uk_order_processes" json:"process_id"`
	StationRole string    `gorm:"size:255;not null;index" json:"station_role"`
	Sort        int       `gorm:"not null;default:0" json:"sort"`
	Status      int8      `gorm:"not null;default:1" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Process     Process   `gorm:"foreignKey:ProcessID" json:"process,omitempty"`
	OrderItem   OrderItem `gorm:"foreignKey:OrderItemID" json:"order_item,omitempty"`
}

type WageRule struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Code          string    `gorm:"size:64;not null;uniqueIndex" json:"code"`
	ComponentType string    `gorm:"size:64;not null;index" json:"component_type"`
	MinDiameter   float64   `gorm:"type:decimal(10,2);not null;default:0" json:"min_diameter"`
	MaxDiameter   float64   `gorm:"type:decimal(10,2);not null;default:0" json:"max_diameter"`
	MinLength     float64   `gorm:"type:decimal(10,2);not null;default:0" json:"min_length"`
	MaxLength     float64   `gorm:"type:decimal(10,2);not null;default:0" json:"max_length"`
	BaseUnitPrice float64   `gorm:"type:decimal(10,4);not null;default:0" json:"base_unit_price"`
	Status        int8      `gorm:"not null;default:1" json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ScanRecord struct {
	ID             uint         `gorm:"primaryKey" json:"id"`
	OrderID        uint         `gorm:"not null;index" json:"order_id"`
	OrderItemID    uint         `gorm:"not null;index" json:"order_item_id"`
	OrderProcessID uint         `gorm:"not null;uniqueIndex" json:"order_process_id"`
	ProcessID      uint         `gorm:"not null;index" json:"process_id"`
	UserID         uint         `gorm:"not null;index" json:"user_id"`
	StationRole    string       `gorm:"size:255;not null;index" json:"station_role"`
	WageRuleID     *uint        `gorm:"index" json:"wage_rule_id,omitempty"`
	WageRuleCode   string       `gorm:"size:64;not null;default:'';index" json:"wage_rule_code"`
	WageUnitPrice  float64      `gorm:"type:decimal(10,4);not null;default:0" json:"wage_unit_price"`
	WageAmount     float64      `gorm:"type:decimal(10,4);not null;default:0" json:"wage_amount"`
	ScannedAt      time.Time    `gorm:"type:datetime(0);not null;index" json:"scanned_at"`
	Source         string       `gorm:"size:32;not null;default:scan" json:"source"`
	CreatedAt      time.Time    `json:"created_at"`
	Order          Order        `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	OrderItem      OrderItem    `gorm:"foreignKey:OrderItemID" json:"order_item,omitempty"`
	OrderProcess   OrderProcess `gorm:"foreignKey:OrderProcessID" json:"order_process,omitempty"`
	Process        Process      `gorm:"foreignKey:ProcessID" json:"process,omitempty"`
	User           User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
