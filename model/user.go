package model

type User struct {
	Id       int64  `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Location string `json:"location,omitempty"`
	Age      int64  `json:"age,omitempty"`
}
