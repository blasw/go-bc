package types

import (
	"encoding/hex"
	"fmt"
)

type Address [20]uint8

func (a Address) ToSlice() []byte {
	return a[:]
}

func (a Address) String() string {
	return hex.EncodeToString(a.ToSlice())
}

func AddressFromBytes(b []byte) Address {
	if len(b) != 20 {
		msg := fmt.Sprintf("given bytes with length %d, expected 20", len(b))
		panic(msg)
	}

	value := [20]uint8{}
	for i := 0; i < 20; i++ {
		value[i] = b[i]
	}

	return Address(value)
}
