// package plat describes 8/16 bits platforms
package plat

import "errors"

// Platform is the platform of a computer. It could be an OS or a machine type.
type Platform string

const (
	SMS Platform = "sms" // Sega Master System (tm)
	NES Platform = "nes" // Nintendo Entertainment System (tm)
	BIN Platform = "bin" // Simple binary platform for testing on emulator.
	PIR Platform = "pir" // PIR platform to produce PIR output only.
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

type Info struct {
	BankSize uint16 // ROM/ram bank size.
	NMI      uint16 // NMI vector address
	INT      uint16 // INT vector address
	Stack    uint16 // Stack address.
	RAM      uint16 // RAM address.
}

var SMSInfo = Info{
	BankSize: 0x4000,
	INT:      0x0038,
	NMI:      0x0066,
	Stack:    0xdff0,
	RAM:      0xc000,
}

var BINInfo = Info{
	BankSize: 0x4000,
	INT:      0x0038,
	NMI:      0x0066,
	Stack:    0xdff0,
	RAM:      0xc000,
}

var NESInfo = Info{
	BankSize: 0x2000,
	INT:      0xFFFC,
	NMI:      0xFFFE,
	Stack:    0xdff0, // XXX
	RAM:      0xc000, // XXX
}

func (p Platform) Info() *Info {
	switch p {
	case SMS:
		return &SMSInfo
	case NES:
		return &NESInfo
	case BIN:
		return &BINInfo
	default:
		break
	}
	return nil
}
