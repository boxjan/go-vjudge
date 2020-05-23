package httpclient

import (
	"errors"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"
)

var ErrEmptyTransportPool = errors.New("empty transport pool")
var ErrNotFoundTransportPoolElement = errors.New("not found id in client pool")

type ClientsPool struct {
	defaultConnectTimeout time.Duration
	defaultUa string
	transports []http.Transport
	transportsMapByProxyUrl map[string]int
	clients []Client
	clientsMapByClientName map[string]int
	sync.RWMutex
}

func NewClientPools(defaultConnectTimeout time.Duration, Ua string,
	proxies ...string) (clientsPools *ClientsPool, errs []error) {
	clientsPools = &ClientsPool{
		defaultConnectTimeout: 		defaultConnectTimeout,
		defaultUa: 			  		Ua,
		transports:			  		make([]http.Transport, 0, 8),
		transportsMapByProxyUrl: 	make(map[string]int),
		clients:              		make([]Client, 0, 8),
		clientsMapByClientName: 	make(map[string]int),
	}

	clientsPools.transports = append(clientsPools.transports, http.Transport{})
	clientsPools.transportsMapByProxyUrl["direct"] = len(clientsPools.transports) - 1

	for _, proxy := range proxies {
		err := clientsPools.AddProxy(proxy)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return
}

func (cp *ClientsPool) AddProxy(proxyUrl string) error {
	if _, ok := cp.transportsMapByProxyUrl[proxyUrl]; ok {
		return nil
	}
	proxy, err := url.Parse(proxyUrl)
	if err != nil {
		return errors.New("the proxy url:" + proxyUrl + "can not be parse")
	}
	cp.Lock()
	defer cp.Unlock()
	cp.transports = append(cp.transports, http.Transport{Proxy:http.ProxyURL(proxy)})
	cp.transportsMapByProxyUrl[proxyUrl] = len(cp.transports) - 1
	return nil
}

func (cp *ClientsPool)GetTransportRandom() (*http.Transport, error) {
	cp.RLock()
	defer cp.RUnlock()
	return cp.GetTransportById(rand.Int() % len(cp.transports))
}

func (cp *ClientsPool)GetTransportById(id int) (*http.Transport, error) {
	cp.RLock()
	defer cp.RUnlock()
	if len(cp.transports) < 1 {
		return nil, ErrEmptyTransportPool
	}
	if len(cp.transports) < id {
		return nil, ErrNotFoundTransportPoolElement
	}
	return &cp.transports[id], nil
}

func (cp *ClientsPool)GetTransportByProxyUrl(url string) (*http.Transport, error) {
	if len(url) == 0 {
		url = "direct"
	}
	cp.RLock()
	defer cp.RUnlock()
	if id, ok := cp.transportsMapByProxyUrl[url]; !ok {
		return nil, ErrNotFoundTransportPoolElement
	} else {
		return cp.GetTransportById(id)
	}
}


func (cp *ClientsPool)ClientByName(r *http.Request, clientName string, proxy...string) (rsp *http.Response, err error) {
	var client *Client
	client, err = cp.GetClientByName(clientName, proxy...)
	if err != nil {
		return
	}

	if r.Header == nil {
		r.Header = http.Header{}
	}

	if len(r.UserAgent()) == 0 {
		r.Header.Add("User-Agent", cp.defaultUa)
	}

	return client.Do(r)
}

func (cp *ClientsPool)NewClientAutoCookie(proxy... string) (hc *Client, err error) {
	var transport *http.Transport
	if len(proxy) < 1 {
		proxy = []string{""}
	}
	transport, err = cp.GetTransportByProxyUrl(proxy[0])
	if err != nil {
		return nil, err
	}
	var jar *cookiejar.Jar
	jar, err =  cookiejar.New(&cookiejar.Options{})
	header := make(http.Header)
	header.Add("User-Agent", cp.defaultUa)
	return &Client{
		Client: &http.Client{
			Transport:     transport,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Jar:           jar,
			Timeout:       cp.defaultConnectTimeout,
		},
		Header: header,
	}, nil
}


func (cp *ClientsPool)GetClientByName(clientName string, proxy...string) (hc *Client, err error) {
	if clientId, ok := cp.clientsMapByClientName[clientName]; ok {
		hc = &cp.clients[clientId]
	} else {
		hc, err = cp.NewClientAutoCookie(proxy...)
		if err != nil {
			return
		}
		cp.Lock()
		cp.clients = append(cp.clients, *hc)
		cp.clientsMapByClientName[clientName] = len(cp.clients) - 1
		cp.Unlock()
	}
	return
}