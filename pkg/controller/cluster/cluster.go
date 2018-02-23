package controller

import (
	//"fmt"
	
	//meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	co_v1aplha1 "github.com/aslanbekirov/cassandra-operator/pkg/apis/cassandra.database.com/v1alpha1"
	utils "github.com/aslanbekirov/cassandra-operator/pkg/utils"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	core_v1 "k8s.io/api/core/v1"
)


// func (c *Cluster) createCluster(spec *co_v1aplha1.CassandraClusterSpec){
	
	
// }

func (c *Cluster) DeleteStatefulSet(ssName string) error{
	err := c.kubeClientset.AppsV1().StatefulSets(c.namespace).Delete(ssName, &meta_v1.DeleteOptions{
		PropagationPolicy: func() *meta_v1.DeletionPropagation {
			foreground := meta_v1.DeletePropagationForeground
			return &foreground
		}(),
	})
	if errors.IsNotFound(err) { 
		err = nil
	}
	return err
}

func (c *Cluster) CreateOrUpdateStatefulSet(ss *v1.StatefulSet) error{

	client := c.kubeClientset.AppsV1().StatefulSets(c.namespace)
    
	statefulSet, err := client.Get(ss.Name, meta_v1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	if errors.IsNotFound(err) {
		_, err = client.Create(ss)
		if err != nil {
			return err
		}
		err = utils.WaitForStatefulSetReady(c.kubeClientset, c.namespace, ss.Name, *ss.Spec.Replicas)
		if err != nil {
			return err
		}
	} else {
		ss.ResourceVersion = statefulSet.ResourceVersion
		_, err := client.Update(ss)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
	}

	return nil
}



func (c *Cluster) BuildStatefulSet(cc *co_v1aplha1.CassandraCluster) *v1.StatefulSet{

	limitCPU, _ := resource.ParseQuantity(cc.Spec.PodSpec.Resources.Limits.Cpu().String())
	limitMemory, _ := resource.ParseQuantity(cc.Spec.PodSpec.Resources.Limits.Memory().String())
	requestCPU, _ := resource.ParseQuantity(cc.Spec.PodSpec.Resources.Requests.Cpu().String())
	requestMemory, _ := resource.ParseQuantity(cc.Spec.PodSpec.Resources.Requests.Memory().String())
	requestDataStorage,_ := resource.ParseQuantity(cc.Spec.PodSpec.PV.VolumeSize)

	var antiAffinity *core_v1.Affinity
	if (cc.Spec.PodSpec.AntiAffinity == true){
		antiAffinity = &core_v1.Affinity{
			PodAntiAffinity: &core_v1.PodAntiAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: []core_v1.PodAffinityTerm{
					{
						LabelSelector: &meta_v1.LabelSelector{
							MatchExpressions: []meta_v1.LabelSelectorRequirement{
								{
									Key:      "cassandraCluster",
									Operator: meta_v1.LabelSelectorOpIn,
									Values:   []string{cc.ObjectMeta.Name},
								},
							},
						},
						TopologyKey: "kubernetes.io/hostname",
					},
				},
			},
		}
	}else{
		antiAffinity = nil
	}

	statefulSet := &v1.StatefulSet{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: cc.Name,
			Labels: map[string]string{
				"cassandraCluster": cc.Name,
				"role": "cassandraCluster",
			},
			Annotations: map[string]string{
				"operatorVersion": co_v1aplha1.SchemeGroupVersion.Version,

			},
		},
		Spec: v1.StatefulSetSpec{
			ServiceName: cc.Name+"-node", 
			Selector: &meta_v1.LabelSelector{
				MatchLabels: map[string]string {
					"cassandraCluster": cc.Name,
				},
			},
			UpdateStrategy: v1.StatefulSetUpdateStrategy{
				Type: v1.RollingUpdateStatefulSetStrategyType,
				RollingUpdate: &v1.RollingUpdateStatefulSetStrategy{
					Partition: func(i int) *int32 { j:=int32(i);return &j}(0),
				},
			},
			Replicas: &cc.Spec.Size,
			//ServiceName: mc.ObjectMeta.Name,
			Template: core_v1.PodTemplateSpec{
				ObjectMeta: meta_v1.ObjectMeta{
					Labels: map[string]string{
						"cassandraCluster": cc.Name,
						"role": "cassandraCluster",
					},
					Annotations: map[string]string{
						"operatorVersion": co_v1aplha1.SchemeGroupVersion.Version,

					},
				},
				Spec: core_v1.PodSpec{

					Affinity: antiAffinity,
					TerminationGracePeriodSeconds: func(i int64) *int64 { return &i}(10),
					/*Volumes: []core_v1.Volume{
						{
							Name:	"secret",
							VolumeSource: core_v1.VolumeSource{
								Secret: &core_v1.SecretVolumeSource{
									SecretName: cc.Name+"-node-token",
									DefaultMode: func(i int32) *int32 {return &i}(256),
								},
							},
						},
					},*/
					Containers: []core_v1.Container{
						{
							Name:            "cassandra",
							Image:           cc.Spec.PodSpec.Image,
							ImagePullPolicy: "Always",
							Env: cc.Spec.PodSpec.Env,
							// Env: []core_v1.EnvVar{
							// 	{
							// 		Name: "NAMESPACE",
							// 		ValueFrom: &core_v1.EnvVarSource{
							// 			FieldRef: &core_v1.ObjectFieldSelector{
							// 				FieldPath: "metadata.namespace",
							// 			},
							// 		},
							// 	},
							// 	{
							// 		Name: "MAX_HEAP_SIZE",
							// 		Value: cc.Spec.PodSpec.Env[0].Value,
							// 	},
							// 	{
							// 		Name: "HEAP_NEW_SIZE",
							// 		Value: cc.Spec.PodSpec.Env[1].Value,
							// 	},
							// 	{
							// 		Name: "CASSANDRA_SEEDS",
							// 		Value: cc.Name+"-0."+cc.Name+"-node."+c.namespace+".svc.cluster.local",
							// 	},
							// 	{
							// 		Name: "CASSANDRA_CLUSTER_NAME",
							// 		Value: cc.Name,
							// 	},
							// 	{
							// 		Name: "CASSANDRA_DC",
							// 		Value: cc.Spec.PodSpec.Env[1].Value,
							// 	},
							// 	{
							// 		Name: "CASSANDRA_RACK",
							// 		Value: cc.Spec.PodSpec.Env[1].Value,
							// 	},
							// 	{
							// 		Name: "CASSANDRA_AUTO_BOOTSTRAP",
							// 		Value: cc.Spec.PodSpec.Env[1].Value,
							// 	},
							// 	{
							// 		Name: "POD_IP",
							// 		ValueFrom: &core_v1.EnvVarSource{
							// 			FieldRef: &core_v1.ObjectFieldSelector{
							// 				FieldPath: "status.podIP",
							// 			},
							// 		},
							// 	},
							//},
							Ports: []core_v1.ContainerPort{
								{
									Name:          "cql",
									ContainerPort: 9042,
								},
								{
									Name:          "jmx",
									ContainerPort: 7199,
								},
								{
									Name:          "tls-intra-node",
									ContainerPort: 7001,
								},
								{
									Name:          "intra-node",
									ContainerPort: 7000,
								},
							},
							SecurityContext: &core_v1.SecurityContext{
								Capabilities: &core_v1.Capabilities{
									Add: []core_v1.Capability{
										"IPC_LOCK",
									},
								},
							},
							Lifecycle: &core_v1.Lifecycle{
								PreStop: &core_v1.Handler{
									Exec: &core_v1.ExecAction{
										Command: []string {
											"/bin/sh",
											"-c",
											"PID=$(pidof java) && kill $PID && while ps -p $PID > /dev/null; do sleep 1; done",
										},
									},
								},
							},
							ReadinessProbe: &core_v1.Probe{
								Handler: core_v1.Handler{
									Exec: &core_v1.ExecAction{
										Command: []string {
											"/bin/bash",
											"-c",
											"/ready-probe.sh",
										},
									},
								},
								InitialDelaySeconds: int32(15),
								TimeoutSeconds: int32(5),
							},
							VolumeMounts: []core_v1.VolumeMount{
								{
									Name:      "data",
									MountPath: "/cassandra_data",
								},
							},
							Resources: core_v1.ResourceRequirements{
								Limits: core_v1.ResourceList{
									"cpu":    limitCPU,
									"memory": limitMemory,
								},
								Requests: core_v1.ResourceList{
									"cpu":    requestCPU,
									"memory": requestMemory,
								},
							},
						},
					},
				},
			},
			VolumeClaimTemplates: []core_v1.PersistentVolumeClaim{
				{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: "data",
						Annotations: map[string]string{
							"volume.beta.kubernetes.io/storage-class": cc.Spec.PodSpec.PV.StorageClass,
						},
						Labels: map[string]string{
							"name":      cc.Name,
						},
					},
					Spec: core_v1.PersistentVolumeClaimSpec{
						AccessModes: []core_v1.PersistentVolumeAccessMode{
							core_v1.ReadWriteOnce,
						},
						Resources: core_v1.ResourceRequirements{
							Requests: core_v1.ResourceList{
								core_v1.ResourceStorage: requestDataStorage,
							},
						},
					},
				},
			},
		},
	}
	return statefulSet
}