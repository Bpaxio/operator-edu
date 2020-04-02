package backuper

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileBackuper) deleteSecret(request reconcile.Request) error {
	logger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	logger.Info("Delete secret duplicate")
	found := &corev1.Secret{}
	if err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      request.Name,
		Namespace: request.Namespace,
	}, found); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Not Found secret")
			return nil
		}
		logger.Info("Failed to fetch")
		return err
	}
	if err := r.client.Delete(context.TODO(), found); err != nil {
		logger.Info("Failed delete")
		return err
	}
	return nil
}

func (r *ReconcileBackuper) ensureSecret(request reconcile.Request,
	instance *corev1.Secret,
) (reconcile.Result, error) {
	// Create the secret
	copy := instance.DeepCopy()

	if err := r.createInAnotherNS(copy); err != nil {
		return reconcile.Result{}, err
	}

	duplication := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:        copy.GetName() + ("-copy"),
			Namespace:   copy.GetNamespace(),
			Annotations: copy.GetAnnotations(),
			Labels:      copy.GetLabels(),
		},
		TypeMeta:   copy.TypeMeta,
		Type:       copy.Type,
		Data:       copy.Data,
		StringData: copy.StringData,
	}
	// duplication = instance.DeepCopy()
	// duplication.Namespace = TARGET_NAMESPACE

	found := &corev1.Secret{}
	if err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      duplication.Name,
		Namespace: duplication.Namespace,
	}, found); err == nil {
		log.Info("Skip reconcile: Secret already exists", "Secret.Namespace", duplication.Namespace, "Secret.Name", duplication.Name)
	}

	log.Info("Creating a new secret", "Secret.Namespace", duplication.Namespace, "Secret.Name", duplication.Name)
	if err := r.client.Create(context.TODO(), duplication); err != nil {
		// Creation failed
		log.Error(err, "Failed to create new Secret", "Secret.Namespace", duplication.Namespace, "Secret.Name", duplication.Name)
		return reconcile.Result{}, err
	}
	// Set Secret instance as the owner of duplication
	if err := controllerutil.SetControllerReference(instance, duplication, r.scheme); err != nil {
		log.Error(err, "Failed to set owner of it", "Secret.Namespace", duplication.Namespace, "Secret.Name", duplication.Name)
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileBackuper) checkSecretInCustomNamespace(request reconcile.Request) (*corev1.Secret, error) {
	logger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name, "Target.Namespace", TARGET_NAMESPACE)
	logger.Info("Search secret duplication")
	found := &corev1.Secret{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      request.Namespace + "-ns-" + request.Name,
		Namespace: TARGET_NAMESPACE,
	}, found)

	return found, err
}

func (r *ReconcileBackuper) createInAnotherNS(copy *corev1.Secret) error {

	duplication := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:        copy.Namespace + "-ns-" + copy.GetName(),
			Namespace:   TARGET_NAMESPACE,
			Annotations: copy.GetAnnotations(),
			Labels:      copy.GetLabels(),
		},
		TypeMeta:   copy.TypeMeta,
		Type:       copy.Type,
		Data:       copy.Data,
		StringData: copy.StringData,
	}
	// duplication = instance.DeepCopy()
	// duplication.Namespace = TARGET_NAMESPACE
	log.Info("Creating a new secret", "Secret.Namespace", duplication.Namespace, "Secret.Name", duplication.Name)
	if err := r.client.Create(context.TODO(), duplication); err != nil {
		// Creation failed
		log.Error(err, "Failed to create new Secret", "Secret.Namespace", duplication.Namespace, "Secret.Name", duplication.Name)
		return err
	}
	// Creation was successful
	return nil
}
