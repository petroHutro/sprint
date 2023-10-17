package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// type BaseURL string

type URLBase struct {
	BaseURL string
}

type Logger struct {
	FilePath  string
	FileFlag  bool
	MultiFlag bool
}

type Storage struct {
	DatabaseDSN     string
	FileStoragePath string
}

type NetAddress struct {
	Host string
	Port int
}

type Flags struct {
	URLBase
	NetAddress
	Logger
	Storage
}

func NewFlags() Flags {
	return Flags{
		URLBase: URLBase{""},
		NetAddress: NetAddress{
			Host: "localhost",
			Port: 8080,
		},
		Logger: Logger{
			FilePath:  "file.log",
			FileFlag:  false,
			MultiFlag: false,
		},
		Storage: Storage{
			// FileStoragePath: "/tmp/short-url-db.json",
			FileStoragePath: "",
			DatabaseDSN:     "",
			// DatabaseDSN:     "host=localhost user=url password=1234 dbname=url sslmode=disable",
		},
	}
}

func LoadServerConfigure() *Flags {
	flags := parseFlags()
	parseENV(flags)
	return flags
}

func parseENV(flags *Flags) {
	if serverAddress := os.Getenv("SERVER_ADDRESS"); serverAddress != "" {
		flags.NetAddress.Set(serverAddress)
	}
	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		flags.BaseURL = baseURL
		// flags.BaseURL.Set(baseURL)
	}
	if fileLoggerPath := os.Getenv("LOGGER_FILE"); fileLoggerPath != "" {
		flags.FilePath = fileLoggerPath
	}
	if fileStoragePath := os.Getenv("FILE_STORAGE_PATH"); fileStoragePath != "" {
		flags.FileStoragePath = fileStoragePath
	}
	if databaseDSN := os.Getenv("DATABASE_DSN"); databaseDSN != "" {
		flags.DatabaseDSN = databaseDSN
	}
}

func parseFlags() *Flags {
	flags := NewFlags()
	flag.Var(&flags.NetAddress, "a", "address and port to run server")
	// flag.Var(&flags.BaseURL, "b", "BaseUrl")
	flag.StringVar(&flags.BaseURL, "b", "http://localhost:8080", "BaseUrl")
	flag.StringVar(&flags.FileStoragePath, "f", "/tmp/short-url-db.json", "FileStoragePath")
	flag.StringVar(&flags.DatabaseDSN, "d", "", "DatabaseDSN")
	flag.BoolVar(&flags.Logger.FileFlag, "l", false, "Logger only file")
	flag.BoolVar(&flags.Logger.MultiFlag, "L", false, "Logger Multi")
	flag.Parse()
	return &flags
}

// func (a BaseURL) String() string {
// 	return string(a)
// }

// func (a *BaseURL) Set(s string) error {
// 	*a = BaseURL(s)
// 	return nil
// }

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
		return fmt.Errorf("cannot atoi port: %w", err)
	}
	a.Host = hp[0]
	a.Port = port
	return nil
}
