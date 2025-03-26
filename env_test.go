package go_env

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func loadEnvFromString(envString string) {
	lines := strings.Split(envString, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Ignore empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Split key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			_ = os.Setenv(key, value)
		}
	}
}

func TestUnmarshal_PrimitiveType(t *testing.T) {
	t.Run("Int", func(t *testing.T) {
		envData := `
			INT_FIELD=10
            INT8_FIELD=8
            INT16_FIELD=16
            INT32_FIELD=32
            INT64_FIELD=64
		`
		loadEnvFromString(envData)
		expectedStruct := struct {
			IntField   int   `env:"INT_FIELD"`
			Int8Field  int8  `env:"INT8_FIELD"`
			Int16Field int16 `env:"INT16_FIELD"`
			Int32Field int32 `env:"INT32_FIELD"`
			Int64Field int64 `env:"INT64_FIELD"`
		}{
			IntField:   10,
			Int8Field:  8,
			Int16Field: 16,
			Int32Field: 32,
			Int64Field: 64,
		}
		actualStruct := &struct {
			IntField   int   `env:"INT_FIELD"`
			Int8Field  int8  `env:"INT8_FIELD"`
			Int16Field int16 `env:"INT16_FIELD"`
			Int32Field int32 `env:"INT32_FIELD"`
			Int64Field int64 `env:"INT64_FIELD"`
		}{}

		err := Unmarshal(actualStruct)

		assert.Nil(t, err)
		assert.Equal(t, expectedStruct, *actualStruct)
	})

	t.Run("Uint", func(t *testing.T) {
		envData := `
			UINT_FIELD=10
            UINT8_FIELD=8
            UINT16_FIELD=16
            UINT32_FIELD=32
            UINT64_FIELD=64
		`
		loadEnvFromString(envData)
		expectedStruct := struct {
			UintField   uint   `env:"UINT_FIELD"`
			Uint8Field  uint8  `env:"UINT8_FIELD"`
			Uint16Field uint16 `env:"UINT16_FIELD"`
			Uint32Field uint32 `env:"UINT32_FIELD"`
			Uint64Field uint64 `env:"UINT64_FIELD"`
		}{
			UintField:   10,
			Uint8Field:  8,
			Uint16Field: 16,
			Uint32Field: 32,
			Uint64Field: 64,
		}
		actualStruct := &struct {
			UintField   uint   `env:"UINT_FIELD"`
			Uint8Field  uint8  `env:"UINT8_FIELD"`
			Uint16Field uint16 `env:"UINT16_FIELD"`
			Uint32Field uint32 `env:"UINT32_FIELD"`
			Uint64Field uint64 `env:"UINT64_FIELD"`
		}{}

		err := Unmarshal(actualStruct)

		assert.Nil(t, err)
		assert.Equal(t, expectedStruct, *actualStruct)
	})

	t.Run("Float", func(t *testing.T) {
		envData := `
            FLOAT_FIELD=3.14
            FLOAT32_FIELD=3.14159265359
            FLOAT64_FIELD=3.14159265358979323846
        `
		loadEnvFromString(envData)
		expectedStruct := struct {
			FloatField   float64 `env:"FLOAT_FIELD"`
			Float32Field float32 `env:"FLOAT32_FIELD"`
			Float64Field float64 `env:"FLOAT64_FIELD"`
		}{
			FloatField:   3.14,
			Float32Field: 3.14159265359,
			Float64Field: 3.14159265358979323846,
		}

		actualStruct := &struct {
			FloatField   float64 `env:"FLOAT_FIELD"`
			Float32Field float32 `env:"FLOAT32_FIELD"`
			Float64Field float64 `env:"FLOAT64_FIELD"`
		}{}
		err := Unmarshal(actualStruct)

		assert.Nil(t, err)
		assert.Equal(t, expectedStruct, *actualStruct)
	})

	t.Run("String", func(t *testing.T) {
		envData := `
            STRING_FIELD=hello world
        `
		loadEnvFromString(envData)
		expectedStruct := struct {
			StringField      string `env:"STRING_FIELD"`
			EmptyStringField string
		}{
			StringField:      "hello world",
			EmptyStringField: "",
		}

		actualStruct := &struct {
			StringField      string `env:"STRING_FIELD"`
			EmptyStringField string
		}{}
		err := Unmarshal(actualStruct)

		assert.Nil(t, err)
		assert.Equal(t, expectedStruct, *actualStruct)
	})

	t.Run("Bool", func(t *testing.T) {
		envData := `
            BOOL_FIELD=true
            BOOL_FALSE_FIELD=false
        `
		loadEnvFromString(envData)
		expectedStruct := struct {
			BoolField      bool `env:"BOOL_FIELD"`
			BoolFalseField bool `env:"BOOL_FALSE_FIELD"`
			BoolEmptyField bool
		}{
			BoolField:      true,
			BoolFalseField: false,
			BoolEmptyField: false,
		}

		actualStruct := &struct {
			BoolField      bool `env:"BOOL_FIELD"`
			BoolFalseField bool `env:"BOOL_FALSE_FIELD"`
			BoolEmptyField bool
		}{}
		err := Unmarshal(actualStruct)

		assert.Nil(t, err)
		assert.Equal(t, expectedStruct, *actualStruct)
	})
}

