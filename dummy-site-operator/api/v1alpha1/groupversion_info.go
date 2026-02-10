package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	GroupeVersion = schema.GroupVersion{
		Group:   "dummy-site-operator.stac.dev",
		Version: "v1alpha1",
	}

	SchemeBuilder = &scheme.Builder{GroupVersion: GroupeVersion}
	AddToScheme   = SchemeBuilder.AddToScheme
)
