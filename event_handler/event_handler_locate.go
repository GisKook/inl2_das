package event_handler

import (
	"github.com/giskook/gotcp"
	"github.com/giskook/inl2_das/buffer_worker"
	"github.com/giskook/inl2_das/pkg"
	"github.com/giskook/inl2_das/protocol"
	"log"
)

func event_handler_blue_tooth_locate(c *gotcp.Conn, p *pkg.Prison_Packet) {
	log.Println("event_handler_blue_tooth_locate")
	locate_pkg := p.Packet.(*protocol.LocatePacket)
	buffer_worker.GetBufferWorker().PushLocate(locate_pkg.Location.RingMac, locate_pkg.Location)
}
