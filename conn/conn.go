package conn

import (
	"bytes"
	"github.com/giskook/gotcp"
	"log"
	"time"
)

type ConnConfig struct {
	ConnCheckInterval uint16
	ReadLimit         uint16
	WriteLimit        uint16
	NsqChanLimit      uint16
}

type Conn struct {
	conn          *gotcp.Conn
	config        *ConnConfig
	recieveBuffer *bytes.Buffer
	ticker        *time.Ticker
	readflag      int64
	writeflag     int64
	chan_read     chan int64
	chan_write    chan int64
	closeChan     chan bool
	index         uint32
	ID            uint64
	Status        uint8
	ReadMore      bool
}

func NewConn(conn *gotcp.Conn, config *ConnConfig) *Conn {
	return &Conn{
		conn:          conn,
		recieveBuffer: bytes.NewBuffer([]byte{}),
		config:        config,
		readflag:      time.Now().Unix(),
		writeflag:     time.Now().Unix(),
		chan_read:     make(chan int64),
		chan_write:    make(chan int64),
		ticker:        time.NewTicker(time.Duration(config.ConnCheckInterval) * time.Second),
		closeChan:     make(chan bool),
		index:         0,
		ReadMore:      true,
	}
}

func (c *Conn) Close() {
	c.closeChan <- true
	c.ticker.Stop()
	c.recieveBuffer.Reset()
	close(c.closeChan)
}

func (c *Conn) GetBuffer() *bytes.Buffer {
	return c.recieveBuffer
}

func (c *Conn) SendToTerm(p gotcp.Packet) {
	log.Printf("<OUT> %x \n", p.Serialize())
	c.conn.AsyncWritePacket(p, time.Second)
}

func (c *Conn) UpdateReadflag() {
	c.chan_read <- time.Now().Unix()
}

func (c *Conn) UpdateWriteflag() {
	c.chan_write <- time.Now().Unix()
}

func (c *Conn) checkHeart() {
	defer func() {
		c.conn.Close()
	}()

	var now int64
	for {
		select {
		case rflag := <-c.chan_read:
			c.readflag = rflag
		case wflag := <-c.chan_write:
			c.writeflag = wflag
		case <-c.ticker.C:
			now = time.Now().Unix()
			if now-c.readflag > int64(c.config.ReadLimit) {
				log.Printf("read limit %x\n", c.ID)
				return
			}
			//			if now-c.writeflag > int64(c.config.WriteLimit) {
			//				log.Println("write limit")
			//				return
			//			}
		case <-c.closeChan:
			log.Println("recv close")
			return
		}
	}
}

func (c *Conn) Do() {
	go c.checkHeart()
}
