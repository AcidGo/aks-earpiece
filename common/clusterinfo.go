package common

import (
    "fmt"
)

type Cluster struct {
    // custom cluster name
    Name        string
    // fit zabbix IP suffix
    IP          string
    // kube config path
    Kubecfg     string
}

type ClusterInfo struct {
    set     map[string]*Cluster
}

func NewClusterInfo() (*ClusterInfo, error) {
    return &ClusterInfo{
        set: make(map[string]*Cluster),
    }, nil
}

func (ci *ClusterInfo) AddInfo(c *Cluster) (error) {
    ci.set[c.Name] = c
    return nil
}

func (ci *ClusterInfo) GetInfo(name string) (*Cluster, error) {
    c, ok := ci.set[name]
    if !ok {
        return nil, fmt.Errorf("not found the cluster %s", name)
    }

    return c, nil
}

func (ci *ClusterInfo) ListInfo() ([]*Cluster) {
    res := make([]*Cluster, len(ci.set))
    idx := 0
    for _, val := range ci.set {
        res[idx] = val
        idx++
    }

    return res
}