package protocol

import (
	"bytes"
	"github.com/giskook/inl2_das/base"
	"log"
)

const (
	PROTOCOL_START_FLAG byte   = 0xce
	PROTOCOL_END_FLAG   byte   = 0xce
	PROTOCOL_MIN_LEN    uint16 = 25
	PROTOCOL_MAX_LEN    uint16 = 1024

	PROTOCOL_ILLEGAL   uint16 = 0
	PROTOCOL_HALF_PACK uint16 = 255

	PROTOCOL_BLUETOOTH_LOCATE uint16 = 0x0001
)

func ParseHeader(buffer []byte) (*bytes.Reader, uint16, uint16) {
	reader := bytes.NewReader(buffer)
	reader.Seek(1, 0)
	length := base.ReadWord(reader)
	protocol_id := base.ReadWord(reader)

	return reader, length, protocol_id
}

func WriteHeader(writer *bytes.Buffer, length uint16, cmdid uint16, cpid uint64) {
	writer.WriteByte(PROTOCOL_START_FLAG)
	base.WriteWord(writer, length)
	base.WriteWord(writer, cmdid)
	base.WriteBcdCpid(writer, cpid)
}

func CalcXor(cmd []byte, cmdlen uint16) byte {
	temp := cmd[0]
	for i := uint16(1); i < cmdlen; i++ {
		temp ^= cmd[i]
	}

	return temp
}

func CheckProtocol(buffer *bytes.Buffer) (uint16, uint16) {
	bufferlen := buffer.Len()
	if bufferlen == 0 {
		return PROTOCOL_ILLEGAL, 0
	}
	if buffer.Bytes()[0] != PROTOCOL_START_FLAG {
		buffer.ReadByte()
		CheckProtocol(buffer)
	} else if bufferlen > 2 {
		pkglen := base.GetWord(buffer.Bytes()[1:3])
		if pkglen < PROTOCOL_MIN_LEN || pkglen > PROTOCOL_MAX_LEN {
			buffer.ReadByte()
			CheckProtocol(buffer)
		}

		if int(pkglen) > bufferlen {
			return PROTOCOL_HALF_PACK, 0
		} else {
			xor_calc := CalcXor(buffer.Bytes()[0:], pkglen-2)
			log.Printf("xor value %x\n", xor_calc)
			if xor_calc == buffer.Bytes()[pkglen-2] && buffer.Bytes()[pkglen-1] == PROTOCOL_END_FLAG {
				protocol_id := base.GetWord(buffer.Bytes()[3:5])
				return protocol_id, pkglen
			} else {
				buffer.ReadByte()
				CheckProtocol(buffer)
			}
		}
	} else {
		return PROTOCOL_HALF_PACK, 0
	}

	return PROTOCOL_HALF_PACK, 0
}
