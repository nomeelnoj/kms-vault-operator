package e2e

import (
	"context"
	"os"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"

	vaultapi "github.com/hashicorp/vault/api"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	apis "github.com/patoarvizu/kms-vault-operator/pkg/apis"
	operator "github.com/patoarvizu/kms-vault-operator/pkg/apis/k8s/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKMSVaultSecretV1(t *testing.T) {
	ctx := framework.NewTestCtx(t)
	defer ctx.Cleanup()
	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Fatalf("could not get namespace: %v", err)
	}

	awsSecret := &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "aws-secrets",
			Namespace: namespace,
		},
		StringData: map[string]string{
			"AWS_ACCESS_KEY_ID":     os.Getenv("AWS_ACCESS_KEY_ID"),
			"AWS_SECRET_ACCESS_KEY": os.Getenv("AWS_SECRET_ACCESS_KEY"),
		},
	}

	framework.Global.Client.Create(context.TODO(), awsSecret, &framework.CleanupOptions{TestContext: ctx, Timeout: time.Second * 5, RetryInterval: time.Second * 1})

	ctx.InitializeClusterResources(&framework.CleanupOptions{TestContext: ctx, Timeout: time.Second * 5, RetryInterval: time.Second * 1})

	err = e2eutil.WaitForOperatorDeployment(t, framework.Global.KubeClient, namespace, "kms-vault-operator", 1, time.Second*5, time.Second*30)
	if err != nil {
		t.Fatal(err)
	}

	kmsVaultSecretList := &operator.KMSVaultSecretList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "KMSVaultSecret",
			APIVersion: "k8s.patoarvizu.dev/v1alpha1",
		},
	}
	err = framework.AddToFrameworkScheme(apis.AddToScheme, kmsVaultSecretList)
	if err != nil {
		t.Fatalf("Failed to add to scheme: %s", err)
	}

	secret := &operator.KMSVaultSecret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "KMSVaultSecret",
			APIVersion: "k8s.patoarvizu.dev/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: namespace,
		},
		Spec: operator.KMSVaultSecretSpec{
			Path:            "secret/test-secret",
			VaultAuthMethod: "k8s",
			KVSettings: operator.KVSettings{
				EngineVersion: "v1",
			},
			Secrets: []operator.Secret{
				operator.Secret{
					Key:             "Hello",
					EncryptedSecret: "AQICAHgKbLYZWOFlPGwA/1foMoxcBOxv7LddQQW9biqG70YNkwF+dKr15L/4Pl/d26uDd7KqAAAAYzBhBgkqhkiG9w0BBwagVDBSAgEAME0GCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMz0gfMT1P5MBTd/fGAgEQgCANG/RycP+0ZXj2qZORafZO4fGdU7KGFINsrs1JDnx1mg==",
				},
			},
		},
	}

	framework.Global.Client.Create(context.TODO(), secret, &framework.CleanupOptions{TestContext: ctx, Timeout: time.Second * 5, RetryInterval: time.Second * 1})
	if err != nil {
		t.Fatalf("failed to create secret: %v", err)
	}

	time.Sleep(time.Second * 5)

	vaultSecret, err := framework.Global.KubeClient.CoreV1().Secrets(namespace).Get("vault-unseal-keys", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("Failed to get Vault root token: %v", err)
	}
	vaultClient, err := vaultapi.NewClient(vaultapi.DefaultConfig())
	if err != nil {
		t.Fatalf("Failed to get Vault client: %v", err)
	}
	vaultClient.SetToken(string(vaultSecret.Data["vault-root"]))
	vaultClient.Auth()

	r, err := vaultClient.Logical().Read("secret/test-secret")
	if err != nil {
		t.Fatalf("Could not read secret from Vault: %v", err)
	}

	if r == nil {
		t.Errorf("Vault result is empty")
	}

	if val, ok := r.Data["Hello"]; ok {
		if val != "World" {
			t.Errorf("Encrypted string wasn't decrypted correctly")
		}
	} else {
		t.Errorf("Secret wasn't successfully put in Vault")
	}
}

