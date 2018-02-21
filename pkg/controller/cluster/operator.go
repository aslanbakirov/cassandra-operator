package controller

import (
	"context"
	"fmt"
	"reflect"

	co_v1aplha1 "github.com/aslanbekirov/cassandra-operator/pkg/apis/cassandra.database.com/v1alpha1"
	utils "github.com/aslanbekirov/cassandra-operator/pkg/utils"
	"github.com/sirupsen/logrus"
	v1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

//Cluster type definition
type Cluster struct {
	logger *logrus.Entry

	namespace string
	// k8s workqueue pattern
	indexer  cache.Indexer
	informer cache.SharedIndexInformer
	queue    workqueue.RateLimitingInterface

	kubeconf string

	createCustomResource bool
}

//New Create new Cluster Instance
func New(createCRD bool, kubeconf string) *Cluster {
	return &Cluster{
		logger:    logrus.WithField("pkg", "controller"),
		namespace: "test",
		kubeconf:   kubeconf,
		createCustomResource: createCRD,
	}
}

// Start starts the Cassandra Cluster operator.
func (c *Cluster) Start(ctx context.Context) error {
	if c.createCustomResource {
		if err := c.createCRD(); err != nil {
			return err
		}
	}

	go c.run(ctx)
	<-ctx.Done()
	return ctx.Err()
}

func (c *Cluster) createCRD() error {

	cassandraCluster := &v1beta1.CustomResourceDefinition{
		ObjectMeta: meta_v1.ObjectMeta{Name: co_v1aplha1.FullCRDName},
		Spec: v1beta1.CustomResourceDefinitionSpec{
			Group:   co_v1aplha1.CRDGroup,
			Version: co_v1aplha1.CRDVersion,
			Scope:   v1beta1.NamespaceScoped,
			Names: v1beta1.CustomResourceDefinitionNames{
				Plural: co_v1aplha1.CRDPlural,
				Kind:   reflect.TypeOf(co_v1aplha1.CassandraCluster{}).Name(),
			},
		},
	}

    r, err:=utils.NewKubeClient(c.kubeconf)

	apiextensionsClient, err := apiextensionsclientset.NewForConfig(r)
	_, err = apiextensionsClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(cassandraCluster)

	if err != nil {
		fmt.Println("Erro occured creating cassandra cluster crd, %v", err)
		panic(err)
	}

	err = utils.WaitCRDReady(apiextensionsClient, co_v1aplha1.FullCRDName)
	if err != nil {
		fmt.Println("Crd is created and ready to use, %s", co_v1aplha1.FullCRDName)
		panic(err)
	}
	return err
}
