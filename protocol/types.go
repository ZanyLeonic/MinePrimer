package protocol

import (
	"errors"
	"io"
)

const ( 
	SegmentBit = 0x7F
	ContinueBit = 0x80
	MaxVarIntBytes = 5 
)

type (
	VarInt int32
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