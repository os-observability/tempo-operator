package manifestutils

import (
	"k8s.io/apimachinery/pkg/labels"
)

// ComponentLabels is a list of all commonLabels including the app.kubernetes.io/component:<component> label.
func ComponentLabels(component, instanceName string) labels.Set {
	return labels.Merge(CommonLabels(instanceName), map[string]string{
		"app.kubernetes.io/component": component,
	})
}

func CommonLabels(instanceName string) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       "tempo",
		"app.kubernetes.io/instance":   instanceName,
		"app.kubernetes.io/managed-by": "tempo-controller",
		"app.kubernetes.io/created-by": "tempo-controller",
	}
}
