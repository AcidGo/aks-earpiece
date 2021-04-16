package common

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

func (ci *ClusterInfo) AddInfo(c *Cluster) (error) {
    set[c.Name] = c
    return nil
}

func (ci *ClusterInfo) GetInfo(name string) (*ClusterInfo, error) {
    c, ok := ci.set[name]
    if !ok {
        return nil, fmt.Errorf("not found the cluster %s", name)
    }

    return c, nil
}

func (ci *ClusterInfo) ListInfo() ([]*ClusterInfo) {
    res := make([]*ClusterInfo, len(ci.set))
    idx := 0
    for _, val := range ci.set {
        res[idx] = val
        idx++
    }

    return res
}