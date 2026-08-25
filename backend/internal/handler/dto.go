package handler

type RegisterRequest struct {
	Username    string `json:"username" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	Password    string `json:"password" binding:"required"`
	StationRole string `json:"station_role"`
}

type LoginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SetWorkerProcessesRequest struct {
	ProcessIDs []uint `json:"process_ids"`
}

type CreateProcessRequest struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
	Sort int    `json:"sort"`
}

type CreateOrderItemRequest struct {
	ItemNo        string            `json:"item_no" binding:"required"`
	ComponentType string            `json:"component_type" binding:"required"`
	PartName      string            `json:"part_name" binding:"required"`
	Model         string            `json:"model"`
	Spec          string            `json:"spec"`
	Dimensions    map[string]string `json:"dimensions"`
	Material      string            `json:"material"`
	Quantity      float64           `json:"quantity" binding:"required"`
	Unit          string            `json:"unit"`
	Remark        string            `json:"remark"`
}

type CreateOrderRequest struct {
	CustomerName string                   `json:"customer_name"`
	ProductName  string                   `json:"product_name"`
	Spec         string                   `json:"spec"`
	Quantity     float64                  `json:"quantity"`
	DeliveryDate string                   `json:"delivery_date"`
	Status       string                   `json:"status"`
	Remark       string                   `json:"remark"`
	ProcessIDs   []uint                   `json:"process_ids"`
	Items        []CreateOrderItemRequest `json:"items"`
}

type RecordScanRequest struct {
	QRToken string `json:"qr_token" binding:"required"`
}
