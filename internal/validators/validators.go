package validators

import (
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-playground/validator/v10"
)

func DOBValidator(f1 validator.FieldLevel) bool {
	minimumDate := time.Now().AddDate(-18, 0, 0)

	dob := f1.Field().Interface().(time.Time)
	log.Info(dob.Date())

	return dob.Before(minimumDate)
}
