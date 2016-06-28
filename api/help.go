package api

var helpString = `
  Usage: manhood COMMAND [args...]
  Version: 0.2.0
  Author:yansmallb
  Email:yanxiaoben@iie.ac.cn

  Commands:
      manage    -kubeconfig=[./kubeconfig] -heapster=[http://127.0.0.1:8082] -host=[127.0.0.1] -port=[2371]
      help

  RestApi
  "GET": {
    "/":                getHelp,  
    "/namespaces/{namespace}/rcs/{rc}/metrics/{metric}":      getRcMetric,
    "/namespaces/{namespace}/rcs/{rc}/get":                 getRc,
  },
  "POST"/"PUT": {
    "/namespaces/{namespace}/rcs/{rc}/{post/put}?replicas={int}&cpu={string}&memory={string}":             postRc,
  },

  Note:
      metric type look here : https://github.com/kubernetes/heapster/blob/master/docs/storage-schema.md
  ` 

func Help() string{
      return helpString
}
