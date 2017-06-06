package inl2_das

import (
	"github.com/giskook/gotcp"
	"github.com/giskook/inl2_das/conn"
	"github.com/giskook/inl2_das/pkg"
	"github.com/giskook/inl2_das/protocol"
	"log"
	"sync"
)

type Pdas_Protocol struct {
}

func (this *Pdas_Protocol) ReadPacket(c *gotcp.Conn) (gotcp.Packet, error) {
	smconn := c.GetExtraData().(*conn.Conn)
	var once sync.Once
	once.Do(smconn.UpdateReadflag)

	buffer := smconn.GetBuffer()

	tcp_conn := c.GetRawConn()
	for {
		if smconn.ReadMore {
			data := make([]byte, 2048)
			readLengh, err := tcp_conn.Read(data)
			log.Printf("<IN>    %x\n", data[0:readLengh])
			if err != nil {
				return nil, err
			}

			if readLengh == 0 {
				return nil, gotcp.ErrConnClosing
			}
			buffer.Write(data[0:readLengh])
		}

		cmdid, pkglen := protocol.CheckProtocol(buffer)
		log.Printf("protocol id %d\n", cmdid)

		pkgbyte := make([]byte, pkglen)
		buffer.Read(pkgbyte)
		switch cmdid {
		case protocol.PROTOCOL_BLUETOOTH_LOCATE:
			log.Println("ReadPacket case")
			p := protocol.ParseLocate(pkgbyte)
			smconn.ReadMore = false

			return pkg.New_Prison_Pkg(protocol.PROTOCOL_BLUETOOTH_LOCATE, p), nil

		case protocol.PROTOCOL_ILLEGAL:
			smconn.ReadMore = true
		case protocol.PROTOCOL_HALF_PACK:
			smconn.ReadMore = true
		}
	}

}
