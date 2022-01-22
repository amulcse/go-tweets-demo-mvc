package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/amulcse/models"

	goaway "github.com/TwiN/go-away"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/google/uuid"
)

var trans = SetupValidation()
var validate = validator.New()

func InternalServerError(w http.ResponseWriter, err error) {
	uid := uuid.New()
	resp := make(map[string]string)
	resp["message"] = err.Error()
	resp["error_id"] = uid.String()
	responseBody, _ := json.Marshal(resp)
	log.Printf("Internal server error: %s", responseBody)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(responseBody)
}

func BadRequestError(w http.ResponseWriter, err error) {
	for _, err := range err.(validator.ValidationErrors) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		responseBody := map[string]string{"error": err.Translate(trans)}
		if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		}
		return
	}
}

func HasSwearWords(fl validator.FieldLevel) bool {
	return !goaway.IsProfane(fl.Field().String())
}

func SetupValidation() ut.Translator {
	_ = validate.RegisterValidation("HasSwearWords", HasSwearWords)

	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, trans)
	validate.RegisterTranslation("HasSwearWords", trans, func(ut ut.Translator) error {
		return ut.Add("HasSwearWords", "{0} has swear words!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("HasSwearWords", fe.Field())
		return t
	})

	return trans
}

func ValidateRequest(tweet models.Tweet, w http.ResponseWriter) bool {
	err := validate.Struct(tweet)
	if err != nil {
		BadRequestError(w, err)
		return false
	} else {
		return true
	}
}
