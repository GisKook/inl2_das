package buffer_worker

import (
	"github.com/giskook/inl2_das/base"
	"github.com/giskook/inl2_das/conf"
	"log"
	"time"
)

type LocatePair struct {
	Key   string
	Value *base.SingleRingRealTimeVector
}

type BufferWorker struct {
	LocateQueue map[uint64]*base.SingleRingRealTimeVector
	LocateAdd   chan *base.SingleRingRealTimeVector
	LocateDel   chan uint64

	ticker *time.Ticker
}

var G_BufferWorker *BufferWorker = nil

func GetBufferWorker() *BufferWorker {
	if G_BufferWorker == nil {
		G_BufferWorker = &BufferWorker{
			LocateQueue: make(map[uint64]*base.SingleRingRealTimeVector),
			LocateAdd:   make(chan *base.SingleRingRealTimeVector),
			LocateDel:   make(chan uint64),
			ticker:      time.NewTicker(time.Duration(conf.GetConf().Nsq.ReportInterval) * time.Second),
		}
	}

	return G_BufferWorker
}

func (z *BufferWorker) Close() {
	close(z.LocateAdd)
	close(z.LocateDel)
	z.ticker.Stop()
}

func (z *BufferWorker) PushLocate(key uint64, locate *base.SingleRingRealTimeVector) {
	z.LocateAdd <- locate
}

func (z *BufferWorker) RemoveLocate(key uint64) {
	z.LocateDel <- key
}

func (z *BufferWorker) Run() {
	go func() {
		for {
			select {
			case <-z.ticker.C:
				log.Println("ticker")
				z.ProccessSendMsg()
			case locate := <-z.LocateAdd:
				z.insert_locate(locate.RingMac, locate)
			}
		}
	}()
}

func (z *BufferWorker) insert_locate(key uint64, value *base.SingleRingRealTimeVector) {
	log.Println("insert_locate")

	_, ok := z.LocateQueue[key]
	if ok {
		z.LocateQueue[key].AccX = value.AccX
		z.LocateQueue[key].AccY = value.AccY
		z.LocateQueue[key].AccZ = value.AccZ
		z.LocateQueue[key].Battery = value.Battery
		z.LocateQueue[key].Alarm = value.Alarm

		//		bHave := false
		//		for i, v := range value.Rssis {
		//			bHave = false
		//			for j, zv := range z.LocateQueue[key].Rssis {
		//				if v.TagMac == zv.TagMac {
		//					z.LocateQueue[key].Rssis[j].Rssi = v.Rssi
		//					bHave = true
		//
		//					break
		//				}
		//			}
		//
		//			if !bHave {
		//				z.LocateQueue[key].Rssis = append(z.LocateQueue[key].Rssis, value.Rssis[i])
		//			}
		//		}
		z.LocateQueue[key].Beacons = append(z.LocateQueue[key].Beacons, value.Beacons...)
	} else {
		z.LocateQueue[key] = value
	}
}
