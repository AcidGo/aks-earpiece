package main

import (
    "flag"
    "log"

    "github.com/AcidGo/aks-earpiece/common"
    "gopkg.in/ini.v1"
)

const (
    KEY_CLUSTER_IP          = "ip"
    KEY_CLUSTER_KUBECFG     = "kubecfg"
)

var (
    earpiece    *common.Earpiece
    cInfo       *common.ClusterInfo
    ops         *common.Options
    cfgPath     string
)

func init() {
    flag.StringVar(&cfgPath, "f", "aks-earpiece.ini", "cluster info configure file path")
    flag.Parse()

    if len(flag.Args()) < 1 {
        log.Fatal("the args is emtpy")
    }

    ops = &common.Options{
        Method:     flag.Arg(0),
        Args:       flag.Args(),
    }

    cInfo = &common.ClusterInfo{}

    cfg, err := ini.Load(cfgPath)
    if err != nil {
        log.Fatal(err)
    }

    for _, sec := range cfg.Sections() {
        clusterName     := sec.Name()
        clusterIP       := sec.Key(KEY_CLUSTER_IP).String()
        clusterKubecfg  := sec.Key(KEY_CLUSTER_KUBECFG).String()
        c := &common.Cluster{
            Name:       clusterName,
            IP:         clusterIP,
            Kubecfg:    clusterKubecfg,
        }

        err = cInfo.AddInfo(c)
        if err != nil {
            log.Fatal(err)
        }
    }
}

func main() {
    var err error

    earpiece, err = common.NewEarpice(cInfo)
    if err != nil {
        log.Fatal(err)
    }

    err = earpiece.Call(ops)
    if err != nil {
        log.Fatal(err)
    }
}