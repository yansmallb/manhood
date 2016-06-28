package client

import (
	"flag"
)

var (
	KubeConfigPath = flag.String("kubeconfig","./kubeconfig.yaml","path to kubeconfig file.")
	ManhoodHost = flag.String("host","127.0.0.1","host to listen.")
	ManhoodPort = flag.Int("port",2371,"port to listen.")
	HeapsterUrl = flag.String("heapster","http://127.0.0.1:8082","heapster url.")
)