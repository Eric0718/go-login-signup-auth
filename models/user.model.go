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
	rows, err := u.DB.Query(`select * from users where `+fieldName+` = ? limit 1`, fieldValue)
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

func (u *UserModel) Create(user entities.User) (int64, error) {
	result, err := u.DB.Exec(
		`insert into users (full_name, email, username, password) values (?,?,?,?)`,
		user.FullName, user.Email, user.Username, user.Password,
	)
	if err != nil {
		return 0, err
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	return lastInsertId, nil
}
