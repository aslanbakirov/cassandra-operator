package v1
import(
   meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const(
    CRDPlural      string = "people"
	CRDGroup       string = "aslangroup.io"
	CRDVersion     string = "v1"
	FullCRDName    string = CRDPlural + "." + CRDGroup
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Person struct{
   meta_v1.TypeMeta `json:",inline"`
   meta_v1.ObjectMeta `json:"metadata"`
   Spec PersonSpec `json:"spec"`
   Status PersonStatus `json:"status, omitempty"`
}

type PersonSpec struct{
    Age string `json:"age"`
    Gender string `json:"gender"`
}

type PersonStatus struct {
    State string `json:"state,omitempty"`
    Message string `json:"message omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PersonList struct{
    meta_v1.TypeMeta `json:",inline"`
    meta_v1.ListMeta `json:"metadata"`
    Items []Person `json:"items"`
}