func TestUnmarshal_SliceType(t *testing.T) {
	t.Run("Int slice", func(t *testing.T) {
		envData := `
		SLICE_INT=1,2,3
		SLICE_INT8=4,5,6
		SLICE_INT16=7,8,9
		SLICE_INT32=10,11,12
		SLICE_INT64=13,14,15
	`
		loadEnvFromString(envData)
		expectedStruct := struct {
			SliceInt   []int   `env:"SLICE_INT"`
			SliceInt8  []int8  `env:"SLICE_INT8"`
			SliceInt16 []int16 `env:"SLICE_INT16"`
			SliceInt32 []int32 `env:"SLICE_INT32"`
			SliceInt64 []int64 `env:"SLICE_INT64"`
		}{
			SliceInt:   []int{1, 2, 3},
			SliceInt8:  []int8{4, 5, 6},
			SliceInt16: []int16{7, 8, 9},
			SliceInt32: []int32{10, 11, 12},
			SliceInt64: []int64{13, 14, 15},
		}
		actualStruct := &struct {
			SliceInt   []int   `env:"SLICE_INT"`
			SliceInt8  []int8  `env:"SLICE_INT8"`
			SliceInt16 []int16 `env:"SLICE_INT16"`
			SliceInt32 []int32 `env:"SLICE_INT32"`
			SliceInt64 []int64 `env:"SLICE_INT64"`
		}{}

		err := Unmarshal(actualStruct)

		assert.Nil(t, err)
		assert.Equal(t, expectedStruct, *actualStruct)
	})

	t.Run("Uint slice", func(t *testing.T) {
		envData := `
            SLICE_UINT=1,2,3
            SLICE_UINT8=4,5,6
            SLICE_UINT16=7,8,9
            SLICE_UINT32=10,11,12
            SLICE_UINT64=13,14,15
        `
		loadEnvFromString(envData)
		expectedStruct := struct {
			SliceUint   []uint   `env:"SLICE_UINT"`
			SliceUint8  []uint8  `env:"SLICE_UINT8"`
			SliceUint16 []uint16 `env:"SLICE_UINT16"`
			SliceUint32 []uint32 `env:"SLICE_UINT32"`
			SliceUint64 []uint64 `env:"SLICE_UINT64"`
		}{
			SliceUint:   []uint{1, 2, 3},
			SliceUint8:  []uint8{4, 5, 6},
			SliceUint16: []uint16{7, 8, 9},
			SliceUint32: []uint32{10, 11, 12},
			SliceUint64: []uint64{13, 14, 15},
		}
		actualStruct := &struct {
			SliceUint   []uint   `env:"SLICE_UINT"`
			SliceUint8  []uint8  `env:"SLICE_UINT8"`
			SliceUint16 []uint16 `env:"SLICE_UINT16"`
			SliceUint32 []uint32 `env:"SLICE_UINT32"`
			SliceUint64 []uint64 `env:"SLICE_UINT64"`
		}{}

		err := Unmarshal(actualStruct)

		assert.Nil(t, err)
		assert.Equal(t, expectedStruct, *actualStruct)
	})

	t.Run("Float slice", func(t *testing.T) {
		envData := `
			SLICE_FLOAT=3.14,3.14159,3.14159265359
			SLICE_FLOAT32=3.14159265359,3.14159265358979323846,3.14159265358979323846
			SLICE_FLOAT64=3.14159265358979323846,3.14159265358979323846,3.14159265358979323846
		`
		loadEnvFromString(envData)
		expectedStruct := struct {
			SliceFloat   []float64 `env:"SLICE_FLOAT"`
			SliceFloat32 []float32 `env:"SLICE_FLOAT32"`
			SliceFloat64 []float64 `env:"SLICE_FLOAT64"`
		}{
			SliceFloat:   []float64{3.14, 3.14159, 3.14159265359},
			SliceFloat32: []float32{3.14159265359, 3.14159265358979323846, 3.14159265358979323846},
			SliceFloat64: []float64{3.14159265358979323846, 3.14159265358979323846, 3.14159265358979323846},
		}
		actualStruct := &struct {
			SliceFloat   []float64 `env:"SLICE_FLOAT"`
			SliceFloat32 []float32 `env:"SLICE_FLOAT32"`
			SliceFloat64 []float64 `env:"SLICE_FLOAT64"`
		}{}

		err := Unmarshal(actualStruct)

		assert.Nil(t, err)
		assert.Equal(t, expectedStruct, *actualStruct)
	})

	t.Run("String slice", func(t *testing.T) {
		envData := `
            SLICE_STRING=hello,world,golang
            SLICE_STRING_EMPTY=
        `
		loadEnvFromString(envData)
		expectedStruct := struct {
			SliceString      []string `env:"SLICE_STRING"`
			SliceStringEmpty []string `env:"SLICE_STRING_EMPTY"`
		}{
			SliceString:      []string{"hello", "world", "golang"},
			SliceStringEmpty: []string{""},
		}
		actualStruct := &struct {
			SliceString      []string `env:"SLICE_STRING"`
			SliceStringEmpty []string `env:"SLICE_STRING_EMPTY"`
		}{}

		err := Unmarshal(actualStruct)

		assert.Nil(t, err)
		assert.Equal(t, expectedStruct, *actualStruct)
	})
}
func TestUnmarshal_NotStructPtr(t *testing.T) {
	t.Run("Non-pointer value", func(t *testing.T) {
		nonPtr := struct {
			Field string `env:"FIELD"`
		}{}

		err := Unmarshal(nonPtr)

		assert.Error(t, err)
		assert.IsType(t, NotStructPtrError{}, err)
		assert.Equal(t, "expected a pointer to a struct (found struct)", err.Error())
	})

	t.Run("Pointer to non-struct", func(t *testing.T) {
		nonStruct := "string"
		err := Unmarshal(&nonStruct)

		assert.Error(t, err)
		assert.IsType(t, NotStructPtrError{}, err)
		assert.Equal(t, "expected a pointer to a struct (found string)", err.Error())
	})
}

