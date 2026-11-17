package core

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

// ErrInvalidRequest is returned by a service Exec when its request fails struct validation.
var ErrInvalidRequest = errors.New("invalid request")

func init() {
	// notblank refuses a string holding only whitespace, which required lets through.
	err := validate.RegisterValidation("notblank", validators.NotBlank)
	if err != nil {
		panic(err)
	}
}
