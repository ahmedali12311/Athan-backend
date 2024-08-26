package tlync

import "time"

const (
	Bearer = "Bearer "
)

type Settings struct {
	Endpoint string `json:"endpoint"`
	Token    string `json:"token"`
	StoreID  string `json:"store_id"`
	FrontURL string `json:"front_url"`
}

// STEP 1: initiate

type InitiateInput struct {
	ID          string `json:"id"`
	Amount      string `json:"amount"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	BackendURL  string `json:"backend_url"`
	FrontendURL string `json:"frontend_url"`
	CustomRef   string `json:"custom_ref"`
}

type InitiateResponse struct {
	Result    string  `json:"result"`
	CustomRef string  `json:"custom_ref"`
	URL       string  `json:"url"`
	Amount    float64 `json:"amount"`
	Message   string  `json:"message,omitempty"`
}

// STEP 2: confirm

type ConfirmInput struct {
	CustomRef string `json:"custom_ref"`
	StoreID   string `json:"store_id"`
}

type ConfirmResponse struct {
	Message string `json:"message,omitempty"`
	Result  string `json:"result"`
	Data    Data   `json:"data"`
}

type NotesToShop struct {
	PaymentStatus string `json:"payment_status"`
}

type Data struct {
	CustomerPhone   string      `json:"customer_phone"`
	CustomerEmail   string      `json:"customer_email"`
	CustomerName    string      `json:"customer_name"`
	OwnerName       string      `json:"owner_name"`
	ShopName        string      `json:"shop_name"`
	ShopLogo        string      `json:"shop_logo"`
	ShopURL         string      `json:"shop_url"`
	OwnerCity       string      `json:"owner_city"`
	OwnerPhone      string      `json:"owner_phone"`
	OwnerEmail      string      `json:"owner_email"`
	Amount          string      `json:"amount"`
	Currency        string      `json:"currency"`
	GatewayName     string      `json:"gateway_name"`
	Gateway         string      `json:"gateway"`
	GatewayRef      string      `json:"gateway_ref"`
	DateTime        string      `json:"date_time"`
	CreatedAt       time.Time   `json:"created_at"`
	OrderStatus     int         `json:"order_status"`
	OrderType       string      `json:"order_type"`
	OrderHistory    []any       `json:"order_history"`
	NotesToCustomer any         `json:"notes_to_customer"`
	NotesToShop     NotesToShop `json:"notes_to_shop"`
	Reference       string      `json:"reference"`
	OrderID         string      `json:"order_id"`
	CustomRef       string      `json:"custom_ref"`
}
