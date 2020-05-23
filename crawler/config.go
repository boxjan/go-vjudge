// Author: Boxjan
// Datetime: 2020/3/19 19:01
// I don't know why, but it can work

package crawler

import (
	"boxjan.li/go-vjudge/crawler/httpclient"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type configStruct struct {
	Timeout uint `json:"timeout"`
	Ua string `json:"user_agent"`
	Proxies []string `json:"proxies"`
	RemoteFileSave string `json:"remote_file_save"`
	Remote []remoteConfig `json:"remote"`
	RemoteMapByName map[string]*remoteConfig
	LocalFilePath string `json:"oj_file_path"`
}

type remoteConfig struct {
	Name string `json:"name"`
	Accounts []remoteAuth `json:"accounts"`
	Proxies []string `json:"proxies"`
	Ua string `json:"user_agent"`
}

type remoteAuth struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
	Token string `json:"token"`
}

func (h *crawler)loadConfig()  {
	log.Println("start load config")
	defer log.Println("end load config")

	if configFile, err := os.Open(h.configPath); err != nil {
		log.Println("open config failed with:", err)
	} else {
		if configFileContext, err := ioutil.ReadAll(configFile); err != nil {
			log.Println("read config failed with:", err)
		} else {
			newConfig := configStruct{RemoteMapByName: make(map[string]*remoteConfig)}
			if err := json.Unmarshal(configFileContext, &(newConfig)); err != nil {
				log.Println("unmarshal json failed with:", err)
			} else {
				for i, v := range newConfig.Remote {
					newConfig.RemoteMapByName[v.Name] = &newConfig.Remote[i]
				}
				for oJName := range h.workerHandle {
					if _, ok := newConfig.RemoteMapByName[oJName]; !ok {
						log.Println("can not find", oJName, "config in config file")
					}
				}
				h.Config = newConfig
				if err := os.MkdirAll(newConfig.LocalFilePath, 0755); err != nil {
					panic(err)
				}
			}
		}
	}
	return
}

func (h *crawler) loadTransports() {
	if h.clients != nil {
		return
	}
	if len(h.Config.Proxies) == 0 {
		h.Config.Proxies = append(h.Config.Proxies, "direct")
	}

	proxyURLs := make([]string, 0, len(h.Config.Proxies) * 2)
	proxyURLs = append(proxyURLs, h.Config.Proxies...)

	for _, remote := range h.Config.Remote {
		proxyURLs = append(proxyURLs, remote.Proxies...)
	}

	clients, errs := httpclient.NewClientPools(time.Millisecond * time.Duration(h.Config.Timeout), h.Config.Ua, proxyURLs...)
	for _, err := range errs {
		log.Println(err)
	}
	h.clients = clients
}

func (h *crawler) reload() {
	h.loadConfig()
	h.loadTransports()
}