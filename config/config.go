package config

import (
	"encoding/json"
	"flag"
	"io"
	"os"
	"strconv"
)

type Config struct {
	Host            string `json:"host"`
	TLS             bool   `json:"tls"`
	Certificate     string `json:"certificate"`
	PrivateKey      string `json:"private"`
	DirectoryFrom   string `json:"dir_from"`
	DirectoryTo     string `json:"dir_to"`
	DB              string `json:"dsn"`
	RefreshInterval int    `json:"refresh_interval"`
	SvgGen          bool   `json:"svg_gen"`
	CFile           string
}

type F struct {
	host            *string
	tls             *bool
	certificate     *string
	privateKey      *string
	directoryFrom   *string
	directoryTo     *string
	db              *string
	refreshInterval *int
	svgGen          *bool
	cFile           *string
}

var f F

const addr = "localhost:8080"

func init() {
	f.host = flag.String("a", addr, "-a=")
	f.tls = flag.Bool("s", false, "-t=")
	f.certificate = flag.String("cert", "", "-cert=path_to_certificate")
	f.privateKey = flag.String("key", "", "-key=path_to_key")
	f.directoryFrom = flag.String("f", "", "-f=from")
	f.db = flag.String("d", "", "-d=db")
	f.directoryTo = flag.String("t", "", "-t=to")
	f.refreshInterval = flag.Int("r", 10, "interval of check")
	f.svgGen = flag.Bool("svg", false, "-svg=")
	f.cFile = flag.String("c", "", "config file")

}

func New() (c Config) {
	flag.Parse()
	if envHost := os.Getenv("HOST"); envHost != "" {
		f.host = &envHost
	}
	if envCert := os.Getenv("CERTIFICATE"); envCert != "" {
		f.certificate = &envCert
	}
	if envKey := os.Getenv("PRIVATE_KEY"); envKey != "" {
		f.privateKey = &envKey
	}
	if envRunDirectoryFrom := os.Getenv("DIRECTORY_FROM"); envRunDirectoryFrom != "" {
		f.directoryFrom = &envRunDirectoryFrom
	}
	if envRunDirectoryTo := os.Getenv("DIRECTORY_TO"); envRunDirectoryTo != "" {
		f.directoryTo = &envRunDirectoryTo
	}
	if envDB := os.Getenv("DATABASE_DSN"); envDB != "" {
		f.db = &envDB
	}
	envRefresh := os.Getenv("REFRESH_INTERVAL")
	if refreshInterval, _ := strconv.Atoi(envRefresh); refreshInterval != 0 {
		f.refreshInterval = &refreshInterval
	}
	c.Host = *f.host
	c.TLS = *f.tls
	c.Certificate = *f.certificate
	c.PrivateKey = *f.privateKey
	c.DirectoryFrom = *f.directoryFrom
	c.DB = *f.db
	c.DirectoryTo = *f.directoryTo
	c.RefreshInterval = *f.refreshInterval
	c.SvgGen = *f.svgGen
	c.CFile = *f.cFile
	file, err := os.Open(c.CFile)
	if err != nil {
		return
	}
	defer file.Close()

	all, err := io.ReadAll(file)
	if err != nil {
		return
	}

	err = json.Unmarshal(all, &c)
	if err != nil {
		return
	}
	return c

}
