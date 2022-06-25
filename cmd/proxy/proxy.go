package main

import (
	"encoding/json"
	"flag"
	"github.com/tsundata/lizard/config"
	"github.com/tsundata/lizard/proxy"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	proxyAddr   string
	proxyConfig string
)

func init() {
	flag.StringVar(&proxyAddr, "addr", ":8080", "proxy address, eg: -addr 0.0.0.0:8080")
	flag.StringVar(&proxyConfig, "config", "config.json", "config path, eg: -config config.json")
}

func main() {
	flag.Parse()

	data, err := ioutil.ReadFile(proxyConfig)
	if err != nil {
		panic(err)
	}

	var conf *config.Gateway
	err = json.Unmarshal(data, &conf)
	if err != nil {
		panic(err)
	}

	p := proxy.NewProxy(proxyAddr, conf)
	go p.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	s := <-c
	log.Println("receive a signal", s.String())

	p.Stop()
}
