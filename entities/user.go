package entities

type User struct {
	Id              int64  `validate:"required" json:"id"`
	FullName        string `validate:"required" label:"Full Name" json:"full_name"`
	Email           string `validate:"required,email,isunique=users-email" json:"email"`
	Username        string `validate:"required,gte=3" json:"username"`
	Password        string `validate:"required,gte=6" json:"password"`
	ConfirmPassword string `validate:"required,eqfield=Password" label:"Confirm Password" json:"confirm_password"`
}
