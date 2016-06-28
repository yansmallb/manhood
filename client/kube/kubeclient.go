package kube

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/fields"
	kclient "k8s.io/kubernetes/pkg/client/unversioned"
	kclientcmd "k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
	kclientcmdapi "k8s.io/kubernetes/pkg/client/unversioned/clientcmd/api"
)

type KubeClient struct {
	client *kclient.Client
}

func NewKubeClient(kubeConfigPath string)(*KubeClient, error){
	c, err := kclientcmd.LoadFromFile(kubeConfigPath)
	if err != nil {
		return  nil, fmt.Errorf("error loading kubeConfig: %v", err.Error())
	}
	if c.CurrentContext == "" || len(c.Clusters) == 0 {
		return  nil, fmt.Errorf("invalid kubeConfig: %+v", *c)
	}

	config, err := kclientcmd.NewDefaultClientConfig(
		*c,
		&kclientcmd.ConfigOverrides{
			ClusterInfo: kclientcmdapi.Cluster{
				APIVersion: "v1",
			},
		}).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("error parsing kubeConfig: %v", err.Error())
	}
	client, _  := kclient.New(config)
	kc := &KubeClient{
		client: client, 
	}
	if err != nil {
		return  nil, fmt.Errorf("error creating client - %q", err)
	}
	return kc, nil
}

func (kc *KubeClient)  GetPodsForRc(controller *api.ReplicationController) ([]api.Pod, error) {
	podSpace := kc.client.Pods(controller.Namespace)

	selector := labels.Set(controller.Spec.Selector).AsSelector()
	podList,_ := podSpace.List(api.ListOptions{
		LabelSelector: selector,
		FieldSelector: fields.Everything(),
	})

	pods := make([]api.Pod,0)
	for _, pod := range podList.Items {
		if isControllerMatch(&pod,controller) {
			pods = append(pods,pod)
		}
	}
	return pods,nil
}

func (kc *KubeClient)  GetReplicationController(namespace,rcname string) (*api.ReplicationController, error) {
	rcs := kc.client.ReplicationControllers(namespace)
	rc, err := rcs.Get(rcname)
	if err!=nil {
		log.Errorf("kube.GetReplicationController: get rc error :%v", err)
	}
	log.Debugf("kube.GetReplicationController: rc :%v", rc)
	return rc, err 
}

func (kc *KubeClient)  UpdateReplicationController(namespace string ,controller *api.ReplicationController) (*api.ReplicationController, error) {
	rcs := kc.client.ReplicationControllers(namespace)
	return rcs.Update(controller)
}

func isControllerMatch(pod *api.Pod, rc *api.ReplicationController) bool {
	if rc.Namespace != pod.Namespace {
		return false
	}
	labelSet := labels.Set(rc.Spec.Selector)
	selector := labels.Set(rc.Spec.Selector).AsSelector()

	// If an rc with a nil or empty selector creeps in, it should match nothing, not everything.
	if labelSet.AsSelector().Empty() || !selector.Matches(labels.Set(pod.Labels)) {
		return false
	}
	return true
}