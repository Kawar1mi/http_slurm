package models

type Car struct {
	Id      int64  `json:"id"`
	Colour  string `json:"colour"`
	Brand   string `json:"brand"`
	License string `json:"license"`
	Owner   User   `json:"owner"`
}
