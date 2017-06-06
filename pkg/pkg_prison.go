package pkg

import (
	//"github.com/giskook/inl2_das/protocol"
	"github.com/giskook/gotcp"
)

type Prison_Packet struct {
	Type   uint16
	Packet gotcp.Packet
}

func (this *Prison_Packet) Serialize() []byte {
	return this.Packet.Serialize()
}

func New_Prison_Pkg(Type uint16, Packet gotcp.Packet) *Prison_Packet {
	return &Prison_Packet{
		Type:   Type,
		Packet: Packet,
	}
}
