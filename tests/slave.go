package main

import (
    	"fmt"
    	"io/ioutil"
    	"net"
    	"os"
    	"time"
)


func main() {
	if len(os.Args) != 2 {
    		fmt.Fprintf(os.Stderr, "Usage: %s serviceName \n", os.Args[0])
    		os.Exit(1)
	}
	serviceName := os.Args[1]
	
	fmt.Println("-----------start manhoodtest slave-----------")

	for{
		time.Sleep(3*time.Second)

		ns, err := net.LookupHost("ep-"+serviceName)  
		if err != nil {
        			fmt.Printf("Fatal error: %s\n", err)
        			continue
    		}
		  
		service := ns[0] + ":2372"
		fmt.Printf("Info service: %s\n", service)
		tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
		if err != nil {
        			fmt.Printf("Fatal error: %s\n", err)
        			continue
    		}
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
        			fmt.Printf("Fatal error: %s\n", err)
        			continue
    		}
		_, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
		if err != nil {
        			fmt.Printf("Fatal error: %s\n", err)
        			continue
    		}
		result, err := ioutil.ReadAll(conn)
		if err != nil {
        			fmt.Printf("Fatal error: %s\n", err)
        			continue
    		}
		fmt.Println(string(result))
    	}
}
