package api

import (
	"strconv"
	"net/http"
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/yansmallb/manhood/client"
)

func getHelp(w http.ResponseWriter, r *http.Request){
	log.Debugln("api.getHelp")
	help := Help()
	json.NewEncoder(w).Encode(help)
}

func getRc(w http.ResponseWriter, r *http.Request){
	namespace := mux.Vars(r)["namespace"]
	rcname := mux.Vars(r)["rc"]
	log.Debugf("api.getRc:%v,%v",namespace,rcname)

	rc,err := client.GetReplicationController(namespace,rcname)
	if err!=nil{
		json.NewEncoder(w).Encode(err)
		return 
	}
	json.NewEncoder(w).Encode(rc)
}

func getRcMetric(w http.ResponseWriter, r *http.Request){
	namespace := mux.Vars(r)["namespace"]
	rcname := mux.Vars(r)["rc"]
	metric := mux.Vars(r)["metric"]
	log.Debugf("api.getRcMetric: namespace=%v,rcname=%v,metric=%v",namespace,rcname,metric)

	rcm, _ := client.GetReplicationControllersMetrio(namespace,rcname,metric)
	json.NewEncoder(w).Encode(rcm)
}

func postRc(w http.ResponseWriter, r *http.Request){
	namespace := mux.Vars(r)["namespace"]
	rcname := mux.Vars(r)["rc"]
	log.Debugf("api.postRc: namespace=%v,rcname=%v",namespace,rcname)
	
	r.ParseForm();
	replicas,err := strconv.Atoi(r.Form.Get("replicas"))
	if err!= nil {
		json.NewEncoder(w).Encode(err)
		return 
	}
	cpu := r.Form.Get("cpu")
	memory := r.Form.Get("memory")
	log.Debugf("api.postRc: cpu=%v,memory=%v",namespace,rcname)

	rc,err := client.UpdateReplicationController(namespace,rcname,replicas,cpu,memory)
	if err!=nil {
		json.NewEncoder(w).Encode(err)
		return 
	}
	json.NewEncoder(w).Encode(rc)
}
