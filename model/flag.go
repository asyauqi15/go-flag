package model

type Flag struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Active bool   `json:"active"`
}
