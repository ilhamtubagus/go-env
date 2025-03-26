package go_env

import (
	"errors"
	"fmt"
)

var InvalidMapKeyError = errors.New("map must have string keys")
var InvalidEnvironmentVariableError = errors.New("environment format is invalid")

// NotStructPtrError The error occurs when pass something that is not a pointer to a struct to Parse
type NotStructPtrError struct {
	actualType string
}

func (e NotStructPtrError) Error() string {
	return fmt.Sprintf("expected a pointer to a struct (found %s)", e.actualType)
}

type NoParserFoundError struct {
	fieldType string
}

func (e NoParserFoundError) Error() string {
	return fmt.Sprintf("no parser found for type %s", e.fieldType)
}
