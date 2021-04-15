package main

import (
    "encoding/json"
    "log"
    "os"

    "k8s.io/client-go/tools/clientcmd"
)

var (
    ops         *options
    clientset   *kubernetes.Clientset
)

type options struct {
    src     string
    args    []string
    f       func (...string) (error)
}

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

func init() {
    if len(os.Args) < 2 {
        log.Fatal("input args must not be empty")
    }

    ops = &options{
        src:    os.Arg[1],
    }

    if len(os.Args) == 3 {
        ops.args = make([]string, 0)
    } else {
        ops.args = os.Args[2:]
    }

    switch ops.src {
    case "discovery_cl":
    case "discovery_ns": 
    case "discovery_nd": 
    case "discovery_pod":
    case "cl_cs":
    case "pod":
    case "no":
    default:
        log.Fatalf("the option src %s is invalid", ops.src)
    }
}

func main() {
    kubeconfig := filepath.Join(home, ".kube", "config")
    config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
    if err != nil {
        log.Fatal(err)
    }

    clientset, err = kubernetes.NewForConfig(config)
    if err != nil {
        log.Fatal(err)
    }

    err = ops.f(ops.args...)
    if err != nil {
        log.Fatal(err)
    }
}

func DiscoveryCluster(args ...string) (error) {
    zd := newZbxDiscovery()

    
}