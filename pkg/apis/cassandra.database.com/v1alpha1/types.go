package v1alpha1
import(
   meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
    Size string `json:"size"`
    Version string `json:"version"`
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
