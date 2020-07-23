package models

type Contact struct {
	ID      uint32 `json:"id"`
	UserID  uint32 `json:"userID"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email"`
}
