package server

import (
	"bytes"
	"encoding/json"
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type SchemaValidator struct {
	validate   *validator.Validate
	translator *ut.Translator
}

func NewSchemaValidator() SchemaValidator {
	var validate = validator.New()
	var translator = en.New()
	var uni = ut.New(translator, translator)

	var trans, found = uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}

	if err := en_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		log.Fatal(err)
	}

	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return SchemaValidator{validate: validate, translator: &trans}
}

func (schemaValidator *SchemaValidator) Validate(dto interface{}) (string, bool) {
	err := schemaValidator.validate.Struct(dto)

	if err == nil {
		return "", true
	}

	var errors []string
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Translate(*schemaValidator.translator))
	}

	return strings.Join(errors, ", "), false
}

func (schemaValidator *SchemaValidator) ValidateStrict(data []byte, dto interface{}) (string, bool) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dto); err != nil {
		return err.Error(), false
	}

	return schemaValidator.Validate(dto)
}
