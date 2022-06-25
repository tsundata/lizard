package proxy

import (
	"bytes"
	"github.com/tsundata/lizard/config"
	"github.com/tsundata/lizard/util"
	"io"
	"log"
	"net/http"
	"sync"
)

const retryStrategyAttempts = 3

type Proxy struct {
	sync.RWMutex

	addr string
	conf *config.Gateway

	dispatcher *dispatcher
}

func NewProxy(addr string, conf *config.Gateway) *Proxy {
	p := &Proxy{addr: addr, conf: conf}
	p.init()
	return p
}

func (p *Proxy) init() {
	p.initDispatcher()

	p.dispatcher.load()
}

func (p *Proxy) initDispatcher() {
	p.dispatcher = newDispatcher(p.conf)
}

func (p *Proxy) Start() {
	log.Println("start proxy server", p.addr)
	p.startHTTP()
}

func (p *Proxy) Stop() {
	log.Println("stop proxy server")
}

func (p *Proxy) startHTTP() {
	s := p.newHTTPServer()
	err := s.ListenAndServe()
	if err != nil {
		log.Println("error start http listen", err)
	}
}

func (p *Proxy) newHTTPServer() *http.Server {
	return &http.Server{
		Addr:    p.addr,
		Handler: p,
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	buf.WriteByte('[')
	buf.WriteString(req.Method)
	buf.WriteString("] ")
	buf.WriteString(req.RequestURI)
	requestTag := util.ByteToString(buf.Bytes())

	dn := p.dispatcher.dispatch(req, requestTag)
	if dn == nil {
		log.Println("not found proxy", requestTag)
		return
	}

	p.doProxy(dn, w)
}

func (p *Proxy) doProxy(dn *dispatcherNode, w http.ResponseWriter) {
	// todo X-Forwarded-For
	client := &http.Client{}
	forwardReq, err := http.NewRequest(dn.req.Method, "http://"+dn.backend+dn.req.RequestURI, nil)
	// todo copy header
	if err != nil {
		log.Println(err)
		return
	}
	resp, err := client.Do(forwardReq)
	if err != nil {
		log.Println(err)
		return
	}

	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	if body := resp.Body; body != nil {
		_, err = io.Copy(w, body)
		if err != nil {
			log.Println("failed to copy body", dn.requestTag, err)
			return
		}
	}
}
