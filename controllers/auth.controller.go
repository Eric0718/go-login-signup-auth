package controllers

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/mhdianrush/go-login-signup-auth/config"
	"github.com/mhdianrush/go-login-signup-auth/entities"
	"github.com/mhdianrush/go-login-signup-auth/models"
	"golang.org/x/crypto/bcrypt"
)

type USerInput struct {
	Username string
	Password string
}

func Index(w http.ResponseWriter, r *http.Request) {
	// checked session login of each client, if empty, will redirect to login page
	session, err := config.Store.Get(r, config.SESSION_ID)
	if err != nil {
		panic(err)
	}

	if len(session.Values) == 0 {
		// empty session
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		if session.Values["LoggedIn"] != true {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			data := map[string]any{
				"full_name_of_user": session.Values["full_name"],
			}
			temp, err := template.ParseFiles("views/index.html")
			if err != nil {
				panic(err)
			}
			temp.Execute(w, data)
		}
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp, err := template.ParseFiles("views/login.html")
		if err != nil {
			panic(err)
		}
		temp.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		// login process
		r.ParseForm()

		userInput := &USerInput{
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		}

		var user entities.User
		err := models.NewUserModel().Find(&user, "username", userInput.Username)
		if err != nil {
			panic(err)
		}

		var message error
		// if login failed

		if user.Username == "" {
			// nothing match data with database
			message = errors.New("username or password doesn't match")
		} else {
			// exist username data in database
			// compare password userInput with password in database with hashing
			errPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password))
			if errPassword != nil {
				// password not match
				message = errors.New("username or password doesn't match")
			}
		}
		// if failed login
		if message != nil {
			data := map[string]any{
				"error": message,
			}
			temp, err := template.ParseFiles("views/login.html")
			if err != nil {
				panic(err)
			}
			temp.Execute(w, data)
		} else {
			// if username and password correct, set session
			session, _ := config.Store.Get(r, config.SESSION_ID)
			session.Values["LoggedIn"] = true
			session.Values["email"] = user.Email
			session.Values["username"] = user.Username
			session.Values["full_name"] = user.FullName

			session.Save(r, w)

			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}

}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := config.Store.Get(r, config.SESSION_ID)
	if err != nil {
		panic(err)
	}
	// delete the session
	session.Options.MaxAge = -1
	// -1 is the default number to delete session

	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp, err := template.ParseFiles("views/register.html")
		if err != nil {
			panic(err)
		}
		temp.Execute(w, nil)
	}

	if r.Method == http.MethodPost {
		r.ParseForm()

		user := entities.User{
			FullName:        r.FormValue("full_name"),
			Email:           r.FormValue("email"),
			Username:        r.FormValue("username"),
			Password:        r.FormValue("password"),
			ConfirmPassword: r.FormValue("confirm_password"),
		}

		errMessage := make(map[string]any)

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
