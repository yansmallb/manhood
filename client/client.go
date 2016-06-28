package client

import (
	log "github.com/Sirupsen/logrus"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/resource"
	"github.com/yansmallb/manhood/client/kube"
	"github.com/yansmallb/manhood/client/heapster"
)

func GetReplicationControllersMetrio(namespace,rcname,metric string) ( heapster.HeapsterPodsMetric,error) {
	log.Infoln("client.GetReplicationControllersMetrio: -------begin------")
	kc,err := kube.NewKubeClient(*KubeConfigPath)
	if err!= nil {
		log.Errorf("client.GetReplicationControllersMetrio: %v\n", err)
		return nil, err
	}
	rc,err := kc.GetReplicationController(namespace,rcname) 
	if err!= nil {
		log.Errorf("client.GetReplicationControllersMetrio: %v\n", err)
		return nil, err
	}
	
	pods,err := kc.GetPodsForRc(rc)
	if err!=nil {
		log.Errorf("client.GetReplicationControllersMetrio: get pods error :%v\n", err)
	}
	log.Debugf("client.GetReplicationControllersMetrio: pods :%v\n", pods)

	hc := heapster.NewHeapsterClient(*HeapsterUrl)
	hpm := make(heapster.HeapsterPodsMetric, 0)
	for _,pod := range pods {
		res, err := hc.GetMetric(namespace,"pods",string(pod.Name),metric)
		if err!= nil {
			log.Errorf("client.GetReplicationControllersMetrio: get metric error :%v\n", err)
		}
		log.Debugf("client.GetReplicationControllersMetrio: hcPods metric :%v\n", res)
		hpm = append(hpm, *res)
	}
	return hpm, nil
}

func GetReplicationController(namespace,rcname string) (*api.ReplicationController, error) {
	log.Infoln("client.GetReplicationController: -------begin------")
	kc, err := kube.NewKubeClient(*KubeConfigPath)
	if err!= nil {
		log.Errorf("client.GetReplicationController: %v\n", err)
		return nil, err
	}
	return kc.GetReplicationController(namespace,rcname)
}

func UpdateReplicationController(namespace,rcname string, replicas int, cpu,memory string) (*api.ReplicationController, error) {
	log.Infoln("client.UpdateReplicationController: -------begin------")
	kc,err := kube.NewKubeClient(*KubeConfigPath)
	if err!= nil {
		log.Errorf("client.UpdateReplicationController: %v\n", err)
		return nil, err
	}
	rc,err := kc.GetReplicationController(namespace,rcname) 
	if err!= nil {
		log.Errorf("client.UpdateReplicationController: %v\n", err)
		return nil, err
	}

	rc.Spec.Replicas = replicas

	rs := make(map[api.ResourceName]resource.Quantity)
	rcpu, err :=resource.ParseQuantity(cpu)
	if err != nil {
		log.Errorf("client.UpdateReplicationController: %v\n", err)
	}else{
		rs[api.ResourceCPU] = *rcpu
	}
	rmem, err := resource.ParseQuantity(memory)
	if err != nil {
		log.Errorf("client.UpdateReplicationController: %v\n", err)
	}else{
		rs[api.ResourceMemory] = *rmem
	}
	
	for index := range rc.Spec.Template.Spec.Containers {
		rc.Spec.Template.Spec.Containers[index].Resources.Limits = rs
	}
	return  kc.UpdateReplicationController(namespace, rc)
}