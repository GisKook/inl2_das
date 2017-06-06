package conf

import (
	"encoding/json"
	"os"
)

type ServerConfiguration struct {
	BindPort          string
	ReadLimit         uint16
	WriteLimit        uint16
	ConnTimeout       uint16
	ConnCheckInterval uint16
	ServerStatistics  uint16
}

type NsqConfiguration struct {
	Addr           string
	ReportInterval uint8
	MaxReportCount uint8
	Count          uint8
	TopicRssis     string
}

type Configuration struct {
	Server *ServerConfiguration
	Nsq    *NsqConfiguration
}

var G_conf *Configuration

func ReadConfig(confpath string) (*Configuration, error) {
	file, _ := os.Open(confpath)
	decoder := json.NewDecoder(file)
	config := Configuration{}
	err := decoder.Decode(&config)
	G_conf = &config

	return &config, err
}

func GetConf() *Configuration {
	return G_conf
}
