package validation

import (
	"errors"
	"strings"
)

func Required(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New(field + " is required")
	}
	return nil
}
