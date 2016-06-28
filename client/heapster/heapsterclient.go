package heapster

import (
	"net/http"
	"fmt"
	"encoding/json"
	"time"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

type MetricPoint struct {
	Timestamp	time.Time 	`json:"timestamp"`
	Value		uint64		`json:"value"`
	FloatValue	*float64 	`json:"floatValue,omitempty"`
}

type HeapsterPodMetric struct {
	Metrics 	[]MetricPoint 	`json:"metrics"`
}

type HeapsterPodsMetric  []HeapsterPodMetric

type HeapsterClient struct{
	url string
	version string
	namespace string
	level	string	//pod,container,node
	id   string 	//podid,containerid,nodeid
	metric  string	//https://github.com/kubernetes/heapster/blob/master/docs/storage-schema.md
}

func NewHeapsterClient(url string) *HeapsterClient {
	return &HeapsterClient{
		url: url,
		version: "/api/v1/model",
	}
}

func (hc *HeapsterClient) SetNamespace(ns string) *HeapsterClient {
	if hc != nil {
		hc.namespace = ns
	}
	return hc
}

func (hc *HeapsterClient) SetLevel(l string) *HeapsterClient {
	if hc != nil {
		hc.level = l
	}
	return hc
}

func (hc *HeapsterClient) SetMetric(m string) *HeapsterClient {
	if hc != nil {
		hc.metric = m
	}
	return hc
}

func (hc *HeapsterClient) SetId(id string) *HeapsterClient {
	if hc != nil {
		hc.id =id
	}
	return hc
}

func (hc *HeapsterClient) GetMetric(namespace,level,id,metric string) ( *HeapsterPodMetric, error) {
	return hc.SetNamespace(namespace).SetLevel(level).SetId(id).SetMetric(metric).Get()
}

func (hc *HeapsterClient) Url() (string,error){
	url := hc.url + hc.version
 	if hc.level != "node" {
 		if hc.namespace != "" {
 			url += "/namespaces/" + hc.namespace
 		}else{
 			return "", fmt.Errorf("heapsterclient namespace is empty")
 		}
 	}
 	url += "/"+ hc.level + "/" +hc.id
 	if hc.metric != ""{
 		url += "/metrics/" + hc.metric
 	}
 	return url,nil
}

func (hc *HeapsterClient) Get() (*HeapsterPodMetric,  error) {
	url ,err := hc.Url()
	if err != nil {
		return nil, err 
	}
	log.Debugf("heapsterclient.Get: url: %v", url)

	resp, err  := http.Get(url)
	result, err := ioutil.ReadAll(resp.Body) 
  	defer resp.Body.Close() 
  	metric := &HeapsterPodMetric{}
  	err = json.Unmarshal(result, metric)
  	log.Debugf("heapsterclient.Get: metric: %v", metric)

  	return metric,err
}