package gateway

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/os-observability/tempo-operator/apis/tempo/v1alpha1"
	"github.com/os-observability/tempo-operator/internal/manifests/manifestutils"
)

const (
	componentName = "gateway"
)

// BuildAll creates objects for Tempo deployment.
func BuildAll(tempo v1alpha1.Microservices) ([]client.Object, error) {
	if tempo.Spec.Tenants == nil || tempo.Spec.Tenants.Mode != v1alpha1.OpenShift {
		return nil, nil
	}

	var objs []client.Object
	objs = append(objs, deployment(tempo))
	objs = append(objs, config(tempo))

	return objs, nil
}

func deployment(tempo v1alpha1.Microservices) client.Object {
	labels := manifestutils.ComponentLabels(componentName, tempo.Name)
	cfg := &v1alpha1.TempoComponentSpec{}
	if userCfg := tempo.Spec.Components.Distributor; userCfg != nil {
		cfg = userCfg
	}
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      manifestutils.Name(componentName, tempo.Name),
			Namespace: tempo.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					NodeSelector: cfg.NodeSelector,
					Tolerations:  cfg.Tolerations,
					Containers: []corev1.Container{
						{
							Name:  "gateway",
							Image: "quay.io/observatorium/api:main-2022-11-22-v0.1.2-287-g0ad5c60",
							Args: []string{
								"--web.listen=0.0.0.0:8080",
								"--web.internal.listen=0.0.0.0:8081",
								fmt.Sprintf("--traces.write.endpoint=%s:4317", manifestutils.Name("distributor", tempo.Name)),
								fmt.Sprintf("--traces.read.endpoint=http://%s:16686", manifestutils.Name("query-frontend", tempo.Name)),
								"--grpc.listen=0.0.0.0:8090",
								"--rbac.config=/etc/observatorium/rbac.yaml",
								"--tenants.config=/etc/observatorium/tenants.yaml",
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "public",
									ContainerPort: 8080,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          "internal",
									ContainerPort: 8081,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          "grpc-public",
									ContainerPort: 8090,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									Exec: &corev1.ExecAction{
										Command: nil,
									},
									HTTPGet: &corev1.HTTPGetAction{
										Path:   "/ready",
										Port:   intstr.FromInt(8081),
										Scheme: "HTTP",
									},
								},
								PeriodSeconds:    5,
								FailureThreshold: 12,
							},
							//VolumeMounts: []corev1.VolumeMount{
							//	{
							//		Name:      configVolumeName,
							//		MountPath: "/conf",
							//		ReadOnly:  true,
							//	},
							//},
						},
					},
				},
			},
		},
	}
}

func config(tempo v1alpha1.Microservices) client.Object {
	return nil
}
