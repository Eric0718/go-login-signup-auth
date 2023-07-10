package controllers

import (
	"errors"
	"html/template"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translation "github.com/go-playground/validator/v10/translations/en"
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
		// manual validation
		// helper.RegistrationManualValidation(w, r)

		user := entities.User{
			FullName:        r.FormValue("full_name"),
			Email:           r.FormValue("email"),
			Username:        r.FormValue("username"),
			Password:        r.FormValue("password"),
			ConfirmPassword: r.FormValue("confirm_password"),
		}
		// validation with validator package
		translator := en.New()
		uni := ut.New(translator, translator)

		trans, _ := uni.GetTranslator("en")

		validate := validator.New()
		en_translation.RegisterDefaultTranslations(validate, trans)

		// change default label tag
		validate.RegisterTagNameFunc(func(field reflect.StructField) string {
			labelName := field.Tag.Get("label")
			return labelName
		})

		// make a custom message validation
		validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
			return ut.Add("required", "{0} can't be empty", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required", fe.Field())
			return t
		})

		// every user can't have a similar email when make a registration
		validate.RegisterValidation("isunique", func(fl validator.FieldLevel) bool {
			params := fl.Param()
			split_params := strings.Split(params, "-")

			tableName := split_params[0]
			// email is index 0
			fieldName := split_params[1]
			// field Email is index 1

			fieldValue := fl.Field().String()
			// fieldValue is used to pooling all the input user

			return checkIsUnique(tableName, fieldName, fieldValue)
		})
		// custome message of the similar email in db
		validate.RegisterTranslation("isunique", trans, func(ut ut.Translator) error {
			return ut.Add("isunique", "{0} already used", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("isunique", fe.Field())
			return t
		})

		vErrors := validate.Struct(user)

		var errMessages = make(map[string]any)

		if vErrors != nil {
			for _, e := range vErrors.(validator.ValidationErrors) {
				errMessages[e.StructField()] = e.Translate(trans)
			}

			data := map[string]any{
				"validation": errMessages,
				// so that the filled text not lost if all of the input text doesn't exist
				"user": user,
			}

			temp, err := template.ParseFiles("views/register.html")
			if err != nil {
				panic(err)
			}
			temp.Execute(w, data)
		}
	}
}

func checkIsUnique(tableName string, fieldName string, fieldValue string) bool {
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}

	row, err := db.Query(`select `+fieldName+` from `+tableName+` where `+fieldName+` = ?`, fieldValue)
	if err != nil {
		panic(err)
	}

	defer row.Close()

	var result string
	for row.Next() {
		row.Scan(&result)
	}

	return result != fieldValue
}
