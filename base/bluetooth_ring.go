package base

type Beacon struct {
	Id   uint64
	Rssi float32
}

type SingleRingRealTimeVector struct {
	RingMac uint64
	Battery int32
	AccX    float64
	AccY    float64
	AccZ    float64
	Alarm   int32
	Beacons []*Beacon
}
