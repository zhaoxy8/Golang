package conf

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

//Conf ...
type Conf struct {
	Hosts []string
	Topic string
	Path  string
}

//NewConf ...
func NewConf(filename string) *Conf {
	cfg, err := ini.Load(filename)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	return &Conf{
		Hosts: strings.Split(cfg.Section("kafka").Key("hosts").String(), ","),
		Topic: cfg.Section("kafka").Key("topic").String(),
		Path:  cfg.Section("taillog").Key("path").String(),
	}
}
