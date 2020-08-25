package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KMSVaultSecretSpec defines the desired state of KMSVaultSecret
// +k8s:openapi-gen=true
type KMSVaultSecretSpec struct {
	Path string `json:"path"`

	// +listType=map
	// +listMapKey=key
	Secrets       []Secret          `json:"secrets"`
	SecretContext map[string]string `json:"secretContext,omitempty"`

	// +listType=set
	IncludeSecrets []string `json:"includeSecrets,omitempty"`

	KVSettings KVSettings `json:"kvSettings"`
}

type KVSettings struct {
	// +kubebuilder:validation:Enum={"v1","v2"}
	EngineVersion string `json:"engineVersion"`
	// +kubebuilder:validation:Minimum=0
	CASIndex int `json:"casIndex,omitempty"`
}

type Secret struct {
	Key             string            `json:"key"`
	EncryptedSecret string            `json:"encryptedSecret,omitempty"`
	SecretContext   map[string]string `json:"secretContext,omitempty"`
	EmptySecret     bool              `json:"emptySecret,omitempty"`
}

// KMSVaultSecretStatus defines the observed state of KMSVaultSecret
// +k8s:openapi-gen=true
type KMSVaultSecretStatus struct {
	Created bool `json:"created,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KMSVaultSecret is the Schema for the kmsvaultsecrets API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=kmsvaultsecrets,scope=Namespaced,shortName=kmsvs
type KMSVaultSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KMSVaultSecretSpec   `json:"spec,omitempty"`
	Status KMSVaultSecretStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KMSVaultSecretList contains a list of KMSVaultSecret
type KMSVaultSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KMSVaultSecret `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KMSVaultSecret{}, &KMSVaultSecretList{})
}
