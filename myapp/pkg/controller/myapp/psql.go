package myapp

import (
	"context"

	bpaxiov1 "bpax.io/ru/cmx/edu/MyOperators/myapp/pkg/apis/bpaxio/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func psqlDeploymentName() string {
	return "psql"
}

func psqlServiceName() string {
	return "psql-service"
}

func psqlAuthName() string {
	return "psql-auth"
}

func (r *ReconcileMyApp) psqlAuthSecret(cr *bpaxiov1.MyApp) *corev1.Secret {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      psqlAuthName(),
			Namespace: cr.Namespace,
		},
		Type: "Opaque",
		StringData: map[string]string{
			"username": "mayapp-user",
			"password": "mayapp-pass",
		},
	}
	controllerutil.SetControllerReference(cr, secret, r.scheme)
	return secret
}

func (r *ReconcileMyApp) psqlDeployment(cr *bpaxiov1.MyApp) *appsv1.Deployment {
	labels := labels(cr, "psql")
	size := int32(1)

	userSecret := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: psqlAuthName()},
			Key:                  "username",
		},
	}

	passwordSecret := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: psqlAuthName()},
			Key:                  "password",
		},
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      psqlDeploymentName(),
			Namespace: cr.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "postgres:9.6.17",
						Name:  cr.Name + "-psql",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 5432,
							Name:          "postgres",
						}},
						Env: []corev1.EnvVar{
							{
								Name:  "PSQL_ROOT_PASSWORD",
								Value: "password",
							},
							{
								Name:  "PG_DATA",
								Value: "/data/postgres",
							},
							{
								Name:  "POSTGRES_DB",
								Value: "example",
							},
							{
								Name:      "POSTGRES_USER",
								ValueFrom: userSecret,
							},
							{
								Name:      "POSTGRES_PASSWORD",
								ValueFrom: passwordSecret,
							},
						},
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(cr, dep, r.scheme)
	return dep
}

func (r *ReconcileMyApp) psqlService(cr *bpaxiov1.MyApp) *corev1.Service {
	labels := labels(cr, "psql")

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      psqlServiceName(),
			Namespace: cr.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Port: 3306,
			}},
			ClusterIP: "None",
		},
	}

	controllerutil.SetControllerReference(cr, s, r.scheme)
	return s
}

// Returns whether or not the PSQL deployment is running
func (r *ReconcileMyApp) isPsqlUp(cr *bpaxiov1.MyApp) bool {
	deployment := &appsv1.Deployment{}

	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      psqlDeploymentName(),
		Namespace: cr.Namespace,
	}, deployment)

	if err != nil {
		log.Error(err, "Deployment psql not found")
		return false
	}

	if deployment.Status.ReadyReplicas == 1 {
		return true
	}

	return false
}
