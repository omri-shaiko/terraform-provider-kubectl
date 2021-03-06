package kubernetes

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"k8s.io/apimachinery/pkg/api/errors"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"kubectl": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccCheckkubectlDestroy(s *terraform.State) error {
	return testAccCheckkubectlStatus(s, false)
}

func testAccCheckkubectlExists(s *terraform.State) error {
	return testAccCheckkubectlStatus(s, true)
}

func testAccCheckkubectlStatus(s *terraform.State, shouldExist bool) error {
	provider := testAccProvider.Meta().(*KubeProvider)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kubectl_manifest" {
			continue
		}

		content, err := provider.MainClientset.RESTClient().Get().AbsPath(rs.Primary.ID).DoRaw()
		if (errors.IsNotFound(err) || errors.IsGone(err)) && shouldExist {
			return fmt.Errorf("Failed to find resource, likely a failure to create occured: %+v %v", err, string(content))
		}

	}

	return nil
}

func TestProviderConfigureWithExec(t *testing.T) {
	if os.Getenv("TF_ACC") != "" {
		t.Skip("The environment variable TF_ACC is set, and this test prevents acceptance tests" +
			" from running as it alters environment variables - skipping")
	}

	resetEnv := unsetEnv(t)
	defer resetEnv()

	const conf string = `
host                   = "testhost"
cluster_ca_certificate = "testcert"

exec {
  api_version = "client.authentication.k8s.io/v1alpha1"
  command     = "aws"

  args = [
    "eks",
    "get-token",
    "--cluster-name",
    "testcluster",
    "--region",
    "us-east-1",
  ]
}`
	var confmap map[string]interface{}
	if err := hcl.Decode(&confmap, conf); err != nil {
		t.Fatal(err)
	}

	c, err := config.NewRawConfig(confmap)
	if err != nil {
		t.Fatal(err)
	}
	rc := terraform.NewResourceConfig(c)
	p := Provider()
	err = p.Configure(rc)
	if err != nil {
		t.Fatal(err)
	}
}

func unsetEnv(t *testing.T) func() {
	e := getEnv()

	if err := os.Unsetenv("KUBECONFIG"); err != nil {
		t.Fatalf("Error unsetting env var KUBECONFIG: %s", err)
	}
	if err := os.Unsetenv("KUBE_CONFIG"); err != nil {
		t.Fatalf("Error unsetting env var KUBE_CONFIG: %s", err)
	}
	if err := os.Unsetenv("KUBE_CTX"); err != nil {
		t.Fatalf("Error unsetting env var KUBE_CTX: %s", err)
	}
	if err := os.Unsetenv("KUBE_CTX_AUTH_INFO"); err != nil {
		t.Fatalf("Error unsetting env var KUBE_CTX_AUTH_INFO: %s", err)
	}
	if err := os.Unsetenv("KUBE_CTX_CLUSTER"); err != nil {
		t.Fatalf("Error unsetting env var KUBE_CTX_CLUSTER: %s", err)
	}
	if err := os.Unsetenv("KUBE_HOST"); err != nil {
		t.Fatalf("Error unsetting env var KUBE_HOST: %s", err)
	}
	if err := os.Unsetenv("KUBE_USER"); err != nil {
		t.Fatalf("Error unsetting env var KUBE_USER: %s", err)
	}
	if err := os.Unsetenv("KUBE_PASSWORD"); err != nil {
		t.Fatalf("Error unsetting env var KUBE_PASSWORD: %s", err)
	}
	if err := os.Unsetenv("KUBE_CLIENT_CERT_DATA"); err != nil {
		t.Fatalf("Error unsetting env var KUBE_CLIENT_CERT_DATA: %s", err)
	}
	if err := os.Unsetenv("KUBE_CLIENT_KEY_DATA"); err != nil {
		t.Fatalf("Error unsetting env var KUBE_CLIENT_KEY_DATA: %s", err)
	}
	if err := os.Unsetenv("KUBE_CLUSTER_CA_CERT_DATA"); err != nil {
		t.Fatalf("Error unsetting env var KUBE_CLUSTER_CA_CERT_DATA: %s", err)
	}

	return func() {
		if err := os.Setenv("KUBE_CONFIG", e.Config); err != nil {
			t.Fatalf("Error resetting env var KUBE_CONFIG: %s", err)
		}
		if err := os.Setenv("KUBECONFIG", e.Config); err != nil {
			t.Fatalf("Error resetting env var KUBECONFIG: %s", err)
		}
		if err := os.Setenv("KUBE_CTX", e.Config); err != nil {
			t.Fatalf("Error resetting env var KUBE_CTX: %s", err)
		}
		if err := os.Setenv("KUBE_CTX_AUTH_INFO", e.CtxAuthInfo); err != nil {
			t.Fatalf("Error resetting env var KUBE_CTX_AUTH_INFO: %s", err)
		}
		if err := os.Setenv("KUBE_CTX_CLUSTER", e.CtxCluster); err != nil {
			t.Fatalf("Error resetting env var KUBE_CTX_CLUSTER: %s", err)
		}
		if err := os.Setenv("KUBE_HOST", e.Host); err != nil {
			t.Fatalf("Error resetting env var KUBE_HOST: %s", err)
		}
		if err := os.Setenv("KUBE_USER", e.User); err != nil {
			t.Fatalf("Error resetting env var KUBE_USER: %s", err)
		}
		if err := os.Setenv("KUBE_PASSWORD", e.Password); err != nil {
			t.Fatalf("Error resetting env var KUBE_PASSWORD: %s", err)
		}
		if err := os.Setenv("KUBE_CLIENT_CERT_DATA", e.ClientCertData); err != nil {
			t.Fatalf("Error resetting env var KUBE_CLIENT_CERT_DATA: %s", err)
		}
		if err := os.Setenv("KUBE_CLIENT_KEY_DATA", e.ClientKeyData); err != nil {
			t.Fatalf("Error resetting env var KUBE_CLIENT_KEY_DATA: %s", err)
		}
		if err := os.Setenv("KUBE_CLUSTER_CA_CERT_DATA", e.ClusterCACertData); err != nil {
			t.Fatalf("Error resetting env var KUBE_CLUSTER_CA_CERT_DATA: %s", err)
		}
	}
}

func getEnv() *currentEnv {
	e := &currentEnv{
		Ctx:               os.Getenv("KUBE_CTX"),
		CtxAuthInfo:       os.Getenv("KUBE_CTX_AUTH_INFO"),
		CtxCluster:        os.Getenv("KUBE_CTX_CLUSTER"),
		Host:              os.Getenv("KUBE_HOST"),
		User:              os.Getenv("KUBE_USER"),
		Password:          os.Getenv("KUBE_PASSWORD"),
		ClientCertData:    os.Getenv("KUBE_CLIENT_CERT_DATA"),
		ClientKeyData:     os.Getenv("KUBE_CLIENT_KEY_DATA"),
		ClusterCACertData: os.Getenv("KUBE_CLUSTER_CA_CERT_DATA"),
	}
	if cfg := os.Getenv("KUBE_CONFIG"); cfg != "" {
		e.Config = cfg
	}
	if cfg := os.Getenv("KUBECONFIG"); cfg != "" {
		e.Config = cfg
	}
	return e
}

type currentEnv struct {
	Config            string
	Ctx               string
	CtxAuthInfo       string
	CtxCluster        string
	Host              string
	User              string
	Password          string
	ClientCertData    string
	ClientKeyData     string
	ClusterCACertData string
}
