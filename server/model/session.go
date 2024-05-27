package model

type Session struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Value     string `json:"value"`
	Expires   int    `json:"expires"`
}
