package handlers

type Request struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Term     string `json:"term"`
	ID       uint32 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
