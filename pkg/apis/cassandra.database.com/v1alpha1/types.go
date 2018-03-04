package v1alpha1

import(
   meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
   "k8s.io/api/core/v1"
)

const(
    CRDPlural      string = "cassandraclusters"
	CRDGroup       string = "cassandra.database.com"
	CRDVersion     string = "v1alpha1"
	FullCRDName    string = CRDPlural + "." + CRDGroup
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CassandraCluster struct{
   meta_v1.TypeMeta `json:",inline"`
   meta_v1.ObjectMeta `json:"metadata"`
   Spec CassandraClusterSpec `json:"spec"`
   Status CassandraClusterStatus `json:"status, omitempty"`
}

type CassandraClusterSpec struct{
    Size int32 `json:"size"`
    Version string `json:"version"`
    PodSpec PodSpec `json:"pod, omitempty"`
    TLS *TLSPolicy `json:"TLS,omitempty"`
}

type CassandraClusterStatus struct {
    State string `json:"state,omitempty"`
    Message string `json:"message omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CassandraClusterList struct{
    meta_v1.TypeMeta `json:",inline"`
    meta_v1.ListMeta `json:"metadata"`
    Items []CassandraCluster `json:"items"`
}

type PodSpec struct{
    Image string `json:"image,omitempty"`
    Labels map[string]string `json:"labels,omitempty"`
    NodeSelector map[string]string `json:"nodeSelector,omitempty"`
    AntiAffinity bool `json:"antiAffinity,omitempty"`
    ServiceAccountName string `json:"serviceAccountName, omitempty"`
    Resources v1.ResourceRequirements `json:"resources,omitempty"`
    Tolerations []v1.Toleration `json:"tolerations,omitempty"`
    PV PVSpec `json:"pv,omitempty"`
    Env []v1.EnvVar `json:"env,omitempty"`
}

type PVSpec struct{
    VolumeSize string `json:"volumeSize"`
    StorageClass string `json:"storageClass"`
    Name string `json:"name, omitempty"`
    MountPath string `json:"mountPath, omitempty"`
}