package protocol

import (
	"errors"
	"io"
)

const ( SegmentBit = 0x7F
ContinueBit = 0x80

MaxVarIntBytes = 5 )

type (
	VarInt int32
)

func (v *VarInt) ReadVarInt(r io.Reader) (numRead int32, err error) {
	var result int32
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

	*v = VarInt(result)
	return numRead, err
}

func (v VarInt) WriteVarInt(w io.Writer) error {
	for {
		if (v & ^SegmentBit) == 0 {
			b := []byte{byte(v)}
			_, err := w.Write(b)
			return err
		}
		b := byte((v & SegmentBit) | ContinueBit)
		_, err := w.Write([]byte{b})
		if err != nil {
			return err
		}
		v >>= 7
	}
}