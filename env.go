package go_env

import (
	conditionals "github.com/ilhamtubagus/go-conditionals"
	"os"
	"reflect"
)

type Options struct {
	// TagName is the tag name used to specify the environment variable name.
	TagName string

	// DefaultTagName is the default tag name to be used if no tag name is specified in the struct fields.
	DefaultTagName string

	// SeparatorTagName is the tag name used to specify the separator for splitting the environment variable value into multiple values.
	SeparatorTagName string

	// separator is the separator used to split the environment variable value into multiple values (used on slices or maps).
	separator string

	// FuncMap is a map of custom parsing functions for specific types.
	FuncMap map[reflect.Type]ParseFunc
}

func defaultOptions() Options {
	return Options{
		TagName:          "env",
		DefaultTagName:   "defaultEnv",
		separator:        ",",
		FuncMap:          nil,
		SeparatorTagName: "envSeparator",
	}
}

// Unmarshal populates the fields of the target struct with values from environment variables.
// It uses reflection to iterate through the struct fields and parse their values.
//
// Parameters:
//   - target: A pointer to a struct whose fields will be populated with environment variable values.
//     The struct fields should be tagged with `env:"VARIABLE_NAME"` to specify the
//     corresponding environment variable.
//
// Returns:
//   - error: An error if any issues occur during the unmarshalling process, such as
//     type conversion errors or missing required environment variables.
//     Returns nil if the unmarshalling is successful.
func Unmarshal(target interface{}) error {
	targetRef := reflect.ValueOf(target)

	if targetRef.Kind() != reflect.Ptr {
		return NotStructPtrError{
			actualType: targetRef.Kind().String(),
		}
	}

	typeRef := reflect.TypeOf(target).Elem()
	if targetRef.Elem().Kind() != reflect.Struct {
		return NotStructPtrError{
			actualType: targetRef.Elem().Kind().String(),
		}
	}

	value := targetRef.Elem()
	options := defaultOptions()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := typeRef.Field(i)

		if err := parseField(field, fieldType, options); err != nil {
			return err
		}
	}
	return nil
}

func parseField(field reflect.Value, fieldType reflect.StructField, options Options) error {
	// Recursively parse nested structs
	if field.Kind() == reflect.Struct {
		return Unmarshal(field.Addr().Interface())
	}

	if err := parseEnv(field, fieldType, options); err != nil {
		return err
	}

	return nil
}

func parseEnv(field reflect.Value, fieldType reflect.StructField, options Options) error {
	envTag := fieldType.Tag.Get(options.TagName)
	// skip parsing when env tag is empty
	if conditionals.IsZeroValue(envTag) {
		return nil
	}

	envValue, isPresent := os.LookupEnv(envTag)
	// use default value if environment variable is not found
	if !isPresent && field.Kind() != reflect.Map {
		return parseDefaultEnv(field, fieldType, options)
	}

	return setFieldValue(field, fieldType, envValue, options)
}

func parseDefaultEnv(field reflect.Value, fieldType reflect.StructField, options Options) error {
	defaultValue := fieldType.Tag.Get(options.DefaultTagName)
	if conditionals.IsZeroValue(defaultValue) {
		return nil
	}

	return setFieldValue(field, fieldType, defaultValue, options)
}

func setFieldValue(field reflect.Value, fieldType reflect.StructField, envValue string, options Options) error {
	parseFunc, ok := options.FuncMap[fieldType.Type]
	if ok {
		parsedValue, err := parseFunc(envValue)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(parsedValue))
	}

	parseFunc, ok = defaultParser[field.Kind()]
	if ok {
		parsedValue, err := parseFunc(envValue)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(parsedValue))
	}
	return nil
}
