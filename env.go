package goenv

import (
	"github.com/ilhamtubagus/condutil"
	"os"
	"reflect"
	"strings"
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
	if condutil.IsZeroValue(envTag) {
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
	if condutil.IsZeroValue(defaultValue) {
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

		return nil
	}

	parseFunc, ok = defaultParser[field.Kind()]
	if ok {
		parsedValue, err := parseFunc(envValue)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(parsedValue))

		return nil
	}

	switch field.Kind() {
	case reflect.Slice:
		err := handleSlice(field, fieldType, envValue, options)
		if err != nil {
			return err
		}
	case reflect.Map:
		err := handleMap(field, fieldType, options)
		if err != nil {
			return err
		}
	default:
		return NoParserFoundError{fieldType.Name}
	}

	return nil
}

func parseMapValue(field reflect.Value, fieldType reflect.StructField, value string, options Options) (reflect.Value, error) {
	parseFunc, ok := options.FuncMap[field.Type().Elem()]
	if ok {
		parsedValue, err := parseFunc(value)
		if err != nil {
			return reflect.Value{}, err
		}

		return reflect.ValueOf(parsedValue), nil
	}

	parseFunc, ok = defaultParser[field.Type().Elem().Kind()]
	if ok {
		parsedValue, err := parseFunc(value)
		if err != nil {
			return reflect.Value{}, err
		}

		return reflect.ValueOf(parsedValue), nil
	}

	return reflect.Value{}, NoParserFoundError{fieldType.Name}
}

func handleMap(field reflect.Value, fieldType reflect.StructField, options Options) error {
	if field.Type().Key().Kind() != reflect.String {
		return InvalidMapKeyError
	}

	matchingEnv := make(map[string]string)

	// Retrieve all environment variables
	envVars := os.Environ()
	envTag := fieldType.Tag.Get(options.TagName)

	for _, env := range envVars {
		// Split into key and value
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			return InvalidEnvironmentVariableError
		}
		key, value := parts[0], parts[1]

		// Check if the key starts with the prefix
		if strings.HasPrefix(key, envTag) {
			mapKey := strings.TrimPrefix(key, envTag+"_")
			mapKey = snakeToCamelCase(mapKey)
			matchingEnv[mapKey] = value
		}
	}

	// Create a new map and set it to the field
	newMap := reflect.MakeMap(field.Type())
	for k, v := range matchingEnv {
		value, err := parseMapValue(newMap, fieldType, v, options)
		if err != nil {
			return err
		}
		newMap.SetMapIndex(reflect.ValueOf(k), value)
	}
	field.Set(newMap)

	return nil
}

func handleSlice(field reflect.Value, fieldType reflect.StructField, value string, options Options) error {
	separator := fieldType.Tag.Get(options.SeparatorTagName)
	if condutil.IsZeroValue(separator) {
		separator = options.separator
	}
	values := strings.Split(value, separator)
	typee := field.Type().Elem()
	if typee.Kind() == reflect.Ptr {
		typee = typee.Elem()
	}

	parserFunc, ok := options.FuncMap[reflect.TypeOf(field)]
	if !ok {
		parserFunc, ok = defaultParser[typee.Kind()]
		if !ok {
			return NoParserFoundError{fieldType: fieldType.Name}
		}
	}
	result := reflect.MakeSlice(fieldType.Type, 0, len(values))
	for _, part := range values {
		r, err := parserFunc(part)
		if err != nil {
			return NoParserFoundError{fieldType: fieldType.Name}
		}
		v := reflect.ValueOf(r).Convert(typee)
		if fieldType.Type.Elem().Kind() == reflect.Ptr {
			v = reflect.New(typee)
			v.Elem().Set(reflect.ValueOf(r).Convert(typee))
		}
		result = reflect.Append(result, v)
	}
	field.Set(result)

	return nil
}
