package main

import (
    "flag"

    "github.com/AcidGo/aks-earpiece/common"
)

const (
    KEY_CLUSTER_IP          "ip"
    KEY_CLUSTER_KUBECFG     "kubecfg"
)

var (
    earpiece    *common.Earpiece
    cInfo       *common.ClusterInfo
    cfgPath     string
)

func init() {
    flag.StringVar(&cfgPath, "f", "aks-earpiece.ini", "cluster info configure file path")
    flag.Parse()

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

    earpiece, err = common.NewEarpice(cInfo)
    if err != nil {
        log.Fatal(err)
    }

    err = earpiece.Call(...flag.Args())
    if err != nil {
        log.Fatal(err)
    }
}