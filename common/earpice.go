package common

import (
    "fmt"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

type Earpiece struct {
    kubeClients     map[string]*kubernetes.Clientset
    clusterInfo     *ClusterInfo
}

func NewEarpice(c *ClusterInfo) (*Earpiece, error) {
    if c == nil {
        return nil, fmt.Errorf("the cluster info is invalid")
    }

    return &Earpiece{
        clusterInfo: c,
    }, nil
}

func (ep *Earpiece) GetClientset(name string) (*kubernetes.Clientset, error) {
    c, ok := ep.kubeClients[name]
    if ok {
        return c, nil
    }

    cInfo, err := ep.clusterInfo.GetInfo(name)
    if err != nil {
        return nil, err
    }

    config, err := clientcmd.BuildConfigFromFlags("", cInfo.Kubecfg)
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

func (ep *Earpiece) Call(ops *Options) (error) {
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

    // case "ev":
    //     err = ep.clusterEvent(ops.Args...)
    case "cs":
        err = ep.componentstatuses(ops.Args...)
    // case "pod":
    //     err = ep.pod(ops.Args...)
    // case "no":
    //     err = ep.node(ops.Args...)

    }

    if err != nil {
        return err
    }

    return nil
}