func TestUnmarshal_MixedType(t *testing.T) {
	envData := `
		INT=1
        STRING=hello
		BOOL=true
        SLICE_INT=1,2,3
		SLICE_STRING=hello,world,golang
	`
	loadEnvFromString(envData)
	expectedStruct := struct {
		Int         int      `env:"INT"`
		String      string   `env:"STRING"`
		Bool        bool     `env:"BOOL"`
		SliceInt    []int    `env:"SLICE_INT"`
		SliceString []string `env:"SLICE_STRING"`
	}{
		Int:         1,
		String:      "hello",
		Bool:        true,
		SliceInt:    []int{1, 2, 3},
		SliceString: []string{"hello", "world", "golang"},
	}

	actualStruct := &struct {
		Int         int      `env:"INT"`
		String      string   `env:"STRING"`
		Bool        bool     `env:"BOOL"`
		SliceInt    []int    `env:"SLICE_INT"`
		SliceString []string `env:"SLICE_STRING"`
	}{}
	err := Unmarshal(actualStruct)

	assert.Nil(t, err)
	assert.Equal(t, expectedStruct, *actualStruct)
}

func TestUnmarshal_NestedStruct(t *testing.T) {
	envData := `
        NESTED_STRUCT_INT=1
        NESTED_STRUCT_STRING=hello
        BOOL=true
        SLICE_INT=1,2,3
		SLICE_STRING=hello,world,golang
    `
	loadEnvFromString(envData)
	type NestedStruct struct {
		Int    int    `env:"NESTED_STRUCT_INT"`
		String string `env:"NESTED_STRUCT_STRING"`
	}
	type MainStruct struct {
		NestedStruct NestedStruct `env:"NESTED_STRUCT"`
		Bool         bool         `env:"BOOL"`
		SliceInt     []int        `env:"SLICE_INT"`
		SliceString  []string     `env:"SLICE_STRING"`
	}
	expectedStruct := MainStruct{
		NestedStruct: NestedStruct{
			Int:    1,
			String: "hello",
		},
		Bool:        true,
		SliceInt:    []int{1, 2, 3},
		SliceString: []string{"hello", "world", "golang"},
	}

	actualStruct := &MainStruct{}
	err := Unmarshal(actualStruct)

	assert.Nil(t, err)
	assert.Equal(t, expectedStruct, *actualStruct)
}

func TestUnmarshal_Map(t *testing.T) {
	t.Run("String map", func(t *testing.T) {
		envData := `
			MAP_FIELD_ONE=value1
			MAP_FIELD_TWO=value2
			MAP_THREE=value3
		`
		loadEnvFromString(envData)
		type MapStruct struct {
			Map map[string]string `env:"MAP"`
		}
		expectedStruct := MapStruct{
			Map: map[string]string{
				"fieldOne": "value1",
				"fieldTwo": "value2",
				"three":    "value3",
			},
		}

		actualStruct := &MapStruct{}
		err := Unmarshal(actualStruct)

		assert.Nil(t, err)
		assert.Equal(t, expectedStruct, *actualStruct)
	})

	t.Run("No parser found error", func(t *testing.T) {
		envData := `
			MAP_FIELD_ONE=value1
			MAP_FIELD_TWO=value2
			MAP_THREE=value3
		`
		loadEnvFromString(envData)
		type MapStruct struct {
			Map map[string]*string `env:"MAP"`
		}

		actualStruct := &MapStruct{}
		err := Unmarshal(actualStruct)

		assert.NotNil(t, err)
		assert.IsType(t, NoParserFoundError{}, err)
	})
}
