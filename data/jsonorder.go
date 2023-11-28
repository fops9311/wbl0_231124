package data
//GENERATED json items data
type RawItemsData struct{
	Brand string `json:"brand"`
	TotalPrice int `json:"total_price"`
	Size string `json:"size"`
	Sale int `json:"sale"`
	Name string `json:"name"`
	Rid string `json:"rid"`
	TrackNumber string `json:"track_number"`
	Status int `json:"status"`
	NmId int `json:"nm_id"`
	Price int `json:"price"`
	ChrtId int `json:"chrt_id"`
LocalID string `json:"-"`
}
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
	DeliveryCost int `json:"delivery_cost"`
	PaymentDt int `json:"payment_dt"`
	Amount int `json:"amount"`
	Currency string `json:"currency"`
	RequestId string `json:"request_id"`
	GoodsTotal int `json:"goods_total"`
	Bank string `json:"bank"`
	Provider string `json:"provider"`
	Transaction string `json:"transaction"`
LocalID string `json:"-"`
}
//GENERATED json Order data
type RawOrderData struct{
	OrderUid string `json:"order_uid"`
	SmId int `json:"sm_id"`
	Shardkey string `json:"shardkey"`
	Payment RawPaymentData `json:"payment"`
	Delivery RawDeliveryData `json:"delivery"`
	TrackNumber string `json:"track_number"`
	DateCreated string `json:"date_created"`
	DeliveryService string `json:"delivery_service"`
	CustomerId string `json:"customer_id"`
	Items []RawItemsData `json:"items"`
	Entry string `json:"entry"`
	InternalSignature string `json:"internal_signature"`
	OofShard string `json:"oof_shard"`
	Locale string `json:"locale"`
LocalID string `json:"-"`
}
