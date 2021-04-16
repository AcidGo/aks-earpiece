package common

import (
    "context"
    "encoding/json"
    "fmt"
    "strings"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    corev1 "k8s.io/api/core/v1"
)

type zbxDiscovery struct {
    data    []map[string]interface{}
}

func newZbxDiscovery() (*zbxDiscovery) {
    return &zbxDiscovery{
        data: make([]map[string]interface{}, 0),
    }
}

func (zd *zbxDiscovery) addItem(m map[string]interface{}) {
    zd.data = append(zd.data, m)
}

func (zd *zbxDiscovery) fmt() (string, error) {
    data := make(map[string]interface{})
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
        zd.addItem(map[string]interface{}{
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

func (ep *Earpiece) discoveryNode(args ...string) (error) {
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
            zd.addItem(map[string]interface{}{
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

func (ep *Earpiece) discoveryNamespace(args ...string) (error) {
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
            zd.addItem(map[string]interface{}{
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
        zd.addItem(map[string]interface{}{
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
        zd.addItem(map[string]interface{}{
            "{#CLUSTER}": clusterName,
            "{#CS}": cs.Name,
        })
    }

    res, err := zd.fmt()
    if err != nil {
        return err
    }

    fmt.Print(res)

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

    b, err := json.Marshal(m)
    if err != nil {
        return err
    }

    fmt.Print(string(b))

    return nil
}

func (ep *Earpiece) pod(args ...string) (error) {
    if len(args) < 3 {
        return fmt.Errorf("must input [cluster, namespace, pod] in args")
    }

    clusterName := args[0]
    namespaceName := args[1]
    podName := args[2]

    cInfo, err := ep.clusterInfo.GetInfo(clusterName)
    if err != nil {
        return err
    }

    clientset, err := ep.GetClientset(cInfo.Name)
    if err != nil {
        return err
    }

    pod, err := clientset.CoreV1().Pods(namespaceName).Get(context.TODO(), podName, metav1.GetOptions{})
    if err != nil {
        return err
    }

    m := make(map[string]interface{})

    tCS := make([]string, len(pod.Status.ContainerStatuses))
    for idx, container := range pod.Status.ContainerStatuses {
        tCS[idx] = container.Name
    }
    m["containers"] = strings.Join(tCS, ", ")

    m["start_time"] = pod.Status.StartTime.Unix()
    m["node"] = pod.Spec.NodeName
    m["phase"] = pod.Status.Phase
    m["message"] = pod.Status.Message
    m["reason"] = pod.Status.Reason
    m["ip"] = pod.Status.PodIP

    tIF := make(map[string]interface{})
    for _, condition := range pod.Status.Conditions {
        if condition.Status != corev1.ConditionTrue {
            tIF[string(condition.Type)] = condition.Message
        }
    }
    _b, err := json.Marshal(tIF)
    if err != nil {
        return err
    }
    m["conditions"] = string(_b)

    tSS := make([]string, 0)
    _i := int32(0)
    for _, cStatus := range pod.Status.ContainerStatuses {
        if !cStatus.Ready {
            tSS = append(tSS, cStatus.Name)
        }
        _i += cStatus.RestartCount
    }
    m["container_not_ready"] = strings.Join(tSS, ", ")
    m["container_restart_count"] = _i

    b, err := json.Marshal(m)
    if err != nil {
        return err
    }

    fmt.Print(string(b))

    return nil
}

func (ep *Earpiece) node(args ...string) (error) {
    if len(args) < 2 {
        return fmt.Errorf("must intpu [cluster, node] in args")
    }

    clusterName := args[0]
    nodeName := args[1]

    cInfo, err := ep.clusterInfo.GetInfo(clusterName)
    if err != nil {
        return err
    }

    clientset, err := ep.GetClientset(cInfo.Name)
    if err != nil {
        return err
    }

    node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
    if err != nil {
        return err
    }

    m := make(map[string]interface{})

    m["unschedulable"] = node.Spec.Unschedulable
    m["cpu"] = node.Status.Capacity.Cpu()

    val, ok := node.Status.Allocatable.Memory().AsInt64()
    if ok {
        m["memory"] = val
    } else {
        m["memory"] = 0
    }

    m["pods"] = node.Status.Capacity.Pods()
    m["ready"] = corev1.ConditionUnknown

    for _, c := range node.Status.Conditions {
        if c.Type == corev1.NodeReady {
            m["ready"] = c.Status
            break
        }
    }

    b, err := json.Marshal(m)
    if err != nil {
        return err
    }

    fmt.Print(string(b))

    return nil
}