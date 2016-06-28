package main

import (
    	"fmt"
    	"net"
    	"net/http"
    	"os"
    	"time"
    	"encoding/json"
    	"io/ioutil"
    	"strconv"
    	"bytes"

    	"k8s.io/kubernetes/pkg/api"
    	"github.com/yansmallb/manhood/client/heapster"
)

func main() {
	if len(os.Args) != 2 {
    		fmt.Fprintf(os.Stderr, "Usage: %s manhoodUrl", os.Args[0])
    		os.Exit(1)
	}
	manhoodUrl := os.Args[1]

    	service := ":2372"
    	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
    	checkError(err)
    	listener, err := net.ListenTCP("tcp", tcpAddr)
    	checkError(err)

    	fmt.Println("-----------start manhoodtest master-----------")

    	go manhoodTest(manhoodUrl)

    	for {
        		conn, err := listener.Accept()
        		if err != nil {
        			fmt.Printf("Error: %s\n", err)
            			continue
        		}
        		daytime := time.Now()
        		conn.Write([]byte(daytime.String())) 	// don't care about return value
        		conn.Close()               	 	// we're finished with this client
    	}
}

func manhoodTest(manhoodUrl string){
	for{
		//get rc
		rc := getRc(manhoodUrl)
		fmt.Printf("rc:%v\n",rc)

		//get rcmetric
		rcmetric := getRcMetric(manhoodUrl)
		fmt.Printf("rcm:%v\n",rcmetric)

		//post rc
		if rcmetric!=nil && len(*rcmetric)<3 {
			replicas := strconv.Itoa(len(*rcmetric)+1)
			cpu := "0"
			if len((*rcmetric)[0].Metrics) > 0 {
				if (*rcmetric)[0].Metrics[0].Value > 300 {
					limitcpu := (*rcmetric)[0].Metrics[0].Value / 2
					cpu = strconv.Itoa(int(limitcpu))  + "m"
				}
			}

			rc := postRc(manhoodUrl,replicas,cpu)
			fmt.Printf("rcpost:%v\n",rc)
		}
		time.Sleep(5*time.Second)
	}
}

func getRc(manhoodUrl string) *api.ReplicationController {
	resp,err := http.Get(manhoodUrl+"/namespaces/default/rcs/mhslave/get")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return nil
	}
	result, _ := ioutil.ReadAll(resp.Body) 
  	defer resp.Body.Close() 
  	rc := &api.ReplicationController{}
  	err = json.Unmarshal(result, rc)
  	if err != nil {
  		fmt.Println("Fatal error: %s\n", err)
		return nil
  	}
  	return rc
}

func getRcMetric(manhoodUrl string) *heapster.HeapsterPodsMetric {
	resp,err := http.Get(manhoodUrl+"/namespaces/default/rcs/mhslave/metrics/cpu/limit/get")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return nil
	}
	result, _ := ioutil.ReadAll(resp.Body) 
  	defer resp.Body.Close() 
  	rcm := &heapster.HeapsterPodsMetric{}
  	err = json.Unmarshal(result,rcm)
  	if err != nil {
  		fmt.Println("Fatal error: %s\n", err)
		return nil
  	}
  	return rcm
}

func postRc(manhoodUrl string,replicas string,cpulimit string) *api.ReplicationController {
	body := bytes.NewBuffer([]byte(""))
	manhoodUrl = manhoodUrl+"/namespaces/default/rcs/mhslave/post?replicas="+replicas+"&cpu="+cpulimit
	resp,err := http.Post(manhoodUrl, "application/json;charset=utf-8", body)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return nil
	}
	result, _ := ioutil.ReadAll(resp.Body) 
  	defer resp.Body.Close() 
  	rc := &api.ReplicationController{}
  	err = json.Unmarshal(result, rc)
  	if err != nil {
  		fmt.Println("Fatal error: %s\n", err)
		return nil
  	}
  	return rc
}

func checkError(err error) {
    	if err != nil {
        		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err)
        		os.Exit(1)
    	}
}