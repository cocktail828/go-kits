package configor

import (
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

type Configor struct {
	EnvPrefix   string
	Unmarshaler func([]byte, any) error
}

// New initialize a Configor
func newConfigor() *Configor {
	return &Configor{
		EnvPrefix:   strings.ToUpper(os.Getenv("CONFIGOR_ENV_PREFIX")),
		Unmarshaler: toml.Unmarshal,
	}
}

func (c *Configor) Load(dst any, payload ...[]byte) (err error) {
	return c.load(dst, payload...)
}

// Load will unmarshal configurations to struct from files that you provide
func (c *Configor) LoadFile(dst any, files ...string) error {
	return c.loadFile(dst, files...)
}

// Load will unmarshal configurations to struct from files that you provide
func Load(dst any, payload ...[]byte) error {
	return newConfigor().Load(dst, payload...)
}

// Load will unmarshal configurations to struct from files that you provide
func LoadFile(dst any, files ...string) error {
	return newConfigor().LoadFile(dst, files...)
}
