# go-env

go-env is a Go package that provides a simple and flexible way to unmarshal environment variables into Go structs. It uses struct tags to map environment variables to struct fields and supports various data types including nested structs, slices, and maps.

## Features

- Unmarshal environment variables into struct fields
- Support for custom tag names
- Default value support
- Custom separators for slice and map values
- Nested struct support
- Custom parsing functions for specific types

## Installation

To install go-env, use `go get`:

```bash
go get github.com/ilhamtubagus/goenv
```

## Usage
Here's a basic example of how to use go-env:

Supposed you have these environment variables: 
```
HOST=127.0.0.0
PORT=3000
DEBUG=false
DB_NAME=TEST
DB_USER=USER
DB_PASSWORD=PASSWORD
ALLOWED_IPS=127.0.0.0,127.1.1.1
OPTIONS_AUTH_SOURCE=admin
OPTIONS_AUTH_DB=admin
```
then in your go file :
```
package main

import (
    "fmt"
    "github.com/ilhamtubagus/go-env"
)

type Config struct {
    Host     string   `env:"HOST"`
    Port     int      `env:"PORT",defaultEnv:"8080"`
    Debug    bool     `env:"DEBUG"`
    Database struct {
        Name     string `env:"DB_NAME"`
        User     string `env:"DB_USER"`
        Password string `env:"DB_PASSWORD"`
    }
    AllowedIPs []string `env:"ALLOWED_IPS",envSeparator:";"`
    Options   map[string]string `env:"OPTIONS"`
}

func main() {
    var cfg Config
    err := goenv.Unmarshal(&cfg)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("%+v\n", cfg)
}
```

## Struct Tags
- `env`: Specifies the name of the environment variable to use for this field.
- `defaultEnv`: Specifies a default value to use if the environment variable is not set.
- `envSeparator`: Specifies a custom separator for slice values (default is `,`).

## Supported Types
go-env supports the following types:
- Basic types: string, bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64
- Slices of basic types
- Maps with string keys and basic type values
- Nested struct


## Custom Parsing
Ongoing development

## Error Handling
go-env returns descriptive errors for various scenarios, such as:
- Invalid struct pointer
- Missing required environment variables
- Type conversion errors
- Invalid map keys
- Invalid environment variable format

## Contributing
Contributions are welcome! Please feel free to submit a Pull Request.

## License
This project is licensed under the LICENSE file in the root directory of this source tree.
