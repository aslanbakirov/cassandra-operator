package main

import (
	"flag"
	"fmt"
	"context"

	crd_controller "github.com/aslanbekirov/cassandra-operator/pkg/controller/cluster"
)

func main() {
	kubeconf := flag.String("kubeconf", "kube.conf", "Path to a kube config. Only required if out-of-cluster.")
	flag.Parse()
	fmt.Println(kubeconf)

	c := crd_controller.New(true, *kubeconf)
	c.Start(context.TODO())
}

