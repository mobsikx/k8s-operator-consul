// package consul_test

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ConsulDatacenters []string
type ConsulInfo []struct {
	ID              string `json:"ID"`
	Node            string `json:"Node"`
	Address         string `json:"Address"`
	Datacenter      string `json:"Datacenter"`
	TaggedAddresses struct {
		Lan string `json:"lan"`
		Wan string `json:"wan"`
	} `json:"TaggedAddresses"`
	NodeMeta struct {
		ConsulNetworkSegment string `json:"consul-network-segment"`
	} `json:"NodeMeta"`
	ServiceKind    string        `json:"ServiceKind"`
	ServiceID      string        `json:"ServiceID"`
	ServiceName    string        `json:"ServiceName"`
	ServiceTags    []interface{} `json:"ServiceTags"`
	ServiceAddress string        `json:"ServiceAddress"`
	ServiceWeights struct {
		Passing int `json:"Passing"`
		Warning int `json:"Warning"`
	} `json:"ServiceWeights"`
	ServiceMeta struct {
	} `json:"ServiceMeta"`
	ServicePort              int    `json:"ServicePort"`
	ServiceEnableTagOverride bool   `json:"ServiceEnableTagOverride"`
	ServiceProxyDestination  string `json:"ServiceProxyDestination"`
	ServiceProxy             struct {
	} `json:"ServiceProxy"`
	ServiceConnect struct {
	} `json:"ServiceConnect"`
	CreateIndex int `json:"CreateIndex"`
	ModifyIndex int `json:"ModifyIndex"`
}

func getConsulDCList() ConsulDatacenters {
	url := "http://consul.service.consul/v1/catalog/datacenters"
	var contents []byte
	var result ConsulDatacenters

	timeout := time.Duration(3 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	contents, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(contents, &result)
	if err != nil {
		panic(err)
	}

	return result
}

func getConsulDCInfo(dc string) ConsulInfo {
	var url string
	if dc == "" {
		url = "http://consul.service.consul/v1/catalog/service/consul"
	} else {
		url = "http://consul.service." + dc + ".consul/v1/catalog/service/consul"
	}
	var contents []byte
	var result ConsulInfo

	timeout := time.Duration(3 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	contents, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(contents, &result)
	if err != nil {
		return nil
	}

	return result
}

func getConsulDCAddress(ci ConsulInfo) string {
	if ci == nil {
		return ""
	}
	return ci[0].Address
}

func getConsulDCName(ci ConsulInfo) string {
	if ci == nil {
		return ""
	}
	return ci[0].Datacenter
}

func main() {
	// get list of DCs
	dcList := getConsulDCList()

	// get info about one particular DC
	for _, dc := range dcList {
		ci := getConsulDCInfo(dc)
		fmt.Println(getConsulDCAddress(ci))
		fmt.Println(getConsulDCName(ci))
	}
}

