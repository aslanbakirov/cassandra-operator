package controller

import(
	//"fmt"
	
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)


func (c *Cluster) CreateService(s *v1.Service) error{
	
	ser := c.buildService("cassandra");
	
	client := c.kubeClientset.CoreV1().Services(c.namespace)
    
	service, err := client.Get(s.Name, meta_v1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	if errors.IsNotFound(err) {
		_, err = client.Create(ser)
		if err != nil {
			return err
		}
		
	} else {
		service.ResourceVersion = ser.ResourceVersion
		_, err := client.Update(ser)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
	}

	return nil
}

func (c *Cluster) buildService(name string) *v1.Service {

	service:= &v1.Service{
        ObjectMeta: meta_v1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"app": name,
				"role": "cassandraCluster",
			},
			
		},
		Spec: v1.ServiceSpec{
			ClusterIP: "None",
			Ports: []v1.ServicePort{
				{
					Port: 9042,
				},
			},
			Selector: map[string]string {
					"app": name,
				},
			},
		}
	return service
}


