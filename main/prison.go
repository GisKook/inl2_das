package main

import (
	"fmt"
	"github.com/giskook/gotcp"
	"github.com/giskook/inl2_das"
	"github.com/giskook/inl2_das/buffer_worker"
	"github.com/giskook/inl2_das/conf"
	"github.com/giskook/inl2_das/event_handler"
	"github.com/giskook/inl2_das/mq"
	"github.com/giskook/inl2_das/server"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// read configuration
	configuration, err := conf.ReadConfig("./conf.json")

	checkError(err)

	// create mq
	mq.NewNsqProducer()
	mq.GetSender().Start()

	// start worker
	buffer_worker.GetBufferWorker().Run()

	// creates a tcp listener
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+configuration.Server.BindPort)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	// creates a tcp server
	config := &gotcp.Config{
		PacketSendChanLimit:    20,
		PacketReceiveChanLimit: 20,
	}
	srv := gotcp.NewServer(config, &event_handler.Callback{}, &inl2_das.Pdas_Protocol{})

	// create inl2_das server
	server_conf := &server.ServerConfig{
		Listener:      listener,
		AcceptTimeout: time.Duration(configuration.Server.ConnTimeout) * time.Second,
	}
	cpd_server := server.NewServer(srv, server_conf)
	server.SetServer(cpd_server)
	// starts service
	fmt.Println("listening:", listener.Addr())
	cpd_server.Start()

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	// stops service
	cpd_server.Stop()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
