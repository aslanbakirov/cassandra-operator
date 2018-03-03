package main

import (
	//"flag"
	"context"
	"os"

	crd_controller "github.com/aslanbekirov/cassandra-operator/pkg/controller/cluster"
	"github.com/sirupsen/logrus"
)

func main() {
	//kubeconf := flag.String("kubeconf", "kube.conf", "Path to a kube config. Only required if out-of-cluster.")
	//flag.Parse()
	
	namespace := os.Getenv("MY_POD_NAMESPACE")
	if len(namespace) == 0 {
		logrus.Fatalf("must set env MY_POD_NAMESPACE", )
	}
	logrus.Infof("Starting Cassandra operator in namespace %s", namespace)
	c := crd_controller.New(true, namespace)
	c.Start(context.Background())
}

