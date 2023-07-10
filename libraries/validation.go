package libraries

import (
	"database/sql"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translation "github.com/go-playground/validator/v10/translations/en"
	"github.com/mhdianrush/go-login-signup-auth/config"
)

type Validation struct {
	DB *sql.DB
}

func NewValidation() *Validation {
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}
	return &Validation{
		DB: db,
	}
}

func (v *Validation) Init() (*validator.Validate, ut.Translator) {
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

		return v.checkIsUnique(tableName, fieldName, fieldValue)
	})
	// custome message of the similar email in db
	validate.RegisterTranslation("isunique", trans, func(ut ut.Translator) error {
		return ut.Add("isunique", "{0} already used", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("isunique", fe.Field())
		return t
	})

	return validate, trans
}

func (v *Validation) Struct(s any) any {
	validate, trans := v.Init()

	vErrors := make(map[string]any)

	err := validate.Struct(s)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			vErrors[e.StructField()] = e.Translate(trans)
		}
	}

	if len(vErrors) > 0 {
		return vErrors
	}
	return nil
}

func (v *Validation) checkIsUnique(tableName string, fieldName string, fieldValue string) bool {
	row, err := v.DB.Query(`select `+fieldName+` from `+tableName+` where `+fieldName+` = ?`, fieldValue)
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
