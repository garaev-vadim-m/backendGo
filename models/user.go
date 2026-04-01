package models

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Country  string `json:"country"`
	Password string `json:"-"`
	Role     Role   `json:"role"`
}
