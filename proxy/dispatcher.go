package proxy

import (
	"github.com/tsundata/lizard/config"
	"github.com/tsundata/lizard/router"
	"log"
	"net/http"
)

type dispatcher struct {
	conf    *config.Gateway
	routers map[string]*config.Endpoint
	router  *router.Router
}

func newDispatcher(conf *config.Gateway) *dispatcher {
	return &dispatcher{
		conf:    conf,
		routers: make(map[string]*config.Endpoint),
		router:  router.NewRouter(),
	}
}

func (d *dispatcher) load() {
	d.loadRouters()
}

func (d *dispatcher) loadRouters() {
	log.Println("load routers")
	for _, endpoint := range d.conf.Endpoints {
		d.addRouter(endpoint)
	}
}

func (d *dispatcher) addRouter(endpoint *config.Endpoint) {
	key := endpoint.Method + ":" + endpoint.Pattern
	if _, ok := d.routers[key]; ok {
		return
	}
	if endpoint.Method == "" || endpoint.Pattern == "" {
		return
	}

	log.Println("add router", endpoint.Method, endpoint.Pattern)
	d.router.AddRoute(endpoint.Method, endpoint.Pattern)
	d.routers[key] = endpoint
}

func (d *dispatcher) dispatch(req *http.Request, requestTag string) *dispatcherNode {
	node, _ := d.router.MatchRoute(req.Method, req.URL.Path)
	if node != nil {
		key := req.Method + ":" + node.Pattern
		if endpoint, ok := d.routers[key]; ok {
			return &dispatcherNode{
				req:        req,
				backend:    endpoint.Backends[0].Target, //todo
				requestTag: requestTag,
				err:        nil,
				code:       0,
			}
		}
	}
	log.Println("not found", requestTag)
	return nil
}

type dispatcherNode struct {
	req        *http.Request
	backend    string
	requestTag string
	err        error
	code       int
}
