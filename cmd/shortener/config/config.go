package config

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"strings"
)

type BaseURL string

type NetAddress struct {
	Host string
	Port int
}

type Flags struct {
	BaseURL
	NetAddress
}

func NewFlags() Flags {
	return Flags{
		BaseURL: "http://localhost:8080",
		NetAddress: NetAddress{
			Host: "localhost",
			Port: 8080,
		},
	}
}

func ConfigureServer() *Flags {
	flags := ParseFlags()
	ParseENV(flags)
	return flags
}

func ParseENV(flags *Flags) {
	if serverAddress := os.Getenv("SERVER_ADDRESS"); serverAddress != "" {
		flags.NetAddress.Set(serverAddress)
	}
	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		flags.BaseURL.Set(baseURL)
	}
}

func ParseFlags() *Flags {
	flags := NewFlags()
	flag.Var(&flags.NetAddress, "a", "address and port to run server")
	flag.Var(&flags.BaseURL, "b", "BaseUrl")
	flag.Parse()
	return &flags
}

func (a BaseURL) String() string {
	return string(a)
}

func (a *BaseURL) Set(s string) error {
	*a = BaseURL(s)
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
