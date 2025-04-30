package payment_gateway

import (
	"app/pkg/pgtypes"
	"time"

	"github.com/google/uuid"
)

const (
	Bearer = "Bearer "
)

type Settings struct {
	APIKey   string `json:"api_key"`
	Endpoint string `json:"endpoint"`
}

type Response struct {
	ID uuid.UUID `db:"id" json:"id"`

	// PaymentWalletID references users.id foreign field
	PaymentWalletID uuid.UUID `db:"payment_wallet_id"       json:"payment_wallet_id"`
	Type            TypeValue `db:"type"                    json:"type"`
	Amount          float64   `db:"amount"                  json:"amount"`
	PaymentMethod   *string   `db:"payment_method"          json:"payment_method"`

	PaymentReference *string `db:"payment_reference" json:"payment_reference"`

	Notes       *string       `db:"notes"          json:"notes"`
	IsConfirmed bool          `db:"is_confirmed"   json:"is_confirmed"`
	TLyncURL    *string       `db:"tlync_url"      json:"tlync_url"`
	Response    pgtypes.JSONB `db:"response"       json:"response"`
	CreatedAt   time.Time     `db:"created_at"     json:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at"     json:"updated_at"`

	User PaymentWalletUser `db:"user"         json:"user"`
}

type PaymentWalletUser struct {
	ID   *uuid.UUID `db:"id"   json:"id"`
	Name *string    `db:"name" json:"name"`
}

// Tlync
type TlyncRequest struct {
	WalletTransactionID uuid.UUID `json:"wallet_transaction_id"`
	Amount              float64   `json:"amount"`
	Phone               string    `json:"phone"`
}

// Masarat
type MasaratInitiateRequest struct {
	WalletTransactionID uuid.UUID `json:"wallet_transaction_id"`
	Amount              float64   `json:"amount"`
	IdentityCard        string    `json:"identity_card"`
	PaymentServiceID    uuid.UUID `json:"payment_service_id"`
}

type MasaratConfirmRequest struct {
	WalletTransactionID uuid.UUID
	Pin                 string `json:"pin"`
}

// Edfali
type EdfaliInitiateRequest struct {
	WalletTransactionID uuid.UUID `json:"wallet_transaction_id"`
	Amount              float64   `json:"amount"`
	Phone               string    `json:"phone"`
}
type EdfaliConfirmRequest struct {
	WalletTransactionID uuid.UUID
	Pin                 string `json:"pin"`
	Phone               string `json:"phone"`
}

// Sadad
type SadadInitiateRequest struct {
	WalletTransactionID uuid.UUID `json:"wallet_transaction_id"`
	Amount              float64   `json:"amount"`
	Phone               string    `json:"phone"`
	Category            int       `json:"category"`
	Birthyear           string    `json:"birthyear"`
}
type SadadConfirmRequest struct {
	WalletTransactionID uuid.UUID
	Pin                 string `json:"pin"`
}
type SadadResendRequest struct {
	WalletTransactionID uuid.UUID
}

type SadadResendResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Type    string `json:"type"`
}

type PaymentServicesResponse struct {
	Meta struct {
		Total         int      `json:"total"`
		PerPage       int      `json:"per_page"`
		CurrentPage   int      `json:"current_page"`
		FirstPage     int      `json:"first_page"`
		LastPage      int      `json:"last_page"`
		From          int      `json:"from"`
		To            int      `json:"to"`
		Columns       []string `json:"columns"`
		SearchColumns []string `json:"search_columns"`
	} `json:"meta"`
	Data []struct {
		ID    string  `json:"id"`
		Img   *string `json:"img"`
		Thumb *string `json:"thumb"`
		Name  struct {
			Ar string `json:"ar"`
			En string `json:"en"`
		} `json:"name"`
		IsDisabled bool      `json:"is_disabled"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
		Gateway    struct {
			ID   *string `json:"id"`   // Nullable field
			Name *string `json:"name"` // Nullable field
		} `json:"gateway"`
	} `json:"data"`
}
