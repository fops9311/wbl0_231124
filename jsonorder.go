package main
//GENERATED json payment data
type rawPaymentData struct{
	CustomFee int `json:"custom_fee"`
	GoodsTotal int `json:"goods_total"`
	Provider string `json:"provider"`
	Currency string `json:"currency"`
	RequestId string `json:"request_id"`
	Transaction string `json:"transaction"`
	DeliveryCost int `json:"delivery_cost"`
	Bank string `json:"bank"`
	PaymentDt int `json:"payment_dt"`
	Amount int `json:"amount"`
}
//GENERATED json delivery data
type rawDeliveryData struct{
	Name string `json:"name"`
	Email string `json:"email"`
	Region string `json:"region"`
	Address string `json:"address"`
	City string `json:"city"`
	Zip string `json:"zip"`
	Phone string `json:"phone"`
}
//GENERATED json items data
type rawItemsData struct{
	Status int `json:"status"`
	TotalPrice int `json:"total_price"`
	Size string `json:"size"`
	Sale int `json:"sale"`
	Price int `json:"price"`
	TrackNumber string `json:"track_number"`
	ChrtId int `json:"chrt_id"`
	Brand string `json:"brand"`
	NmId int `json:"nm_id"`
	Name string `json:"name"`
	Rid string `json:"rid"`
}
//GENERATED json Order data
type rawOrderData struct{
	TrackNumber string `json:"track_number"`
	OrderUid string `json:"order_uid"`
	SmId int `json:"sm_id"`
	Items []rawItemsData `json:"items"`
	Delivery rawDeliveryData `json:"delivery"`
	Shardkey string `json:"shardkey"`
	CustomerId string `json:"customer_id"`
	Payment rawPaymentData `json:"payment"`
	Entry string `json:"entry"`
	OofShard string `json:"oof_shard"`
	DeliveryService string `json:"delivery_service"`
	Locale string `json:"locale"`
	DateCreated string `json:"date_created"`
	InternalSignature string `json:"internal_signature"`
}
