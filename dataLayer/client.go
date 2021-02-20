package dataLayer

type Client struct {
	Model
	FirstName   string `json:"firstName"`
	MiddleName  string `json:"middleName"`
	Description string `json:"description"`
	UserID      uint64 `json:"userId"`
}
