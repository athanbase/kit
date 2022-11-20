package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Protocol struct {
	Version  uint8
	BodyLen  uint16
	Reserved [2]byte
	Unit     uint8
	Value    uint32
	A        uint64
}

func main() {

	p := Protocol{
		Version: 1,
		BodyLen: 2,
		Unit:    3,
		Value:   5,
		A:       67,
	}
	buffer := new(bytes.Buffer)

	err := binary.Write(buffer, binary.LittleEndian, p)
	if err != nil {
		panic(err)
	}

	binaryData := buffer.Bytes()
	fmt.Printf("protocol obj encoded: %+v\n", binaryData)

	var b Protocol
	err = binary.Read(bytes.NewReader(binaryData), binary.LittleEndian, &b)
	if err != nil {
		panic(err)
	}

	fmt.Printf("decoded binary to protocol obj: %+v\n", b)
}
