// package plat describes 8/16 bits platforms
package plat

import "errors"

// Platform is the platform of a computer. It could be an OS or a machine type.
type Platform string

const (
	SMS Platform = "sms" // Sega Master System (tm)
	NES Platform = "nes" // Nintendo Entertainment System (tm)
	BIN Platform = "bin" // Simple binary platform for testing on emulator.
)

func (p Platform) String() string {
	return string(p)
}

func (p Platform) MarshalText() ([]byte, error) {
	return []byte(p), nil
}

func (p *Platform) UnmarshalText(text []byte) error {
	switch Platform(text) {
	case SMS:
		*p = SMS
	case NES:
		*p = NES
	case BIN:
		*p = BIN
	default:
		break
	}
	return errors.New("Unknown Platform")
}
