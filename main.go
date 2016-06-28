package main

import (
	"os"
	"fmt"
	"flag"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/yansmallb/manhood/client"
	"github.com/yansmallb/manhood/api"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println(api.Help())
		return
	}
	flag.Parse()
	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)
	log.Infof("main.main(): KubeConfigPath=%v,HeapsterUrl=%v",*client.KubeConfigPath,*client.HeapsterUrl)
	log.Infof("main.main(): ManhoodHost=%v,ManhoodPort=%v", *client.ManhoodHost,*client.ManhoodPort)
	log.Infoln("main.main():Start Manage")

	// start API listener
	host := *client.ManhoodHost + ":" + strconv.Itoa(*client.ManhoodPort)
	if !strings.Contains(host,"http://") {
		host = "http://" + host
	}
	hosts := []string{host}
	log.Infof("main: hosts %+v\n", hosts)

	server := api.NewServer(hosts)
	server.SetHandler(api.NewPrimary())
	server.ListenAndServe()
}
