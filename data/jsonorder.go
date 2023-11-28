package data
//GENERATED json delivery data
type RawDeliveryData struct{
	Address string `json:"address"`
	City string `json:"city"`
	Zip string `json:"zip"`
	Phone string `json:"phone"`
	Name string `json:"name"`
	Email string `json:"email"`
	Region string `json:"region"`
LocalID string `json:"-"`
}
//GENERATED json payment data
type RawPaymentData struct{
	CustomFee int `json:"custom_fee"`
	Bank string `json:"bank"`
	Amount int `json:"amount"`
	Provider string `json:"provider"`
	Currency string `json:"currency"`
	Transaction string `json:"transaction"`
	GoodsTotal int `json:"goods_total"`
	DeliveryCost int `json:"delivery_cost"`
	PaymentDt int `json:"payment_dt"`
	RequestId string `json:"request_id"`
LocalID string `json:"-"`
}
//GENERATED json items data
type RawItemsData struct{
	Status int `json:"status"`
	Sale int `json:"sale"`
	Name string `json:"name"`
	Rid string `json:"rid"`
	Price int `json:"price"`
	TrackNumber string `json:"track_number"`
	Brand string `json:"brand"`
	NmId int `json:"nm_id"`
	TotalPrice int `json:"total_price"`
	Size string `json:"size"`
	ChrtId int `json:"chrt_id"`
LocalID string `json:"-"`
}
//GENERATED json Order data
type RawOrderData struct{
	CustomerId string `json:"customer_id"`
	Entry string `json:"entry"`
	OrderUid string `json:"order_uid"`
	SmId int `json:"sm_id"`
	InternalSignature string `json:"internal_signature"`
	Items []RawItemsData `json:"items"`
	Payment RawPaymentData `json:"payment"`
	DeliveryService string `json:"delivery_service"`
	Locale string `json:"locale"`
	TrackNumber string `json:"track_number"`
	OofShard string `json:"oof_shard"`
	DateCreated string `json:"date_created"`
	Shardkey string `json:"shardkey"`
	Delivery RawDeliveryData `json:"delivery"`
LocalID string `json:"-"`
}
