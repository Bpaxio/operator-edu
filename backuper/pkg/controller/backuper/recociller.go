package backuper

import (
	"context"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ReconcileBackuper reconciles a Backuper object
type ReconcileBackuper struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

var TARGET_NAMESPACE = "dev-duplication"

// Reconcile reads that state of the cluster for a Backuper object and makes changes based on the state read
// and what is in the Backuper.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileBackuper) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Backuper")
	if request.Namespace == TARGET_NAMESPACE ||
		request.Namespace == "default" ||
		request.Namespace == "kube-node-lease" ||
		request.Namespace == "kube-public" ||
		request.Namespace == "kube-system" ||
		strings.Contains(request.Name, "-copy") ||
		strings.Contains(request.Name, TARGET_NAMESPACE+"-ns-") {
		reqLogger.Info("Skip reconcile: Secret is special")
		// secrets of target and default namespaces should be skipped
		return reconcile.Result{}, nil
	}

	// Fetch the Secret instance
	instance := &corev1.Secret{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("Not Found secret")
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Check if this Secret duplication already exists
	found, fail := r.checkSecretInCustomNamespace(request)
	if fail != nil {
		if errors.IsNotFound(fail) {
			reqLogger.Info("Not Found secret", "Secret.Namespace", found.Namespace, "Secret.Name", found.Name)
			// Define a new Secret object
			return r.ensureSecret(request, instance)
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, fail
	}

	// secret already exists - don't requeue
	reqLogger.Info("Skip reconcile: Secret already exists", "Secret.Namespace", found.Namespace, "Secret.Name", found.Name)
	return reconcile.Result{}, nil
}
