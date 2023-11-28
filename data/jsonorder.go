package data
//GENERATED json payment data
type RawPaymentData struct{
	Provider string `json:"provider"`
	Currency string `json:"currency"`
	DeliveryCost int `json:"delivery_cost"`
	RequestId string `json:"request_id"`
	Transaction string `json:"transaction"`
	CustomFee int `json:"custom_fee"`
	GoodsTotal int `json:"goods_total"`
	Bank string `json:"bank"`
	PaymentDt int `json:"payment_dt"`
	Amount int `json:"amount"`
LocalID string `json:"-"`
}
//GENERATED json delivery data
type RawDeliveryData struct{
	Zip string `json:"zip"`
	Phone string `json:"phone"`
	Name string `json:"name"`
	Email string `json:"email"`
	Region string `json:"region"`
	Address string `json:"address"`
	City string `json:"city"`
LocalID string `json:"-"`
}
//GENERATED json items data
type RawItemsData struct{
	Sale int `json:"sale"`
	Price int `json:"price"`
	ChrtId int `json:"chrt_id"`
	Brand string `json:"brand"`
	Name string `json:"name"`
	Rid string `json:"rid"`
	TrackNumber string `json:"track_number"`
	Status int `json:"status"`
	NmId int `json:"nm_id"`
	TotalPrice int `json:"total_price"`
	Size string `json:"size"`
LocalID string `json:"-"`
}
//GENERATED json Order data
type RawOrderData struct{
	Shardkey string `json:"shardkey"`
	Locale string `json:"locale"`
	Items []RawItemsData `json:"items"`
	Delivery RawDeliveryData `json:"delivery"`
	DateCreated string `json:"date_created"`
	CustomerId string `json:"customer_id"`
	InternalSignature string `json:"internal_signature"`
	OofShard string `json:"oof_shard"`
	SmId int `json:"sm_id"`
	Payment RawPaymentData `json:"payment"`
	Entry string `json:"entry"`
	DeliveryService string `json:"delivery_service"`
	TrackNumber string `json:"track_number"`
	OrderUid string `json:"order_uid"`
LocalID string `json:"-"`
}
