package btesting

import (
	goctx "context"
	"fmt"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"k8s.io/apimachinery/pkg/types"
)

var (
	retryInterval        = time.Second * 5
	timeout              = time.Second * 60
	cleanupRetryInterval = time.Second * 1
	cleanupTimeout       = time.Second * 5
)

func TestSecret(t *testing.T) {
	// s := &corev1.Secret{}
	// err := framework.AddToFrameworkScheme(apis.AddToScheme, s)
	// if err != nil {
	// 	t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	// }
	// run subtests
	t.Run("backuper-group", func(t *testing.T) {
		t.Run("SecretCreation", SecretCreation)
		t.Run("SecretDeletion", SecretDeletion)
		t.Run("SecretEdition", SecretEdition)
	})
}

func SecretCreation(t *testing.T) {
	f, ctx, err := setupCtx(t)
	if err != nil {
		t.Fatal("failed init ctx", err)
	}
	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Fatal("could not get namespace: %v", err)
	}
	secret := &corev1.Secret{}

	// use TestCtx's create helper to create the object and add a cleanup function for the new object
	err = f.Client.Create(goctx.TODO(), secret, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		t.Fatal(err)
	}
	// wait for deployment(but there is no deployment)
	// err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "ex-secret", 1, retryInterval, timeout)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	err = f.Client.Get(goctx.TODO(), types.NamespacedName{Name: "ex-secret", Namespace: namespace}, secret)
	if err != nil {
		t.Fatal(err)
	}
	err = f.Client.Get(goctx.TODO(), types.NamespacedName{Name: "ex-secret", Namespace: "customNS"}, secret)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("done")
}

func secretCreatedDuplicationTest(t *testing.T, f *framework.Framework, ctx *framework.TestCtx) error {
	return nil
}

func setupCtx(t *testing.T) (f *framework.Framework, ctx *framework.TestCtx, err error) {
	t.Parallel()
	ctx = framework.NewTestCtx(t)
	defer ctx.Cleanup()
	err = ctx.InitializeClusterResources(&framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return f, ctx, fmt.Errorf("failed to initialize cluster resources: %v", err)
	}
	t.Log("Initialized cluster resources")
	namespace, err := ctx.GetNamespace()
	if err != nil {
		return f, ctx, fmt.Errorf("failed to get ns: %v", err)
	}
	t.Logf("ns: %v", namespace)
	// get global framework variables
	f = framework.Global
	// wait for backuper-operator to be ready
	err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, "backuper-operator", 1, retryInterval, timeout)
	if err != nil {
		return f, ctx, fmt.Errorf("failed to deploy operator", err)
	}
	return f, ctx, err
}

func SecretDeletion(t *testing.T) {
	t.Log("done")
}

func SecretEdition(t *testing.T) {
	t.Log("done")
}
