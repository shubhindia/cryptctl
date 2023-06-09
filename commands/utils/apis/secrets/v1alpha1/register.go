// Package v1beta2 contains API Schema definitions for the secrets v1beta2 API group
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=package,register
// +k8s:conversion-gen=github.com/shubhindia/hcictl/commands/utils/apis/secrets
// +k8s:defaulter-gen=TypeMeta
// +groupName=secrets.shubhindia.xyz
package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: "secrets.shubhindia.xyz", Version: "v1alpha"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}
)
