package common

import (
    "endcoding/json"
    "fmt"
)

type zbxDiscovery struct {
    data    []map[interface{}]interface{}
}

func newZbxDiscovery() (*zbxDiscovery) {
    return &zbxDiscovery{
        data: make([]map[interface{}]interface{}),
    }
}

func (zd *zbxDiscovery) addItem(m map[interface{}]interface{}) {
    zd.data = append(zd.data, m)
}

func (zd *zbxDiscovery) fmt() (string, error) {
    data := make(map[interface{}]interface{})
    data["data"] = zd.data
    return json.Marshal(data)
}

func (ep *Earpice) discoveryCluster(args ...string) (error) {
    zd := newZbxDiscovery()

    for _, val := range ep.clusterInfo.ListInfo() {
        zd.addItem(map[interface{}]interface{}{
            "{#CLUSTER}": val.Name,
            "{#IP}": val.IP,
        })
    }

    res, err := zd.fmt()
    if err != nil {
        return err
    }

    fmt.Print(res)
}

func (ep *Earpice) discoveryNamespace(args ...string) (error) {
    zd := newZbxDiscovery()

    for _, val := range ep.clusterInfo.ListInfo() {
        
    }
}