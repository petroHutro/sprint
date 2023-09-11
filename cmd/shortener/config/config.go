package config

import (
	"errors"
	"flag"
	"strconv"
	"strings"
)

type BaseUrl string

type NetAddress struct {
	Host string
	Port int
}

type Flags struct {
	BaseUrl
	NetAddress
}

func NewFlags() Flags {
	return Flags{
		BaseUrl: "http://localhost:8080",
		NetAddress: NetAddress{
			Host: "localhost",
			Port: 8080,
		},
	}
}

func ParseFlags() *Flags {
	flags := NewFlags()
	flag.Var(&flags.NetAddress, "a", "address and port to run server")
	flag.Var(&flags.BaseUrl, "b", "BaseUrl")
	flag.Parse()
	return &flags
}

func (a BaseUrl) String() string {
	return string(a)
}

func (a *BaseUrl) Set(s string) error {
	*a = BaseUrl(s)
	return nil
}

func (a NetAddress) String() string {
	return a.Host + ":" + strconv.Itoa(a.Port)
}

func (a *NetAddress) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("need address in a form host:port")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	a.Host = hp[0]
	a.Port = port
	return nil
}
