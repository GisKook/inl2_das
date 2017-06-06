package protocol

import (
	"github.com/giskook/inl2_das/base"
)

type LocatePacket struct {
	Location *base.SingleRingRealTimeVector
}

func (p *LocatePacket) Serialize() []byte {
	return nil
}

func ParseLocate(buffer []byte) *LocatePacket {
	reader, _, _ := ParseHeader(buffer)
	ring_mac := base.ReadMacInt(reader)
	_degree_x, _ := reader.ReadByte()
	degree_x := int8(_degree_x)
	_degree_y, _ := reader.ReadByte()
	degree_y := int8(_degree_y)
	_degree_z, _ := reader.ReadByte()
	degree_z := int8(_degree_z)
	bett, _ := reader.ReadByte()
	alarm, _ := reader.ReadByte()
	rssi_count, _ := reader.ReadByte()

	var tag_mac uint64
	var _rssi uint8
	rssis := make([]*base.Beacon, rssi_count)

	for i := uint8(0); i < rssi_count; i++ {
		tag_mac = base.ReadMacInt(reader)
		_rssi, _ = reader.ReadByte()
		rssis[i] = &base.Beacon{
			Id:   tag_mac,
			Rssi: float32(int8(_rssi)),
		}
	}

	return &LocatePacket{
		Location: &base.SingleRingRealTimeVector{
			RingMac: ring_mac,
			AccX:    float64(degree_x) / 32,
			AccY:    float64(degree_y) / 32,
			AccZ:    float64(degree_z) / 32,
			Battery: int32(bett),
			Alarm:   int32(alarm),
			Beacons: rssis,
		},
	}
}
