package common

import (
    "fmt"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

type Earpice struct {
    kubeClients     map[string]*kubeClient.Clientset
    clusterInfo     *ClusterInfo
}

func NewEarpice(c *ClusterInfo) (*Earpice, error) {
    if c == nil {
        return nil, fmt.Errorf("the cluster info is invalid")
    }

    return &Earpice{
        clusterInfo: c,
    }, nil
}

func (ep *Earpice) GetClientset(name string) (*kubeClient.Clientset, error) {
    c, ok := ep.kubeClients[name]
    if ok {
        return c
    }

    cInfo, err := ep.clusterInfo.GetInfo(name)
    if err != nil {
        return nil, err
    }

    config, err := clientcmd.BuildConfigFromFlags("", cInfo.kubecfg)
    if err != nil {
        return nil, err
    }

    c, err = kubernetes.NewForConfig(config)
    if err != nil {
        return nil, err
    }

    ep.kubeClients[name] = c
    return c, nil
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
    case "discovery_cs":
        err = ep.discoveryComponentstatuses(ops.Args...)

    case "ev":
        err = ep.clusterEvent(ops.Args...)
    case "cs":
        err = ep.componentstatuses(ops.Args...)
    case "pod":
        err = ep.pod(ops.Args...)
    case "no":
        err = ep.node(ops.Args...)

    }

    if err != nil {
        return err
    }
}