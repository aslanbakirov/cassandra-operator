package utils

import(
	"fmt"
	"time"
	
	
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes"
)

func  WaitForStatefulSetReady(clientset kubernetes.Interface, ns string, name string, size int32) error {
	return Retry(10*time.Second, 10, func() (bool, error) {
		
		statefulSet,err := clientset.AppsV1beta1().StatefulSets(ns).Get(name, meta_v1.GetOptions{})
		if err != nil {
			return false,err
		}
		if statefulSet.Status.ReadyReplicas < size {
			return false,nil
		}
		return true, nil
	})
}

func WaitCRDReady(clientset apiextensionsclient.Interface, crdName string) error {
	err := Retry(5*time.Second, 20, func() (bool, error) {
		crd, err := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Get(crdName, meta_v1.GetOptions{})
		if err != nil {
			return false, err
		}
		for _, cond := range crd.Status.Conditions {
			switch cond.Type {
			case apiextensionsv1beta1.Established:
				if cond.Status == apiextensionsv1beta1.ConditionTrue {
					return true, nil
				}
			case apiextensionsv1beta1.NamesAccepted:
				if cond.Status == apiextensionsv1beta1.ConditionFalse {
					return false, fmt.Errorf("Name conflict: %v", cond.Reason)
				}
			}
		}
		return false, nil
	})
	if err != nil {
		return fmt.Errorf("wait CRD created failed: %v", err)
	}
	return nil
}

func NewKubeClient(kubeconfig string) (*rest.Config, error){
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}
