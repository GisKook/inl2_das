package mq

import (
	"log"
	"sync"

	"github.com/bitly/go-nsq"
	"github.com/giskook/inl2_das/conf"
)

var G_NsqSender *NsqProducer
var once sync.Once

type NsqProducer struct {
	producer *nsq.Producer
}

func NewNsqProducer() *NsqProducer {
	once.Do(func() {
		G_NsqSender = &NsqProducer{}
	})

	return G_NsqSender
}

func GetSender() *NsqProducer {
	return G_NsqSender
}

func (s *NsqProducer) Send(topic string, value []byte) error {
	err := s.producer.PublishAsync(topic, value, nil, nil)
	log.Printf("<OUT_NSQ> topic %s %s\n", topic, value)
	if err != nil {
		log.Println("error occur")
		log.Println(err.Error())
	}

	return err
}

func (s *NsqProducer) Start() {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("err found")
		}

	}()
	config := nsq.NewConfig()

	var errmsg error
	s.producer, errmsg = nsq.NewProducer(conf.GetConf().Nsq.Addr, config)

	if errmsg != nil {
		panic("create producer error " + errmsg.Error())
	} else {
		log.Println("producer start ok")
	}
}

func (s *NsqProducer) Stop() {
	s.producer.Stop()
}
