package data

import "encoding/gob"

func init() { gob.Register(RawOrderData{}) }

//GENERATED json delivery data
type RawDeliveryData struct {
	Address string `json:"address"`
	City    string `json:"city"`
	Zip     string `json:"zip"`
	Phone   string `json:"phone"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Region  string `json:"region"`
	LocalID string `json:"-"`
}

// GENERATED json items data
type RawItemsData struct {
	Price       int    `json:"price"`
	ChrtId      int    `json:"chrt_id"`
	Status      int    `json:"status"`
	Brand       string `json:"brand"`
	TotalPrice  int    `json:"total_price"`
	Size        string `json:"size"`
	Sale        int    `json:"sale"`
	Name        string `json:"name"`
	Rid         string `json:"rid"`
	TrackNumber string `json:"track_number"`
	NmId        int    `json:"nm_id"`
	LocalID     string `json:"-"`
}

// GENERATED json payment data
type RawPaymentData struct {
	CustomFee    int    `json:"custom_fee"`
	PaymentDt    int    `json:"payment_dt"`
	Provider     string `json:"provider"`
	Currency     string `json:"currency"`
	Transaction  string `json:"transaction"`
	GoodsTotal   int    `json:"goods_total"`
	DeliveryCost int    `json:"delivery_cost"`
	Bank         string `json:"bank"`
	Amount       int    `json:"amount"`
	RequestId    string `json:"request_id"`
	LocalID      string `json:"-"`
}

// GENERATED json Order data
type RawOrderData struct {
	SmId              int             `json:"sm_id"`
	Shardkey          string          `json:"shardkey"`
	DeliveryService   string          `json:"delivery_service"`
	CustomerId        string          `json:"customer_id"`
	InternalSignature string          `json:"internal_signature"`
	DateCreated       string          `json:"date_created"`
	Locale            string          `json:"locale"`
	Payment           RawPaymentData  `json:"payment"`
	Entry             string          `json:"entry"`
	Items             []RawItemsData  `json:"items"`
	Delivery          RawDeliveryData `json:"delivery"`
	OofShard          string          `json:"oof_shard"`
	TrackNumber       string          `json:"track_number"`
	OrderUid          string          `json:"order_uid"`
	LocalID           string          `json:"-"`
}
