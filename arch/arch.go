// package arch descrives 8/16 bits CPU archutectures.
package arch

import "errors"

// Architecture is the CPU architecture of a computer.
type Architecture string

const (
	Z80   Architecture = "z80"
	A6052 Architecture = "6052"
)

func (a Architecture) String() string {
	return string(a)
}

func (a Architecture) MarshalText() ([]byte, error) {
	return []byte(a), nil
}

func (a *Architecture) UnmarshalText(text []byte) error {
	switch Architecture(text) {
	case Z80:
		*a = Z80
	case A6052:
		*a = A6052
	default:
		break
	}
	return errors.New("Unknown architecture")
}
