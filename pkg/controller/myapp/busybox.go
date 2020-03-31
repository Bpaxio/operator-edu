package myapp

import (
	bpaxiov1 "bpax.io/ru/cmx/edu/MyApp/pkg/apis/bpaxio/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// newBusyBoxPod returns a busybox pod with the same name/namespace as the cr
func newBusyBoxPod(cr *bpaxiov1.MyApp) *corev1.Pod {
	labels := labels(cr, "busybox")
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}
