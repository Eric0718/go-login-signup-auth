package models

import (
	"database/sql"

	"github.com/mhdianrush/go-login-signup-auth/config"
	"github.com/mhdianrush/go-login-signup-auth/entities"
)

type UserModel struct {
	DB *sql.DB
}

func NewUserModel() *UserModel {
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}

	return &UserModel{
		DB: db,
	}
}

func (u *UserModel) Find(user *entities.User, fieldName string, fieldValue string) error {
	rows, err := u.DB.Query(`select * from users where `+fieldName+` = ? limit = 1`, fieldValue)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(
			&user.Id,
			&user.FullName,
			&user.Email,
			&user.Username, 
			&user.Password,
		)
	}
	return nil
}