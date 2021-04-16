package common

import (
    "context"
    "encoding/json"
    "fmt"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type zbxDiscovery struct {
    data    []map[interface{}]interface{}
}

func newZbxDiscovery() (*zbxDiscovery) {
    return &zbxDiscovery{
        data: make([]map[interface{}]interface{}, 0),
    }
}

func (zd *zbxDiscovery) addItem(m map[interface{}]interface{}) {
    zd.data = append(zd.data, m)
}

func (zd *zbxDiscovery) fmt() (string, error) {
    data := make(map[interface{}]interface{})
    data["data"] = zd.data

    b, err := json.Marshal(data)
    if err != nil {
        return "", err
    }

    return string(b), nil
}

func (ep *Earpiece) discoveryCluster(args ...string) (error) {
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

    return nil
}

func (ep *Earpiece) discoveryNamespace(args ...string) (error) {
    zd := newZbxDiscovery()

    for _, val := range ep.clusterInfo.ListInfo() {
        clientset, err := ep.GetClientset(val.Name)
        if err != nil {
            continue
        }

        nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
        if err != nil {
            continue
        }

        for _, n := range nodes.Items {
            zd.addItem(map[interface{}]interface{}{
                "{#CLUSTER}": val.Name,
                "{#NODE}": n.Name,
            })
        }
    }

    res, err := zd.fmt()
    if err != nil {
        return err
    }

    fmt.Print(res)

    return nil
}

func (ep *Earpiece) discoveryNode(args ...string) (error) {
    zd := newZbxDiscovery()

    for _, val := range ep.clusterInfo.ListInfo() {
        clientset, err := ep.GetClientset(val.Name)
        if err != nil {
            continue
        }

        namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
        if err != nil {
            continue
        }

        for _, n := range namespaces.Items {
            zd.addItem(map[interface{}]interface{}{
                "{#CLUSTER}": val.Name,
                "{#NAMESPACE}": n.Name,
                "{#IP}": val.IP,
            })
        }
    }

    res, err := zd.fmt()
    if err != nil {
        return err
    }

    fmt.Print(res)

    return nil
}

func (ep *Earpiece) discoveryPod(args ...string) (error) {
    if len(args) < 2 {
        return fmt.Errorf("must input [cluster, namespace] in args")
    }

    clusterName := args[0]
    namespaceName := args[1]

    cInfo, err := ep.clusterInfo.GetInfo(clusterName)
    if err != nil {
        return err
    }

    clientset, err := ep.GetClientset(cInfo.Name)
    if err != nil {
        return err
    }

    zd := newZbxDiscovery()

    pods, err := clientset.CoreV1().Pods(namespaceName).List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        return err
    }

    for _, pod := range pods.Items {
        zd.addItem(map[interface{}]interface{}{
            "{#POD}": pod.Name,
        })
    }

    res, err := zd.fmt()
    if err != nil {
        return err
    }

    fmt.Print(res)

    return nil
}

func (ep *Earpiece) discoveryComponentstatuses(args ...string) (error) {
    if len(args) < 1 {
        return fmt.Errorf("must input [cluster] in args")
    }

    clusterName := args[0]

    cInfo, err := ep.clusterInfo.GetInfo(clusterName)
    if err != nil {
        return err
    }

    clientset, err := ep.GetClientset(cInfo.Name)
    if err != nil {
        return err
    }

    zd := newZbxDiscovery()

    css, err := clientset.CoreV1().ComponentStatuses().List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        return err
    }

    for _, cs := range css.Items {
        zd.addItem(map[interface{}]interface{}{
            "{#CLUSTER}": clusterName,
            "{#CS}": cs.Name,
        })
    }

    return nil
}

func (ep *Earpiece) componentstatuses(args ...string) (error) {
    if len(args) < 2 {
        return fmt.Errorf("must input [cluster, cs] in args")
    }

    clusterName := args[0]
    csName := args[1]

    cInfo, err := ep.clusterInfo.GetInfo(clusterName)
    if err != nil {
        return err
    }

    clientset, err := ep.GetClientset(cInfo.Name)
    if err != nil {
        return err
    }

    cs, err := clientset.CoreV1().ComponentStatuses().Get(context.TODO(), csName, metav1.GetOptions{})
    if err != nil {
        return err
    }

    m := make(map[string]string)

    m["status"] = string(cs.Conditions[len(cs.Conditions)-1].Type)
    m["error"] = cs.Conditions[len(cs.Conditions)-1].Error

    res, err := json.Marshal(m)
    if err != nil {
        return err
    }

    fmt.Print(res)

    return nil
}