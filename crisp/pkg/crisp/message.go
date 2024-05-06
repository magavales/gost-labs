package crisp

import (
	"encoding/hex"
	"fmt"
)

var (
	ExternalKeyIdFlagWithVersion = []byte{ // 1 bit + 15 bits
		0x00, 0x00,
	}
	CS = []byte{ // 8 bits
		0xf8,
	}
	KeyId = []byte{ // 8 bits
		0x80,
	}
)

type Message struct {
	ExternalKeyIdFlagWithVersion []byte
	CS                           []byte
	KeyId                        []byte
	SeqNum                       []byte
	Payload                      []byte
	ICV                          []byte
	Digits                       []byte
}

func (m *Message) String() string {
	format :=
		`Message:
    ExternalKeyIdFlagWithVersion: %s
    CS:                           %s
    KeyId:                        %s
    SeqNum:                       %s
    Payload:                      %s
    ICV:                          %s
    As block:                     %s`

	return fmt.Sprintf(format,
		hex.EncodeToString(m.ExternalKeyIdFlagWithVersion),
		hex.EncodeToString(m.CS),
		hex.EncodeToString(m.KeyId),
		hex.EncodeToString(m.SeqNum),
		hex.EncodeToString(m.Payload),
		hex.EncodeToString(m.ICV),
		hex.EncodeToString(m.Digits))
}