func TestKMSVaultSecretFinalizers(t *testing.T) {
	ctx := framework.NewTestCtx(t)
	defer ctx.Cleanup()
	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Fatalf("could not get namespace: %v", err)
	}

	awsSecret := &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "aws-secrets",
			Namespace: namespace,
		},
		StringData: map[string]string{
			"AWS_ACCESS_KEY_ID":     os.Getenv("AWS_ACCESS_KEY_ID"),
			"AWS_SECRET_ACCESS_KEY": os.Getenv("AWS_SECRET_ACCESS_KEY"),
		},
	}

	framework.Global.Client.Create(context.TODO(), awsSecret, &framework.CleanupOptions{TestContext: ctx, Timeout: time.Second * 5, RetryInterval: time.Second * 1})

	ctx.InitializeClusterResources(&framework.CleanupOptions{TestContext: ctx, Timeout: time.Second * 5, RetryInterval: time.Second * 1})

	err = e2eutil.WaitForOperatorDeployment(t, framework.Global.KubeClient, namespace, "kms-vault-operator", 1, time.Second*5, time.Second*30)
	if err != nil {
		t.Fatal(err)
	}

	kmsVaultSecretList := &operator.KMSVaultSecretList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "KMSVaultSecret",
			APIVersion: "k8s.patoarvizu.dev/v1alpha1",
		},
	}
	err = framework.AddToFrameworkScheme(apis.AddToScheme, kmsVaultSecretList)
	if err != nil {
		t.Fatalf("Failed to add to scheme: %s", err)
	}

	secret := &operator.KMSVaultSecret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "KMSVaultSecret",
			APIVersion: "k8s.patoarvizu.dev/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: namespace,
			Finalizers: []string{
				"delete.k8s.patoarvizu.dev",
			},
		},
		Spec: operator.KMSVaultSecretSpec{
			Path:            "secret/test-secret",
			VaultAuthMethod: "k8s",
			KVSettings: operator.KVSettings{
				EngineVersion: "v1",
			},
			Secrets: []operator.Secret{
				operator.Secret{
					Key:             "Hello",
					EncryptedSecret: "AQICAHgKbLYZWOFlPGwA/1foMoxcBOxv7LddQQW9biqG70YNkwF+dKr15L/4Pl/d26uDd7KqAAAAYzBhBgkqhkiG9w0BBwagVDBSAgEAME0GCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMz0gfMT1P5MBTd/fGAgEQgCANG/RycP+0ZXj2qZORafZO4fGdU7KGFINsrs1JDnx1mg==",
				},
			},
		},
	}

	framework.Global.Client.Create(context.TODO(), secret, nil)
	if err != nil {
		t.Fatalf("failed to create secret: %v", err)
	}

	time.Sleep(time.Second * 5)

	vaultSecret, err := framework.Global.KubeClient.CoreV1().Secrets(namespace).Get("vault-unseal-keys", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("Failed to get Vault root token: %v", err)
	}
	vaultClient, err := vaultapi.NewClient(vaultapi.DefaultConfig())
	if err != nil {
		t.Fatalf("Failed to get Vault client: %v", err)
	}
	vaultClient.SetToken(string(vaultSecret.Data["vault-root"]))
	vaultClient.Auth()

	r, err := vaultClient.Logical().Read("secret/test-secret")
	if err != nil {
		t.Fatalf("Could not read secret from Vault: %v", err)
	}

	if r == nil {
		t.Errorf("Vault result is empty")
	}

	if val, ok := r.Data["Hello"]; ok {
		if val != "World" {
			t.Errorf("Encrypted string wasn't decrypted correctly")
		}
	} else {
		t.Errorf("Secret wasn't successfully put in Vault")
	}

	framework.Global.Client.Delete(context.TODO(), secret)
}
