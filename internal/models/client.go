package models

type Client struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Service string `json:"service"`
	Target  string `json:"target"`
}
