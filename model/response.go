package model

type Response struct {
	Id      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}
