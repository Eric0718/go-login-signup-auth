package entities

type User struct {
	Id       int64  `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}


