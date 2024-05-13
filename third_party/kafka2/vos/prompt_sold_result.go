package vos

import "time"

type PromptSoldResult struct {
	OrderId        int64     `json:"orderId"`
	PromptId       int       `json:"promptId"`
	PromptTitle    string    `json:"promptTitle"`
	SellerId       int       `json:"sellerId"`
	BuyerId        int       `json:"buyerId"`
	Price          float64   `json:"price"`
	CreateTime     time.Time `json:"createTime"`
	SellerUserId   int       `json:"sellerUserId"`
	IncomeChannel  int       `json:"incomeChannel"`
	OutcomeChannel int       `json:"outcomeChannel"`
}
