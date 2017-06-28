package buffer_worker

import (
	"github.com/giskook/inl2_das/base"
	"github.com/giskook/inl2_das/conf"
	"github.com/giskook/inl2_das/mq"
	"github.com/giskook/inl2_das/pb"
	"github.com/golang/protobuf/proto"
	"time"
)

type _tag_mac_rssi struct {
	TagMac uint64
	Rssi   float32
	Count  int
}

func (z *BufferWorker) GetAvgRssi(in []*base.Beacon) []*base.Beacon {
	if in == nil {
		return nil
	}
	tag_mac_rssi := []*_tag_mac_rssi{
		&_tag_mac_rssi{
			TagMac: in[0].Id,
			Rssi:   in[0].Rssi,
			Count:  1,
		},
	}
	bhave := false
	for i := 1; i < len(in); i++ {
		bhave = false
		for j := 0; j < len(tag_mac_rssi); j++ {
			if tag_mac_rssi[j].TagMac == in[i].Id {
				bhave = true
				tag_mac_rssi[j].Rssi += in[i].Rssi
				tag_mac_rssi[j].Count++
			}
		}

		if !bhave {
			tag_mac_rssi = append(tag_mac_rssi, &_tag_mac_rssi{
				TagMac: in[i].Id,
				Rssi:   in[i].Rssi,
				Count:  1,
			})
		}
	}

	length := len(tag_mac_rssi)

	mac_rssis := make([]*base.Beacon, length)

	for k, value := range tag_mac_rssi {
		mac_rssis[k] = &base.Beacon{
			Id:   value.TagMac,
			Rssi: value.Rssi / float32(value.Count),
		}
	}

	return mac_rssis
}

func (z *BufferWorker) PreProccessMsg() {
	for i, m := range z.LocateQueue {
		z.LocateQueue[i].Beacons = z.GetAvgRssi(m.Beacons)
	}
}

func (z *BufferWorker) ProccessSendMsg() {
	if len(z.LocateQueue) > 0 {
		z.PreProccessMsg()
		time_recv := time.Now().Unix() * 1000
		rtvs := &RealTimeVector.RealTimeVectors{
			TimeRecv:      time_recv,
			SingleRingRtv: make([]*RealTimeVector.SingleRingRealTimeVector, 0),
		}
		for i, m := range z.LocateQueue {
			tag_mac_rssi_count := len(m.Beacons)
			if tag_mac_rssi_count > int(conf.GetConf().Nsq.MaxReportCount) {
				tag_mac_rssi_count = int(conf.GetConf().Nsq.MaxReportCount)
			}
			indoor := &RealTimeVector.SingleRingRealTimeVector{
				RingMac: m.RingMac,
				AccX:    m.AccX,
				AccY:    m.AccY,
				AccZ:    m.AccZ,
				Battery: m.Battery,
				Alarm:   m.Alarm,
				Beacons: make([]*RealTimeVector.Beacon, tag_mac_rssi_count),
			}

			for j := 0; j < tag_mac_rssi_count; j++ {
				indoor.Beacons[j] = &RealTimeVector.Beacon{
					Id:   m.Beacons[j].Id,
					Rssi: m.Beacons[j].Rssi,
				}
			}
			rtvs.SingleRingRtv = append(rtvs.SingleRingRtv, indoor)
			delete(z.LocateQueue, i)
		}
		data, _ := proto.Marshal(rtvs)

		mq.GetSender().Send(conf.GetConf().Nsq.TopicRssis, data)
	}
}
