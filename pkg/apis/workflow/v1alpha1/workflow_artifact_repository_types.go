package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// WorkflowArtifactRepository is an artifact repository that workflows can refer to by name.
// One named "default" is used by workflows that do not set artifactRepositoryRef.
// +genclient
// +kubebuilder:resource:shortName=wfar
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type WorkflowArtifactRepository struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Spec              ArtifactRepository `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

// WorkflowArtifactRepositoryList is list of WorkflowArtifactRepository resources
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type WorkflowArtifactRepositoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []WorkflowArtifactRepository `json:"items" protobuf:"bytes,2,opt,name=items"`
}
