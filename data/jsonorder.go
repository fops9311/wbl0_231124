package data

import "encoding/gob"

func init() { gob.Register(RawOrderData{}) }

//GENERATED json delivery data
type RawDeliveryData struct {
	Region  string `json:"region"`
	Address string `json:"address"`
	City    string `json:"city"`
	Zip     string `json:"zip"`
	Phone   string `json:"phone"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	LocalID string `json:"-"`
}

// GENERATED json payment data
type RawPaymentData struct {
	Bank         string `json:"bank"`
	PaymentDt    int    `json:"payment_dt"`
	Amount       int    `json:"amount"`
	RequestId    string `json:"request_id"`
	CustomFee    int    `json:"custom_fee"`
	GoodsTotal   int    `json:"goods_total"`
	DeliveryCost int    `json:"delivery_cost"`
	Provider     string `json:"provider"`
	Currency     string `json:"currency"`
	Transaction  string `json:"transaction"`
	LocalID      string `json:"-"`
}

// GENERATED json items data
type RawItemsData struct {
	Sale        int    `json:"sale"`
	Name        string `json:"name"`
	Rid         string `json:"rid"`
	Price       int    `json:"price"`
	TrackNumber string `json:"track_number"`
	ChrtId      int    `json:"chrt_id"`
	Status      int    `json:"status"`
	Brand       string `json:"brand"`
	NmId        int    `json:"nm_id"`
	TotalPrice  int    `json:"total_price"`
	Size        string `json:"size"`
	LocalID     string `json:"-"`
}

// GENERATED json Order data
type RawOrderData struct {
	Items             []RawItemsData  `json:"items"`
	OofShard          string          `json:"oof_shard"`
	DateCreated       string          `json:"date_created"`
	Shardkey          string          `json:"shardkey"`
	Locale            string          `json:"locale"`
	Payment           RawPaymentData  `json:"payment"`
	OrderUid          string          `json:"order_uid"`
	SmId              int             `json:"sm_id"`
	CustomerId        string          `json:"customer_id"`
	Delivery          RawDeliveryData `json:"delivery"`
	DeliveryService   string          `json:"delivery_service"`
	InternalSignature string          `json:"internal_signature"`
	Entry             string          `json:"entry"`
	TrackNumber       string          `json:"track_number"`
	LocalID           string          `json:"-"`
}
