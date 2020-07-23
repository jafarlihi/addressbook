package models

type ContactList struct {
	ID     uint32 `json:"id"`
	UserID uint32 `json:"userID"`
	Name   string `json:"name"`
}
