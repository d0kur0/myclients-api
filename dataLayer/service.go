package dataLayer

type Service struct {
	Model
	Name   string `json:"name"`
	Price  int64  `json:"price"`
	UserID uint64 `json:"userId"`
}
