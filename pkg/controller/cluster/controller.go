package controller

import (
	"context"
	"time"

	external_versions "github.com/aslanbekirov/cassandra-operator/pkg/client/informers/externalversions"
	clientset "github.com/aslanbekirov/cassandra-operator/pkg/client/clientset/versioned"
	utils "github.com/aslanbekirov/cassandra-operator/pkg/utils"
	
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)


func (c *Cluster) run(ctx context.Context) {

	c.queue = workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "cassandra-operator")
	
	r, err:=utils.NewKubeClient(c.kubeconf)

	clientset, err := clientset.NewForConfig(r)
	if err != nil {
		panic(err.Error())
	}
	
	factory := external_versions.NewSharedInformerFactory(clientset, time.Minute*3)
	
	c.informer = factory.Cassandra().V1alpha1().CassandraClusters().Informer()
    c.indexer = factory.Cassandra().V1alpha1().CassandraClusters().Informer().GetIndexer()
	 
	c.informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				c.onAdd(obj)
			},
			DeleteFunc: func(obj interface{}) {
				c.onDelete(obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				c.onUpdate(oldObj,newObj)
			},
		},
	)
	stop := make(chan struct{})

	go c.informer.Run(stop)

	if !cache.WaitForCacheSync(ctx.Done(), c.informer.HasSynced) {
		return
	}

	const numWorkers = 1
	for i := 0; i < numWorkers; i++ {
		go wait.Until(c.runWorker, time.Second, ctx.Done())
	}

	<-ctx.Done()
	c.logger.Info("stopping cassandra cluster controller")

}

func (c *Cluster) onAdd(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		panic(err)
	}
	c.queue.Add(key)
}

func (c *Cluster) onUpdate(oldObj, newObj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(newObj)
	if err != nil {
		panic(err)
	}
	c.queue.Add(key)
}

func (c *Cluster) onDelete(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		panic(err)
	}
	c.queue.Add(key)
}
