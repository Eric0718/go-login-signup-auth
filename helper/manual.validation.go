package helper

import (
	"html/template"
	"net/http"

	"github.com/mhdianrush/go-login-signup-auth/entities"
	"github.com/mhdianrush/go-login-signup-auth/models"
	"golang.org/x/crypto/bcrypt"
)

func RegistrationManualValidation(w http.ResponseWriter, r *http.Request) {
	// Manual Validation
	if r.Method == http.MethodPost {
		r.ParseForm()

		user := entities.User{
			FullName:        r.FormValue("full_name"),
			Email:           r.FormValue("email"),
			Username:        r.FormValue("username"),
			Password:        r.FormValue("password"),
			ConfirmPassword: r.FormValue("confirm_password"),
		}

		// Manual Validation
		errMessage := map[string]any{}

		if user.FullName == "" {
			errMessage["FullName"] = "full name must be filled in"
		}
		if user.Email == "" {
			errMessage["Email"] = "email must be filled in"
		}
		if user.Username == "" {
			errMessage["Username"] = "username must be filled in"
		}
		if user.Password == "" {
			errMessage["Password"] = "password must be filled in"
		}
		if user.ConfirmPassword == "" {
			errMessage["ConfirmPassword"] = "confirm password must be filled in"
		} else {
			// if not same with previous input password
			if user.ConfirmPassword != user.Password {
				errMessage["ConfirmPassword"] = "password confirmation does not match"
			}
		}

		if len(errMessage) > 0 {
			// means exist failed registration
			data := map[string]any{
				"validation": errMessage,
			}
			temp, err := template.ParseFiles("views/register.html")
			if err != nil {
				panic(err)
			}
			temp.Execute(w, data)
		} else {
			// success registration
			hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				panic(err)
			}
			user.Password = string(hashPassword)

			// insert to db
			_, err = models.NewUserModel().Create(user)

			var message string

			if err != nil {
				message = "Registration Failed " + err.Error()
			} else {
				message = "Registration Successfully, Please Login"
			}

			data := map[string]any{
				"message": message,
			}

			temp, err := template.ParseFiles("views/register.html")
			if err != nil {
				panic(err)
			}
			temp.Execute(w, data)
		}
	}
}
