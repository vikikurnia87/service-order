package utils

// Facade ke service-utils/validator (paket `validator` bentrok nama dengan
// go-playground/validator/v10 di handler, jadi tetap di-facade lokal).
import suvalidator "github.com/vikikurnia87/service-utils/validator"

type FieldError = suvalidator.FieldError

var (
	NewValidator          = suvalidator.New
	TranslateErrorMessage = suvalidator.TranslateErrorMessage
	ServiceMessageError   = suvalidator.ServiceMessageError
	SetValidatorLang      = suvalidator.SetLang
)
