package data

import "encoding/gob"

func init() { gob.Register(RawOrderData{}) }

//GENERATED json items data
type RawItemsData struct {
	Price       int    `json:"price"`
	ChrtId      int    `json:"chrt_id"`
	Status      int    `json:"status"`
	TotalPrice  int    `json:"total_price"`
	Sale        int    `json:"sale"`
	Name        string `json:"name"`
	Rid         string `json:"rid"`
	Size        string `json:"size"`
	TrackNumber string `json:"track_number"`
	Brand       string `json:"brand"`
	NmId        int    `json:"nm_id"`
	LocalID     string `json:"-"`
}

// GENERATED json delivery data
type RawDeliveryData struct {
	Phone   string `json:"phone"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Region  string `json:"region"`
	Address string `json:"address"`
	City    string `json:"city"`
	Zip     string `json:"zip"`
	LocalID string `json:"-"`
}

// GENERATED json payment data
type RawPaymentData struct {
	DeliveryCost int    `json:"delivery_cost"`
	Bank         string `json:"bank"`
	PaymentDt    int    `json:"payment_dt"`
	Provider     string `json:"provider"`
	Currency     string `json:"currency"`
	Transaction  string `json:"transaction"`
	CustomFee    int    `json:"custom_fee"`
	GoodsTotal   int    `json:"goods_total"`
	Amount       int    `json:"amount"`
	RequestId    string `json:"request_id"`
	LocalID      string `json:"-"`
}

// GENERATED json Order data
type RawOrderData struct {
	OofShard          string          `json:"oof_shard"`
	Payment           RawPaymentData  `json:"payment"`
	DateCreated       string          `json:"date_created"`
	DeliveryService   string          `json:"delivery_service"`
	InternalSignature string          `json:"internal_signature"`
	Locale            string          `json:"locale"`
	Delivery          RawDeliveryData `json:"delivery"`
	Entry             string          `json:"entry"`
	CustomerId        string          `json:"customer_id"`
	Items             []RawItemsData  `json:"items"`
	OrderUid          string          `json:"order_uid"`
	SmId              int             `json:"sm_id"`
	Shardkey          string          `json:"shardkey"`
	TrackNumber       string          `json:"track_number"`
	LocalID           string          `json:"-"`
}
