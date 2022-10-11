package models

type Project struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Title  string `json:"title"`
	Target string `json:"target"`
	Env    string `json:"env"`
}
