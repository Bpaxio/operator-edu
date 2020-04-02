package backuper

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileBackuper) checkSecretInCustomNamespace(request reconcile.Request, instance *corev1.Secret) (*corev1.Secret, error) {
	logger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	logger.Info("Search secret duplication")
	found := &corev1.Secret{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      request.Name,
		Namespace: TARGET_NAMESPACE,
	}, found)

	return found, err
}

func (r *ReconcileBackuper) ensureSecret(request reconcile.Request,
	instance *corev1.Secret,
	duplication *corev1.Secret,
) (reconcile.Result, error) {
	// Create the secret
	log.Info("Creating a new secret", "Secret.Namespace", duplication.Namespace, "Secret.Name", duplication.Name)
	if err := r.client.Create(context.TODO(), duplication); err != nil {
		// Creation failed
		log.Error(err, "Failed to create new Secret", "Secret.Namespace", duplication.Namespace, "Secret.Name", duplication.Name)
		return reconcile.Result{}, err
	}
	// Creation was successful

	// Set Secret instance as the owner of duplication
	if err := controllerutil.SetControllerReference(instance, duplication, r.scheme); err != nil {
		log.Error(err, "Failed to set owner of it", "Secret.Namespace", duplication.Namespace, "Secret.Name", duplication.Name)
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}
