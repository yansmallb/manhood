package api

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"net/http"
)

type handler func(w http.ResponseWriter, r *http.Request)

var routes = map[string]map[string]handler{
	"GET": {
		"/":		getHelp,
		"/namespaces/{namespace:.*}/rcs/{rc:.*}/metrics/{metric:.*}/get":              getRcMetric,
		"/namespaces/{namespace:.*}/rcs/{rc:.*}/get":              getRc,
	},
	"POST": {
		"/namespaces/{namespace:.*}/rcs/{rc:.*}/post":              postRc,
	},
	"PUT": {
		"/namespaces/{namespace:.*}/rcs/{rc:.*}/put":              postRc,
	},
}

// NewPrimary creates a new API router.
func NewPrimary() *mux.Router {
	// Register the API events handler in the cluster.
	r := mux.NewRouter()
	for method, mappings := range routes {
		for route, fct := range mappings {
			//log.WithFields(log.Fields{"method": method, "route": route}).Debug("api.NewPrimary():Registering HTTP route")

			localRoute := route
			localFct := fct
			wrap := func(w http.ResponseWriter, r *http.Request) {
				log.WithFields(log.Fields{"method": r.Method, "uri": r.RequestURI}).Debug("api.NewPrimary():HTTP request received")
				localFct(w, r)
			}
			localMethod := method

			r.Path(localRoute).Methods(localMethod).HandlerFunc(wrap)
		}
	}
	return r
}
