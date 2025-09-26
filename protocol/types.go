package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const ( 
	SegmentBit = 0x7F
	ContinueBit = 0x80

	MaxVarIntBytes = 5 
	MaxStringLength = 32767
)

type (
	String string
	VarInt int32
	UShort uint16
)

func ReadVarInt(r io.Reader) (VarInt, error) {
	var result int32
	var numRead int

	for {
		if numRead >= MaxVarIntBytes {
			return 0, errors.New("varint too big")
		}

		var buf [1]byte
		_, err := io.ReadFull(r, buf[:])
		if err != nil {
			return 0, err
		}
		
		fByte := buf[0]
		value := int32(fByte & SegmentBit)
		result |= value << (7 * numRead)
		numRead++

		if (fByte & ContinueBit) == 0 {
			break
		}
	}

	return VarInt(result), nil
}

func WriteVarInt(w io.Writer, value VarInt) error {
	for {
		if (value & ^SegmentBit) == 0 {
			b := []byte{byte(value)}
			_, err := w.Write(b)
			return err
		}
		b := byte((value & SegmentBit) | ContinueBit)
		_, err := w.Write([]byte{b})
		if err != nil {
			return err
		}
		value >>= 7
	}
}

func ReadString(r io.Reader) (String, error) {
	length, err := ReadVarInt(r)
	if err != nil {
		return "", err
	}
	
	if length < 0 {
		return "", errors.New("string length is negative")
	}

	if length > MaxStringLength {
		return "", fmt.Errorf("string length %d exceeds maximum", length)
	}

	buf := make([]byte, length)
	_, err = io.ReadFull(r, buf)

	return String(buf), err
}

func WriteString(w io.Writer, content String) error {
	if err := WriteVarInt(w, VarInt(len(content))); err != nil {
		return err
	}
	_, err := w.Write([]byte(content))
	return err
}

func ReadUnsignedShort(r io.Reader) (UShort, error) {
	var buf [2]byte
	_, err := io.ReadFull(r, buf[:])
	return UShort(binary.BigEndian.Uint16(buf[:])), err
}

func WriteUnsignedShort(w io.Writer, v UShort) error {
	var buf [2]byte
	binary.BigEndian.PutUint16(buf[:], uint16(v))

	_, err := w.Write(buf[:])
	return err
}