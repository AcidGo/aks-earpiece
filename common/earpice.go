package common

import (
    "fmt"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

type Earpice struct {
    kubeClient      *kubeClient.Clientset
    clusterInfo     *ClusterInfo
}

func NewEarpice(c *ClusterInfo) (*Earpice, error) {
    if c == nil || c.Name == "" || c.kubecfg == "" {
        return nil, fmt.Errorf("the cluster info is invalid")
    }

    config, err := clientcmd.BuildConfigFromFlags("", c.kubecfg)
    if err != nil {
        return nil, err
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return nil, err
    }

    return &Earpice{
        kubeClient:     clientset,
        clusterInfo:    c,
    }, nil
}

func (ep *Earpice) Call(ops *Options) (error) {
    var err error

    switch ops.Method {
    case "discovery_cl":
        err = ep.discoveryCluster(ops.Args...)
    case "discovery_ns":
        err = ep.discoveryNamespace(ops.Args...)
    case "discovery_no":
        err = ep.discoveryNode(ops.Args...)
    case "discovery_pod":
        err = ep.discoveryPod(ops.Args...)

    case "ev":
        err = ep.clusterEvent(ops.Args...)
    case "cl_cs":
        err = ep.clusterComponentstatuses(ops.Args...)
    case "pod":
        err = ep.pod(ops.Args...)
    case "no":
        err = ep.node(ops.Args...)

    }

    if err != nil {
        return err
    }
}