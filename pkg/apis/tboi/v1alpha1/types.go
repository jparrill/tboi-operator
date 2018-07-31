package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ItemList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Item `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Item struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              ItemSpec   `json:"spec"`
	Status            ItemStatus `json:"status,omitempty"`
}

type ItemSpec struct {
	ItemName string `json:"itemName"`
	ItemUrl  string `json:"itemUrl"`
	ItemSize int32  `json:"itemSize"`
}

type ItemStatus struct {
	ApiPods []string `"json:apiPods"`
}
